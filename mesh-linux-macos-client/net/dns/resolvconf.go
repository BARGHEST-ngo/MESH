// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build linux || freebsd || openbsd

package dns

import (
	"bytes"
	"os/exec"
)

func resolvconfStyle() string {
	if _, err := exec.LookPath("resolvconf"); err != nil {
		return ""
	}
	output, err := exec.Command("resolvconf", "--version").CombinedOutput()
	if err != nil {
		// Debian resolvconf doesn't understand --version, and
		// exits with a specific error code.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 99 {
			return "debian"
		}
	}
	if bytes.HasPrefix(output, []byte("Debian resolvconf")) {
		return "debian"
	}
	// Treat everything else as openresolv, by far the more popular implementation.
	return "openresolv"
}
