// Copyright (c) 2021 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkctr builds the Tailscale OCI containers.
package main

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type logf func(msg string, args ...interface{})

func withPrefix(f logf, prefix string) logf {
	return func(format string, args ...interface{}) {
		f(prefix+format, args...)
	}
}

// parseFiles parses a comma-separated list of colon-separated pairs
// into a map of filePathOnDisk -> filePathInContainer.
func parseFiles(s string) (map[string]string, error) {
	ret := map[string]string{}
	if len(s) == 0 {
		return ret, nil
	}
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		fs := strings.Split(f, ":")
		if len(fs) != 2 {
			return nil, fmt.Errorf("unparseable file field %q", f)
		}
		ret[fs[0]] = fs[1]
	}
	return ret, nil
}

func parseRepos(reg, tags []string) ([]name.Tag, error) {
	var refs []name.Tag
	for _, rs := range reg {
		r, err := name.NewRepository(rs)
		if err != nil {
			return nil, err
		}
		for _, t := range tags {
			refs = append(refs, r.Tag(t))
		}
	}
	return refs, nil
}

type buildParams struct {
	baseImage   string
	goPaths     map[string]string
	staticFiles map[string]string
	imageRefs   []name.Tag
	publish     bool
	ldflags     string
	gotags      string
	target      string
	verbose     bool
	annotations map[string]string // OCI image annotations
}

func main() {
	var (
		baseImage   = flag.String("base", "", "base image for container")
		gopaths     = flag.String("gopaths", "", "comma-separated list of go paths in src:dst form")
		files       = flag.String("files", "", "comma-separated list of static files in src:dst form")
		repos       = flag.String("repos", "", "comma-separated list of image registries")
		tagArg      = flag.String("tags", "", "comma-separated tags")
		ldflagsArg  = flag.String("ldflags", "", "the --ldflags value to pass to go")
		gotags      = flag.String("gotags", "", "the --tags value to pass to go")
		push        = flag.Bool("push", false, "publish the image")
		target      = flag.String("target", "", "build for a specific env (options: flyio, local)")
		verbose     = flag.Bool("v", false, "verbose build output")
		annotations = flag.String("annotations", "", `OCI image annotations https://github.com/opencontainers/image-spec/blob/main/annotations.md.
		Annotations must be comma separated key=value pairs, i.e key1=val1,key2=val2. For a single image manifest annotations will get added to the image manifest.
		For an image index (a multi-platform manifest list) annotations will get added to each image manifest as well as the image index.
		Annotations with empty values are not supported.`)
	)
	flag.Parse()
	if *tagArg == "" {
		log.Fatal("tags must be set")
	}
	if *repos == "" {
		log.Fatal("registries must be set")
	}
	if *baseImage == "" {
		log.Fatal("baseImage must be set")
	}
	switch *target {
	case "", "flyio", "local":
	default:
		log.Fatalf("unsupported target %q", *target)
	}
	refs, err := parseRepos(strings.Split(*repos, ","), strings.Split(*tagArg, ","))
	if err != nil {
		log.Fatal(err)
	}
	paths, err := parseFiles(*gopaths)
	if err != nil {
		log.Fatal(err)
	}
	staticFiles, err := parseFiles(*files)
	if err != nil {
		log.Fatal(err)
	}
	if len(paths) == 0 && len(staticFiles) == 0 {
		log.Fatal("at least one of --files or --gopaths must be set")
	}

	bp := &buildParams{
		baseImage:   *baseImage,
		goPaths:     paths,
		staticFiles: staticFiles,
		imageRefs:   refs,
		publish:     *push,
		ldflags:     *ldflagsArg,
		gotags:      *gotags,
		target:      *target,
		verbose:     *verbose,
		annotations: parseAnnotations(*annotations),
	}

	if err := fetchAndBuild(bp); err != nil {
		log.Fatal(err)
	}
}

func fetchBaseImage(baseImage string, opts ...remote.Option) (*remote.Descriptor, error) {
	baseRef, err := name.ParseReference(baseImage)
	if err != nil {
		return nil, err
	}
	desc, err := remote.Get(baseRef, opts...)
	if err != nil {
		return nil, err
	}
	return desc, nil
}

// canRunLocal reports whether the platform can run the binary locally, to be
// used by the local target.
func canRunLocal(p v1.Platform) bool {
	if p.OS != "linux" {
		return false
	}
	if runtime.GOOS == "linux" {
		return p.Architecture == runtime.GOARCH
	}
	if runtime.GOOS == "darwin" {
		// macOS can run amd64 linux binaries in docker.
		return p.Architecture == "amd64"
	}
	return false
}

func verifyPlatform(p v1.Platform, target string) error {
	if p.OS != "linux" {
		return fmt.Errorf("unsupported OS: %v", p.OS)
	}
	if target == "local" && !canRunLocal(p) {
		return fmt.Errorf("not required for target %q", target)
	}
	if target == "flyio" && p.Architecture != "amd64" {
		return fmt.Errorf("not required for target %q", target)
	}
	switch p.Architecture {
	case "arm", "arm64", "amd64", "386":
	default:
		return fmt.Errorf("unsupported arch: %v", p.Architecture)
	}
	return nil
}

func fetchAndBuild(bp *buildParams) error {
	ctx := context.Background()
	logf := log.Printf
	remoteOpts := []remote.Option{
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
		remote.WithContext(ctx),
	}
	baseDesc, err := fetchBaseImage(bp.baseImage, remoteOpts...)
	if err != nil {
		return err
	}

	switch baseDesc.MediaType {
	case types.OCIManifestSchema1, types.DockerManifestSchema2:
		// baseRef is an image.
		// Special case to make it only build for that one platform.
		baseImage, err := baseDesc.Image()
		if err != nil {
			return err
		}

		config, err := baseImage.ConfigFile()
		if err != nil {
			return fmt.Errorf("error getting config: %w", err)
		}
		if config.Architecture == "" || config.OS == "" {
			return fmt.Errorf("unknown platform for image: %v", bp.baseImage)
		}

		p := v1.Platform{
			OS:           config.OS,
			Architecture: config.Architecture,
		}
		if config.Variant != "" {
			p.Variant = config.Variant
		}

		if err := verifyPlatform(p, bp.target); err != nil {
			return err
		}
		logf := withPrefix(logf, fmt.Sprintf("%v/%v: ", p.OS, p.Architecture))
		img, err := createImageForBase(bp, logf, baseImage, p)
		if err != nil {
			return err
		}
		if !bp.publish {
			logf("not pushing")
			return nil
		}

		img = mutate.Annotations(img, bp.annotations).(v1.Image) // OCI annotations

		for _, r := range bp.imageRefs {
			if bp.target == "local" {
				if err := loadLocalImage(logf, r, img); err != nil {
					return err
				}
				continue
			}
			logf("pushing to %v", r)
			if err := remote.Write(r, img, remoteOpts...); err != nil {
				return err
			}
		}
		return nil
	case types.OCIImageIndex, types.DockerManifestList:
		// baseRef is a multi-platform index, rest of the method handles this.
	default:
		return fmt.Errorf("failed to interpret base as index or image: %v", baseDesc.MediaType)
	}
	baseIndex, err := baseDesc.ImageIndex()
	if err != nil {
		return err
	}

	im, err := baseIndex.IndexManifest()
	if err != nil {
		return fmt.Errorf("failed to interpret base as index: %w", err)
	}
	var adds []mutate.IndexAddendum
	// Try to build images for all supported platforms.
	for _, id := range im.Manifests {
		logf := withPrefix(logf, fmt.Sprintf("%v/%v: ", id.Platform.OS, id.Platform.Architecture))
		if id.Platform == nil {
			return fmt.Errorf("unknown platform for image: %v", bp.baseImage)
		}
		if err := verifyPlatform(*id.Platform, bp.target); err != nil {
			logf("skipping: %v", err)
			continue
		}
		logf("base digest: %v", id.Digest)
		bi, err := baseIndex.Image(id.Digest)
		if err != nil {
			return err
		}
		logf("building")
		img, err := createImageForBase(bp, logf, bi, *id.Platform)
		if err != nil {
			return err
		}

		// Ensure that any provided OCI annotations are added to each OCI image manifest.
		img = mutate.Annotations(img, bp.annotations).(v1.Image)

		if args := flag.Args(); len(args) > 0 {
			img, err = mutate.Config(img, v1.Config{
				Cmd: args,
			})
			if err != nil {
				return err
			}
		}
		d, err := img.Digest()
		if err != nil {
			return err
		}
		logf("new digest: %v", d)
		adds = append(adds, mutate.IndexAddendum{
			Add: img,
			Descriptor: v1.Descriptor{
				MediaType: id.MediaType,
				URLs:      id.URLs,
				Platform:  id.Platform,
			},
		})
	}
	switch len(adds) {
	case 0:
		logf("no images")
		return nil
	case 1:
		// Don't use a manifest for a single image.
		img := adds[0].Add.(v1.Image)
		d, err := img.Digest()
		if err != nil {
			return err
		}
		logf("image digest: %v", d)
		if !bp.publish {
			logf("not pushing")
			return nil
		}

		for _, r := range bp.imageRefs {
			if bp.target == "local" {
				if err := loadLocalImage(logf, r, img); err != nil {
					return err
				}
				continue
			}
			logf("pushing to %v", r)
			if err := remote.Write(r, img, remoteOpts...); err != nil {
				return err
			}
		}
		return nil
	}
	if bp.target == "local" {
		return fmt.Errorf("cannot build multi-platform images for local target")
	}
	// Generate a new 'fat manifest' with all the platform images. If we are
	// at this point the base was either a Dokcer manifest list or an OCI
	// image index- make sure the new manifest of that type.
	idx := mutate.AppendManifests(mutate.IndexMediaType(empty.Index, baseDesc.MediaType), adds...)
	d, err := idx.Digest()
	if err != nil {
		return err
	}

	// Add any provided OCI annotations to the image index.
	idx = mutate.Annotations(idx, bp.annotations).(v1.ImageIndex)

	logf("index digest: %v", d)
	if !bp.publish {
		logf("not pushing")
		return nil
	}

	for _, r := range bp.imageRefs {
		logf("pushing to %v", r)
		if err := remote.WriteIndex(r, idx, remoteOpts...); err != nil {
			return err
		}
	}

	return nil
}

func goarm(platform v1.Platform) (string, error) {
	if platform.Architecture != "arm" {
		return "", fmt.Errorf("not arm: %v", platform.Architecture)
	}
	v := platform.Variant
	if len(v) != 2 {
		return "", fmt.Errorf("unexpected varient: %v", v)
	}
	if v[0] != 'v' || !('0' <= v[1] && v[1] <= '9') {
		return "", fmt.Errorf("unexpected varient: %v", v)
	}
	return string(v[1]), nil
}

func createImageForBase(bp *buildParams, logf logf, base v1.Image, platform v1.Platform) (v1.Image, error) {
	tmpDir, err := os.MkdirTemp("", "mkctr")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	env := append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS="+platform.OS,
		"GOARCH="+platform.Architecture,
	)
	if platform.Architecture == "arm" {
		v, err := goarm(platform)
		if err != nil {
			return nil, err
		}
		env = append(env, "GOARM="+v)
	}

	files := map[string]string{}
	for src, dst := range bp.staticFiles {
		files[src] = dst
	}

	// Compile all the goPaths
	for gp, dst := range bp.goPaths {
		logf("compiling %v", gp)
		n, err := compileGoBinary(gp, tmpDir, env, bp.ldflags, bp.gotags, bp.verbose)
		if err != nil {
			return nil, err
		}
		logf("output %v -> %v", gp, n)
		files[n] = dst
	}
	// Determine media type of the base image.
	var layerMediaType types.MediaType
	mt, err := base.MediaType()
	if err != nil {
		return nil, fmt.Errorf("error determining base image media type: %w", err)
	}
	switch mt {
	case types.OCIManifestSchema1:
		layerMediaType = types.OCILayer
	case types.DockerManifestSchema2:
		layerMediaType = types.DockerLayer
	default:
		return nil, fmt.Errorf("unknown base image media type %v, accepted types are OCI image manifest v1 (%s) and Docker image manifest v2 (%s)", mt, types.OCIManifestSchema1, types.DockerManifestSchema2)
	}
	layer, err := layerFromFiles(logf, files, layerMediaType)
	if err != nil {
		return nil, err
	}
	return mutate.AppendLayers(base, layer)
}

func compileGoBinary(what, where string, env []string, ldflags, gotags string, verbose bool) (string, error) {
	f, err := os.CreateTemp(where, "out")
	if err != nil {
		return "", err
	}
	out := f.Name()
	if err := f.Close(); err != nil {
		return "", err
	}
	args := []string{
		"build",
		"-trimpath",
	}
	if verbose {
		args = append(args, "-v")
	}
	if len(gotags) > 0 {
		args = append(args, "--tags="+gotags)
	}
	if len(ldflags) > 0 {
		args = append(args, "--ldflags="+ldflags)
	}
	args = append(args,
		"-o="+out,
		what,
	)
	cmd := exec.Command("go", args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func layerFromFiles(logf logf, files map[string]string, layerMediaType types.MediaType) (v1.Layer, error) {
	buf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dirs := make(map[string]bool)
	writeDir := func(dir string) error {
		if dirs[dir] {
			return nil
		}
		logf("creating dir %v", dir)
		if err := tw.WriteHeader(&tar.Header{
			Name:     dir,
			Typeflag: tar.TypeDir,
			Mode:     0555,
			// Set time to 0 to make the images reproducible.
			ModTime: time.Time{},
		}); err != nil {
			return err
		}
		dirs[dir] = true
		return nil
	}
	for src, dst := range files {
		err := filepath.WalkDir(src, func(srcWalk string, d fs.DirEntry, err error) error {
			path := strings.TrimPrefix(srcWalk, src)
			dstWalk := filepath.Join(dst, path)
			writeDir(filepath.Dir(dstWalk))
			if d.IsDir() {
				return writeDir(dstWalk)
			}
			logf("copying %v -> %v", srcWalk, dstWalk)
			return tarFile(tw, srcWalk, dstWalk)
		})
		if err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	binaryLayerBytes := buf.Bytes()
	// An alternative to using tarball.LayerFromOpener would be to use
	// stream.NewLayer
	// https://pkg.go.dev/github.com/google/go-containerregistry@v0.17.0/pkg/v1/stream#NewLayer.
	// This would, however, require us to restructure the code to write each
	// layer to the upstream repository immediately after producing it. At
	// this point we (irbekrm) are not sure if there would be any benefits
	// to switching to stream.NewLayer.
	// https://github.com/google/go-containerregistry/tree/main/pkg/v1/stream#caveats
	return tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(binaryLayerBytes)), nil
	}, tarball.WithCompressedCaching, tarball.WithMediaType(layerMediaType))
}

func tarFile(tw *tar.Writer, src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if err := tw.WriteHeader(&tar.Header{
		Name:     dst,
		Size:     stat.Size(),
		Typeflag: tar.TypeReg,
		Mode:     0555,
		// Set time to 0 to make the images reproducible.
		ModTime: time.Time{},
	}); err != nil {
		return err
	}
	if _, err := io.Copy(tw, file); err != nil {
		return err
	}
	return nil
}

func loadLocalImage(logf logf, tag name.Tag, img v1.Image) error {
	if _, err := daemon.Write(tag, img); err == nil {
		return nil
	}

	// Assume we failed because the docker daemon API is not available, try a
	// CLI option instead.
	var bin string
	if p, err := exec.LookPath("docker"); err == nil {
		bin = p
	} else if p, err = exec.LookPath("podman"); err == nil {
		bin = p
	} else if p, err = exec.LookPath("nerdctl"); err == nil {
		bin = p
	} else {
		return errors.New("no suitable docker CLI-compatible binary found")
	}

	cmd := exec.Command(bin, "image", "load")
	imgReader, imgWriter := io.Pipe()
	defer imgReader.Close()
	go func() {
		defer imgWriter.Close()
		tarball.Write(tag, img, imgWriter)
	}()
	cmd.Stdin = imgReader
	logf("running command: %s", cmd.String())
	out, err := cmd.CombinedOutput()
	logf("output: %s", string(out))
	if err != nil {
		return err
	}

	return nil
}

// parseAnnotations accepts a string with comma separated key=value pairs of annotations i.e key1=val1,key2=val2 and
// returns them as a parsed map.
func parseAnnotations(s string) map[string]string {
	ss := strings.Split(s, ",")
	annotations := make(map[string]string)
	for _, annot := range ss {
		kv := strings.SplitN(annot, "=", 2)
		if len(kv) != 2 {
			continue
		}
		annotations[kv[0]] = kv[1]
	}
	return annotations
}
