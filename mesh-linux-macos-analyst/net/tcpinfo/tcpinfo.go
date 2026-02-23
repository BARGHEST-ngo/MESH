// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package tcpinfo provides platform-agnostic accessors to information about a
// TCP connection (e.g. RTT, MSS, etc.).
package tcpinfo

import (
	"errors"
	"net"
	"time"
)

var (
	ErrNotTCP        = errors.New("tcpinfo: not a TCP conn")
	ErrUnimplemented = errors.New("tcpinfo: unimplemented")
)

// RTT returns the RTT for the given net.Conn.
//
// If the net.Conn is not a *net.TCPConn and cannot be unwrapped into one, then
// ErrNotTCP will be returned. If retrieving the RTT is not supported on the
// current platform, ErrUnimplemented will be returned.
func RTT(conn net.Conn) (time.Duration, error) {
	tcpConn, err := unwrap(conn)
	if err != nil {
		return 0, err
	}

	return rttImpl(tcpConn)
}

// netConner is implemented by crypto/tls.Conn to unwrap into an underlying
// net.Conn.
type netConner interface {
	NetConn() net.Conn
}

// unwrap attempts to unwrap a net.Conn into an underlying *net.TCPConn
func unwrap(nc net.Conn) (*net.TCPConn, error) {
	for {
		switch v := nc.(type) {
		case *net.TCPConn:
			return v, nil
		case netConner:
			nc = v.NetConn()
		default:
			return nil, ErrNotTCP
		}
	}
}
