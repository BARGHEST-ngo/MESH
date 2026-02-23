// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package sysresources

// TotalMemory returns the total accessible system memory, in bytes. If the
// value cannot be determined, then 0 will be returned.
func TotalMemory() uint64 {
	return totalMemoryImpl()
}
