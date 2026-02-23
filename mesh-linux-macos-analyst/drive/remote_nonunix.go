// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !unix

package drive

func doAllowShareAs() bool {
	// On non-UNIX platforms, we use the GUI application (e.g. Windows taskbar
	// icon) to access the filesystem as whatever unprivileged user is running
	// the GUI app, so we cannot allow sharing as a different user.
	return false
}
