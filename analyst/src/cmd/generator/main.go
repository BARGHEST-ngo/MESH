// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tsDir := filepath.Join(root, "tailscale")

	info, err := os.Stat(tsDir)
	if err != nil || !info.IsDir() {
		log.Fatalf("tailscale submodule missing: %s", tsDir)
	}

	log.Println("tailscale submodule present")
	log.Println("No patching required")
}
