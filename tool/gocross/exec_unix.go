// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build unix

package main

import "golang.org/x/sys/unix"

func doExec(cmd string, args []string, env []string) error {
	return unix.Exec(cmd, args, env)
}
