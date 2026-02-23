// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
