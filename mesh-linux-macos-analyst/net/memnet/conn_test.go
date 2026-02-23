// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package memnet

import (
	"net"
	"testing"

	"golang.org/x/net/nettest"
)

func TestConn(t *testing.T) {
	nettest.TestConn(t, func() (c1 net.Conn, c2 net.Conn, stop func(), err error) {
		c1, c2 = NewConn("test", bufferSize)
		return c1, c2, func() {
			c1.Close()
			c2.Close()
		}, nil
	})
}
