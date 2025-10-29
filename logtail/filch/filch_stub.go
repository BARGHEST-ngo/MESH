// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build wasm || plan9 || tamago

package filch

import (
	"os"
)

func saveStderr() (*os.File, error) {
	return os.Stderr, nil
}

func unsaveStderr(f *os.File) error {
	os.Stderr = f
	return nil
}

func dup2Stderr(f *os.File) error {
	return nil
}
