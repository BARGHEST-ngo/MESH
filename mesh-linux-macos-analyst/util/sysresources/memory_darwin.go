// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build darwin

package sysresources

import "golang.org/x/sys/unix"

func totalMemoryImpl() uint64 {
	val, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0
	}
	return val
}
