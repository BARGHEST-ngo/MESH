// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ios || darwin

package captivedetection

import (
	"syscall"

	"golang.org/x/sys/unix"
	"tailscale.com/types/logger"
)

// setSocketInterfaceIndex sets the IP_BOUND_IF socket option on the given RawConn.
// This forces the socket to use the given interface.
func setSocketInterfaceIndex(c syscall.RawConn, ifIndex int, logf logger.Logf) error {
	return c.Control((func(fd uintptr) {
		err := unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, ifIndex)
		if err != nil {
			logf("captivedetection: failed to set IP_BOUND_IF (ifIndex=%d): %v", ifIndex, err)
		}
	}))
}
