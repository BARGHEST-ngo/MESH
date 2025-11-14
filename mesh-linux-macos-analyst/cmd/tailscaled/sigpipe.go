// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build go1.21 && !plan9

package main

import "syscall"

func init() {
	sigPipe = syscall.SIGPIPE
}
