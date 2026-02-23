// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package tcpinfo

import (
	"net"
	"time"

	"golang.org/x/sys/unix"
)

func rttImpl(conn *net.TCPConn) (time.Duration, error) {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return 0, err
	}

	var (
		tcpInfo *unix.TCPConnectionInfo
		sysErr  error
	)
	err = rawConn.Control(func(fd uintptr) {
		tcpInfo, sysErr = unix.GetsockoptTCPConnectionInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_CONNECTION_INFO)
	})
	if err != nil {
		return 0, err
	} else if sysErr != nil {
		return 0, sysErr
	}

	return time.Duration(tcpInfo.Rttcur) * time.Millisecond, nil
}
