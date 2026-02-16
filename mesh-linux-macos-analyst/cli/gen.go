//go:build ignore

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

const (
	tailscaleCliPkg            = "tailscale.com/cmd/tailscale/cli"
	tailscaleInternalClientPkg = "tailscale.com/internal/client/tailscale"
)

var (
	// Additional CLI source files to copy from Tailscale repo
	cliFiles = []string{
		"appcroutes.go",
		"configure_apple-all.go",
		"configure_apple.go",
		"configure_linux-all.go",
		"configure_linux.go",
		"configure.go",
		"debug.go",
		"dns-query.go",
		"dns-status.go",
		"dns.go",
		"down.go",
		"exitnode.go",
		"id-token.go",
		"ip.go",
		"logout.go",
		"metrics.go",
		"nc.go",
		"netcheck.go",
		"ping.go",
		"risks.go",
		"set.go",
		"ssh_exec_windows.go",
		"ssh_exec.go",
		"ssh_unix.go",
		"ssh.go",
		"up.go",
		"update.go",
		"version.go",
		"whois.go",
	}
)

func main() {
	// go mod verify has already been run by the go:generate directive
	// so we can be sure that the module files we're about to copy
	// haven't been tampered with.
	// TODO: have the pentesters verify this. Do GONOSUMDB or GOSUMDB
	// env vars affect this?

	clean := flag.Bool("clean", false, "Clean the analyst/cli directory instead of copying files")
	flag.Parse()

	// Get absolute path to current module's go.mod file
	output, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		log.Fatalf("failed to get GOMOD: %v", err)
	}
	goModPath := string(output[:len(output)-1]) // remove trailing newline
	rootModPath := filepath.Dir(goModPath)

	if *clean {
		log.Println("Cleaning analyst/cli directory...")
		cliDir := filepath.Join(rootModPath, "analyst", "cli")
		for _, f := range cliFiles {
			path := filepath.Join(cliDir, f)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				log.Fatalf("failed to remove %q: %v", path, err)
			}
			log.Printf("Removed %q\n", path)
		}
		log.Println("Done.")
		return
	}

	pkgs, err := packages.Load(&packages.Config{Dir: rootModPath}, tailscaleCliPkg)
	if err != nil {
		log.Fatalf("failed to load package %q: %v", tailscaleCliPkg, err)
	} else if len(pkgs) != 1 {
		log.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	pkg := pkgs[0]

	// Make sure all expected files are present
	var toCopy []string
	for _, f := range cliFiles {
		found := false
		for _, sf := range append(pkg.GoFiles, pkg.IgnoredFiles...) {
			if filepath.Base(sf) == f {
				found = true
				toCopy = append(toCopy, sf)
				break
			}
		}
		if !found {
			log.Fatalf("file %q not found in package %q", f, tailscaleCliPkg)
		}
	}

	// Copy source files to our analyst/cli directory
	log.Println("Copying files...")
	analystCliPkgPath := filepath.Join(rootModPath, "analyst", "cli")
	for _, srcPath := range toCopy {
		destPath := filepath.Join(analystCliPkgPath, filepath.Base(srcPath))
		cmd := exec.Command("cp", srcPath, destPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("failed to copy %q to %q: %v\nOutput: %s", srcPath, destPath, err, output)
		}
		if err := os.Chmod(destPath, 0644); err != nil {
			log.Fatalf("failed to set permissions on %q: %v", destPath, err)
		}
		log.Printf("Copied %q to %q\n", srcPath, destPath)
	}

	// Apply any necessary patches
	log.Println("Patching files...")
	patches, err := os.ReadDir(filepath.Join(analystCliPkgPath, "patches"))
	if err != nil {
		log.Fatalf("failed to read patches directory: %v", err)
	}
	for _, p := range patches {
		patchPath := filepath.Join("patches", p.Name())
		cmd := exec.Command("patch", "-p1", "-i", patchPath, "-d", analystCliPkgPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("failed to apply patch %q: %v\nOutput: %s", patchPath, err, output)
		}
		log.Printf("Applied patch %q\n", patchPath)
	}

	log.Println("Copying tailscale internal client package...")
	pkgs, err = packages.Load(&packages.Config{Dir: rootModPath}, tailscaleInternalClientPkg)
	if err != nil {
		log.Fatalf("failed to load package %q: %v", tailscaleInternalClientPkg, err)
	} else if len(pkgs) != 1 {
		log.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	pkg = pkgs[0]
	if err := os.MkdirAll(filepath.Join(rootModPath, "analyst", "internal", "client", "tailscale"), 0755); err != nil {
		log.Fatalf("failed to create tailscale internal client directory: %v", err)
	}
	files, err := os.ReadDir(pkg.Dir)
	if err != nil {
		log.Fatalf("failed to read tailscale internal client directory: %v", err)
	}
	internalClientPkgPath := filepath.Join(rootModPath, "analyst", "internal", "client", "tailscale")
	for _, f := range files {
		destPath := filepath.Join(internalClientPkgPath, f.Name())
		cmd := exec.Command("cp", filepath.Join(pkg.Dir, f.Name()), destPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("failed to copy %q to %q: %v\nOutput: %s", f.Name(), destPath, err, output)
		}
		if err := os.Chmod(destPath, 0644); err != nil {
			log.Fatalf("failed to set permissions on %q: %v", destPath, err)
		}
		log.Printf("Copied %q to %q\n", f.Name(), destPath)
	}

	log.Println("Done.")
}

// TODO
// embed version info (including tailscale version) via ldflags -X
// go list -f '{{.Module.Version}}' tailscale.com/cmd/tailscale/cli
