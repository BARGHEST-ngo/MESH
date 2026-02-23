// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package childproc allows other packages to register "tailscaled be-child"
// child process hook code. This avoids duplicating build tags in the
// tailscaled package. Instead, the code that needs to fork/exec the self
// executable (when it's tailscaled) can instead register the code
// they want to run.
package childproc

var Code = map[string]func([]string) error{}

// Add registers code f to run as 'tailscaled be-child <typ> [args]'.
func Add(typ string, f func(args []string) error) {
	if _, dup := Code[typ]; dup {
		panic("dup hook " + typ)
	}
	Code[typ] = f
}
