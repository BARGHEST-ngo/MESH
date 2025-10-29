// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ((linux && !android) || windows || (darwin && !ios) || freebsd) && !ts_omit_cliconndiag

package safesocket

import (
	"strings"

	ps "github.com/mitchellh/go-ps"
)

func init() {
	tailscaledProcExists.Set(func() bool {
		procs, err := ps.Processes()
		if err != nil {
			return false
		}
		for _, proc := range procs {
			name := proc.Executable()
			const tailscaled = "tailscaled"
			if len(name) < len(tailscaled) {
				continue
			}
			// Do case insensitive comparison for Windows,
			// notably, and ignore any ".exe" suffix.
			if strings.EqualFold(name[:len(tailscaled)], tailscaled) {
				return true
			}
		}
		return false
	})
}
