// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux && !darwin

package tcpinfo

import (
	"net"
	"time"
)

func rttImpl(conn *net.TCPConn) (time.Duration, error) {
	return 0, ErrUnimplemented
}
