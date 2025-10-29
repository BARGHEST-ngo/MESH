// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build darwin

package hostinfo

import (
	"os"
	"path/filepath"
)

func init() {
	packageType = packageTypeDarwin
}

func packageTypeDarwin() string {
	// Using tailscaled or IPNExtension?
	exe, _ := os.Executable()
	return filepath.Base(exe)
}
