// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build linux

package sysresources

import "golang.org/x/sys/unix"

func totalMemoryImpl() uint64 {
	var info unix.Sysinfo_t

	if err := unix.Sysinfo(&info); err != nil {
		return 0
	}

	// uint64 casts are required since these might be uint32s
	return uint64(info.Totalram) * uint64(info.Unit)
}
