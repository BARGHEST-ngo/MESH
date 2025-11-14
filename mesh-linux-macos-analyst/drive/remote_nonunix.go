// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !unix

package drive

func doAllowShareAs() bool {
	// On non-UNIX platforms, we use the GUI application (e.g. Windows taskbar
	// icon) to access the filesystem as whatever unprivileged user is running
	// the GUI app, so we cannot allow sharing as a different user.
	return false
}
