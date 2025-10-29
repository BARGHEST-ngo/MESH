// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package sysresources

import (
	"runtime"
	"testing"
)

func TestTotalMemory(t *testing.T) {
	switch runtime.GOOS {
	case "linux":
	case "freebsd", "openbsd", "dragonfly", "netbsd":
	case "darwin":
	default:
		t.Skipf("not supported on runtime.GOOS=%q yet", runtime.GOOS)
	}

	mem := TotalMemory()
	if mem == 0 {
		t.Fatal("wanted TotalMemory > 0")
	}
	t.Logf("total memory: %v bytes", mem)
}
