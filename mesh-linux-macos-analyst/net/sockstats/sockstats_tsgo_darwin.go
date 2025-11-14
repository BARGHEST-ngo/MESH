// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
