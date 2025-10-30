// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
