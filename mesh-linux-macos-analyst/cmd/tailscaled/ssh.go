// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build (linux || darwin || freebsd || openbsd || plan9) && !ts_omit_ssh

package main

// Force registration of tailssh with LocalBackend.
import _ "tailscale.com/ssh/tailssh"
