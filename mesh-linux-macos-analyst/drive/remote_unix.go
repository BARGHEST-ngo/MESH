// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build unix

package drive

import "tailscale.com/version"

func doAllowShareAs() bool {
	// All UNIX platforms use user servers (sub-processes) to access the OS
	// filesystem as a specific unprivileged users, except for sandboxed macOS
	// which doesn't support impersonating users and instead accesses files
	// through the macOS GUI app as whatever unprivileged user is running it.
	return !version.IsSandboxedMacOS()
}
