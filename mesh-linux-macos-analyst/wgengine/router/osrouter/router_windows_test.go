// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package osrouter

import (
	"path/filepath"
	"testing"
)

func TestGetNetshPath(t *testing.T) {
	ft := &firewallTweaker{
		logf: t.Logf,
	}
	path := ft.getNetshPath()
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path for netsh.exe: %q", path)
	}
}
