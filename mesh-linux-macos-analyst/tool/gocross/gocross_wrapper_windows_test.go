// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGocrossWrapper(t *testing.T) {
	for i := range 2 { // once to build gocross; second to test it's cached
		cmd := exec.Command("pwsh", "-NoProfile", "-ExecutionPolicy", "Bypass", ".\\gocross-wrapper.ps1", "version")
		cmd.Env = append(os.Environ(), "CI=true", "NOPWSHDEBUG=false", "TS_USE_GOCROSS=1") // for Set-PSDebug verbosity
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("gocross-wrapper.ps1 failed: %v\n%s", err, out)
		}
		if i > 0 && !strings.Contains(string(out), "$gocrossOk = $true\r\n") {
			t.Errorf("expected to find '$gocrossOk = $true'; got output:\n%s", out)
		}
	}
}
