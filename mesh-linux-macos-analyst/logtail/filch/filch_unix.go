// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !windows && !wasm && !plan9 && !tamago

package filch

import (
	"os"

	"golang.org/x/sys/unix"
)

func saveStderr() (*os.File, error) {
	fd, err := unix.Dup(stderrFD)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), "stderr"), nil
}

func unsaveStderr(f *os.File) error {
	err := dup2Stderr(f)
	f.Close()
	return err
}

func dup2Stderr(f *os.File) error {
	return unix.Dup2(int(f.Fd()), stderrFD)
}
