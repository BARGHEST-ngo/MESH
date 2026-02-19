package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	// go mod verify should already have been run by go:generate
	// so we can be sure that the module files we're about to copy
	// haven't been tampered with.
	// TODO: have the pentesters verify this. Do GONOSUMDB or GOSUMDB
	// env vars affect this?

	// Get absolute path to current module's go.mod file
	output, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		log.Fatalf("failed to get GOMOD: %v", err)
	}
	goModPath := string(output[:len(output)-1]) // remove trailing newline
	rootModPath := filepath.Dir(goModPath)
	analystDir := filepath.Join(rootModPath, "analyst")
	patchesDir := filepath.Join(analystDir, "patches")
	buildDir := filepath.Join(analystDir, "build")
	tailscaleDir := filepath.Join(buildDir, "tailscale.com")

	// Only proceed if tailscaleDir does not exist
	if _, err := os.Stat(tailscaleDir); err == nil {
		log.Printf("build directory already exists, skipping generation")
		return
	}

	// Load the tailscale.com package to find its module directory
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedModule,
			Dir:  rootModPath,
		},
		"tailscale.com",
	)
	if err != nil {
		log.Fatalf("failed to load package tailscale.com: %v", err)
	} else if len(pkgs) != 1 {
		log.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	modPath := pkgs[0].Module.Dir
	fmt.Println("Located tailscale.com module:", modPath)

	// Copy the tailscale.com module to build directory
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		log.Fatalf("failed to create build directory: %v", err)
	}
	cmd := exec.Command("cp", "-r", modPath, tailscaleDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to copy tailscale.com module to build directory: %v\nOutput: %s", err, output)
	}

	// chmod -R 0755 the copied directory to ensure we have read/write permissions
	if err := filepath.Walk(tailscaleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, 0755)
	}); err != nil {
		log.Fatalf("failed to set permissions on copied tailscale.com module: %v", err)
	}
	log.Printf("Copied tailscale.com module to %q\n", tailscaleDir)

	// Add files from cli directory to the tailscale.com/cmd/tailscale/cli package in the build directory
	analystCliDir := filepath.Join(rootModPath, "analyst", "cli")
	tailscaleCliDir := filepath.Join(tailscaleDir, "cmd", "tailscale", "cli")
	analystFiles, err := os.ReadDir(analystCliDir)
	if err != nil {
		log.Fatalf("failed to read analyst/cli directory: %v", err)
	}

	log.Println("Copying analyst/cli files to tailscale.com/cmd/tailscale/cli in build directory...")
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
	analystWgcfgDir := filepath.Join(rootModPath, "analyst", "wgcfg")
	tailscaleWgcfgDir := filepath.Join(tailscaleDir, "wgengine", "wgcfg")
	analystFiles, err = os.ReadDir(analystWgcfgDir)
	if err != nil {
		log.Fatalf("failed to read analyst/wgcfg directory: %v", err)
	}

	log.Println("Copying analyst/wgcfg files to tailscale.com/wgengine/wgcfg in build directory...")
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

	// Run go mod tidy in the build directory
	log.Println("Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tailscaleDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to run go mod tidy: %v\nOutput: %s", err, output)
	}

	log.Println("Done.")
}
