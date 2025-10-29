// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build (linux || darwin || freebsd || openbsd || plan9) && !ts_omit_ssh

package main

// Force registration of tailssh with LocalBackend.
import _ "tailscale.com/ssh/tailssh"
