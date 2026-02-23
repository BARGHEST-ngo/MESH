// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package sockopts contains logic for applying socket options.
package sockopts

import (
	"net"
	"runtime"

	"tailscale.com/types/nettype"
)

// BufferDirection represents either the read/receive or write/send direction
// of a socket buffer.
type BufferDirection string

const (
	ReadDirection  BufferDirection = "read"
	WriteDirection BufferDirection = "write"
)

func portableSetBufferSize(pconn nettype.PacketConn, direction BufferDirection, size int) error {
	if runtime.GOOS == "plan9" {
		// Not supported. Don't try. Avoid logspam.
		return nil
	}
	var err error
	if c, ok := pconn.(*net.UDPConn); ok {
		if direction == WriteDirection {
			err = c.SetWriteBuffer(size)
		} else {
			err = c.SetReadBuffer(size)
		}
	}
	return err
}
