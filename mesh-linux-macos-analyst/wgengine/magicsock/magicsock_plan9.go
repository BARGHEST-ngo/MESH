// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build plan9

package magicsock

// shouldRebind returns if the error is one that is known to be healed by a
// rebind, and if so also returns a resason string for the rebind.
func shouldRebind(err error) (ok bool, reason string) {
	return false, ""
}
