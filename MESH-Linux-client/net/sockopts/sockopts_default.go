// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux

package sockopts

import (
	"tailscale.com/types/nettype"
)

// SetBufferSize sets pconn's buffer to size for direction. size may be silently
// capped depending on platform.
//
// errForce is only relevant for Linux, and will always be nil otherwise,
// but we maintain a consistent cross-platform API.
//
// If pconn is not a [*net.UDPConn], then SetBufferSize is no-op.
func SetBufferSize(pconn nettype.PacketConn, direction BufferDirection, size int) (errForce error, errPortable error) {
	return nil, portableSetBufferSize(pconn, direction, size)
}
