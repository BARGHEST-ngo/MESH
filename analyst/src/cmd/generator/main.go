// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

func main() {
	// Get absolute path to current module's go.mod file
	output, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		log.Fatalf("failed to get GOMOD: %v", err)
	}
	goModPath := string(output[:len(output)-1]) // remove trailing newline
	rootModPath := filepath.Dir(goModPath)
	analystSrcDir := filepath.Join(rootModPath, "analyst", "src")
	patchesDir := filepath.Join(analystSrcDir, "patches")
	tailscaleDir := filepath.Join(rootModPath, "tailscale")

	// If tailscale dir already exists, copy out .git directory and remove it
	var gitDir string
	if _, err := os.Stat(tailscaleDir); err == nil {
		gitDir = filepath.Join(tailscaleDir, ".git")
		if err := os.Rename(gitDir, filepath.Join(analystSrcDir, "tailscale.git")); err != nil {
			log.Fatalf("failed to move .git directory: %v", err)
		}
		if err := os.RemoveAll(tailscaleDir); err != nil {
			log.Fatalf("failed to remove tailscale directory: %v", err)
		}
	}

	// Parse go.mod to get tailscale.com version
	modContents, err := os.ReadFile(goModPath)
	if err != nil {
		log.Fatalf("failed to read go.mod: %v", err)
	}
	modFile, err := modfile.Parse("go.mod", modContents, nil)
	if err != nil {
		log.Fatalf("failed to parse go.mod: %v", err)
	}
	var tailscaleVersion string
	for _, req := range modFile.Require {
		if req.Mod.Path == "tailscale.com" {
			tailscaleVersion = req.Mod.Version
			break
		}
	}
	if tailscaleVersion == "" {
		log.Fatalf("tailscale.com not found in go.mod")
	}
	log.Printf("Found tailscale.com version %q in go.mod", tailscaleVersion)

	// Download the module and get the path to its verified zip file.
	// go mod verify checksums the zip (not the extracted cache directory),
	// so extracting directly from the zip avoids trusting the cache.
	output, err = exec.Command("go", "mod", "download", "-json", "tailscale.com@"+tailscaleVersion).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to download tailscale.com module: %v\nOutput: %s", err, output)
	}
	var modInfo struct {
		Zip string `json:"Zip"`
	}
	if err := json.Unmarshal(output, &modInfo); err != nil {
		log.Fatalf("failed to parse go mod download output: %v", err)
	}
	if modInfo.Zip == "" {
		log.Fatalf("go mod download did not return a Zip path")
	}
	log.Printf("Downloaded module zip: %s", modInfo.Zip)

	// Verify module checksums against go.sum before extracting.
	cmd := exec.Command("go", "mod", "verify")
	if verifyOutput, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("go mod verify failed (module may have been tampered with): %v\nOutput: %s", err, verifyOutput)
	}
	log.Println("Module checksums verified against go.sum")

	// Extract the verified zip directly to the build directory.
	// Module zips prefix all entries with "module@version/"; we strip that.
	if err := extractModuleZip(modInfo.Zip, tailscaleDir); err != nil {
		log.Fatalf("failed to extract module zip: %v", err)
	}
	log.Printf("Extracted module zip to %q", tailscaleDir)

	// Add files from cli directory to the tailscale.com/cmd/tailscale/cli package in the build directory
	analystCliDir := filepath.Join(analystSrcDir, "cli")
	tailscaleCliDir := filepath.Join(tailscaleDir, "cmd", "tailscale", "cli")
	analystFiles, err := os.ReadDir(analystCliDir)
	if err != nil {
		log.Fatalf("failed to read analyst/src/cli directory: %v", err)
	}

	log.Println("Copying analyst/src/cli files to tailscale.com/cmd/tailscale/cli in build directory...")
	for _, f := range analystFiles {
		// Skip directories and non-Go files
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}

		srcPath := filepath.Join(analystCliDir, f.Name())
		destPath := filepath.Join(tailscaleCliDir, f.Name())

		cmd := exec.Command("cp", srcPath, destPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("failed to copy %q to %q: %v\nOutput: %s", srcPath, destPath, err, output)
		}
		if err := os.Chmod(destPath, 0644); err != nil {
			log.Fatalf("failed to set permissions on %q: %v", destPath, err)
		}
		log.Printf("Copied %q to %q\n", f.Name(), destPath)
	}

	// Add files from wgcfg directory to the tailscale.com/wgengine/wgcfg package in the build directory
	analystWgcfgDir := filepath.Join(analystSrcDir, "wgcfg")
	tailscaleWgcfgDir := filepath.Join(tailscaleDir, "wgengine", "wgcfg")
	analystFiles, err = os.ReadDir(analystWgcfgDir)
	if err != nil {
		log.Fatalf("failed to read analyst/src/wgcfg directory: %v", err)
	}

	log.Println("Copying analyst/src/wgcfg files to tailscale.com/wgengine/wgcfg in build directory...")
	for _, f := range analystFiles {
		// Skip directories and non-Go files
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}

		srcPath := filepath.Join(analystWgcfgDir, f.Name())
		destPath := filepath.Join(tailscaleWgcfgDir, f.Name())

		cmd := exec.Command("cp", srcPath, destPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("failed to copy %q to %q: %v\nOutput: %s", srcPath, destPath, err, output)
		}
		if err := os.Chmod(destPath, 0644); err != nil {
			log.Fatalf("failed to set permissions on %q: %v", destPath, err)
		}
		log.Printf("Copied %q to %q\n", f.Name(), destPath)
	}

	// Apply any necessary patches
	log.Println("Patching files...")
	patches, err := os.ReadDir(patchesDir)
	if err != nil {
		log.Fatalf("failed to read patches directory: %v", err)
	}
	for _, p := range patches {
		patchPath := filepath.Join(patchesDir, p.Name())
		cmd := exec.Command("patch", "-p1", "-i", patchPath, "-d", tailscaleDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("failed to apply patch %q: %v\nOutput: %s", patchPath, err, output)
		}
		log.Printf("Applied patch %q\n", patchPath)
	}

	log.Println("Replacing upstream Tailscale DNS fallback servers with empty set...")
	dnsFallbackPath := filepath.Join(tailscaleDir, "net", "dnsfallback", "dns-fallback-servers.json")
	if err := os.WriteFile(dnsFallbackPath, []byte(`{"Regions": {}}`+"\n"), 0644); err != nil {
		log.Fatalf("failed to write dns-fallback-servers.json: %v", err)
	}

	log.Println("Replacing upstream Tailscale hostnames with example.com...")
	hostReplacements := map[string]string{
		"log.tailscale.com":          "log.example.com",
		"controlplane.tailscale.com": "controlplane.example.com",
		"login.tailscale.com":        "login.example.com",
	}
	err = filepath.WalkDir(tailscaleDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "tstest" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		replaced, err := replaceOutsideComments(content, hostReplacements)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		if !bytes.Equal(replaced, content) {
			if err := os.WriteFile(path, replaced, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			log.Printf("Patched Tailscale hostnames in %s", path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to replace Tailscale hostnames: %v", err)
	}

	// Run go mod tidy in the build directory
	log.Println("Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tailscaleDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to run go mod tidy: %v\nOutput: %s", err, output)
	}

	// Replace .git directory if we moved it out earlier
	if gitDir != "" {
		if err := os.Rename(filepath.Join(analystSrcDir, "tailscale.git"), gitDir); err != nil {
			log.Fatalf("failed to move .git directory back: %v", err)
		}
		log.Printf("Restored .git directory\n")
	}

	// Reset file permissions (https://stackoverflow.com/a/4408378)
	cmd = exec.Command(
		"/bin/sh", "-c",
		"git diff -p -R --no-ext-diff --no-color --diff-filter=M | grep -A 1 -B 1 \"old mode\" --color=never | grep -E \"^(diff|(old|new) mode)\" --color=never | git apply --allow-empty",
	)
	cmd.Dir = tailscaleDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to reset permissions on copied tailscale.com module: %v\nOutput: %s", err, output)
	}

	log.Println("Done.")
}

// replaceOutsideComments applies string replacements to Go source code,
// skipping comments so that documentation and notes are preserved unchanged.
func replaceOutsideComments(src []byte, replacements map[string]string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Collect comment byte ranges, sorted by position.
	type span struct{ start, end int }
	var comments []span
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			comments = append(comments, span{
				start: fset.Position(c.Pos()).Offset,
				end:   fset.Position(c.End()).Offset,
			})
		}
	}
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].start < comments[j].start
	})

	// replaceSegment applies all replacements to a code segment.
	replaceSegment := func(segment string) string {
		for old, repl := range replacements {
			segment = strings.ReplaceAll(segment, old, repl)
		}
		return segment
	}

	// Build result: apply replacements only to non-comment segments.
	var result []byte
	pos := 0
	for _, c := range comments {
		result = append(result, replaceSegment(string(src[pos:c.start]))...)
		result = append(result, src[c.start:c.end]...)
		pos = c.end
	}
	result = append(result, replaceSegment(string(src[pos:]))...)

	return result, nil
}

// extractModuleZip extracts a Go module zip to destDir, stripping the
// "module@version/" prefix that all entries share per the module zip spec.
func extractModuleZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		// Strip the "module@version/" prefix from each entry.
		_, relPath, ok := strings.Cut(f.Name, "/")
		if !ok || relPath == "" {
			continue
		}
		destPath := filepath.Join(destDir, relPath)

		// Verify the path doesn't escape the destination (zip slip protection).
		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("zip entry %q would escape destination directory", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			continue
		}

		if err := extractZipFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a zip archive to destPath,
// creating parent directories as needed.
func extractZipFile(f *zip.File, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("creating parent directory for %s: %w", destPath, err)
	}

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening zip entry %s: %w", f.Name, err)
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating %s: %w", destPath, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return fmt.Errorf("writing %s: %w", destPath, err)
	}

	return nil
}
