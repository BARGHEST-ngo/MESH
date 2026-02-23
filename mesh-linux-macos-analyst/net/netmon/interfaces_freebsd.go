// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// This might work on other BSDs, but only tested on FreeBSD.

//go:build freebsd

package netmon

import (
	"syscall"

	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
)

// fetchRoutingTable calls route.FetchRIB, fetching NET_RT_DUMP.
func fetchRoutingTable() (rib []byte, err error) {
	return route.FetchRIB(syscall.AF_UNSPEC, unix.NET_RT_DUMP, 0)
}

func parseRoutingTable(rib []byte) ([]route.Message, error) {
	return route.ParseRIB(syscall.NET_RT_IFLIST, rib)
}

func getDelegatedInterface(ifIndex int) (int, error) {
	return 0, nil
}
