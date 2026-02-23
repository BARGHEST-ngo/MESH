// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
