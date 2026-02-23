// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build tailscale_go && (darwin || ios)

package sockstats

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func init() {
	tcpConnStats = darwinTcpConnStats
}

func darwinTcpConnStats(c syscall.RawConn) (tx, rx uint64) {
	c.Control(func(fd uintptr) {
		if rawInfo, err := unix.GetsockoptTCPConnectionInfo(
			int(fd),
			unix.IPPROTO_TCP,
			unix.TCP_CONNECTION_INFO,
		); err == nil {
			tx = uint64(rawInfo.Txbytes)
			rx = uint64(rawInfo.Rxbytes)
		}
	})
	return
}
