// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package magicsock

import (
	"errors"
	"net"
	"net/netip"
	"sync"
	"syscall"
	"time"
)

// blockForeverConn is a net.PacketConn whose reads block until it is closed.
type blockForeverConn struct {
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
}

func (c *blockForeverConn) ReadFromUDPAddrPort(p []byte) (n int, addr netip.AddrPort, err error) {
	c.mu.Lock()
	for !c.closed {
		c.cond.Wait()
	}
	c.mu.Unlock()
	return 0, netip.AddrPort{}, net.ErrClosed
}

func (c *blockForeverConn) WriteToUDPAddrPort(p []byte, addr netip.AddrPort) (int, error) {
	// Silently drop writes.
	return len(p), nil
}

func (c *blockForeverConn) LocalAddr() net.Addr {
	// Return a *net.UDPAddr because lots of code assumes that it will.
	return new(net.UDPAddr)
}

func (c *blockForeverConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return net.ErrClosed
	}
	c.closed = true
	c.cond.Broadcast()
	return nil
}

func (c *blockForeverConn) SetDeadline(t time.Time) error         { return errors.New("unimplemented") }
func (c *blockForeverConn) SetReadDeadline(t time.Time) error     { return errors.New("unimplemented") }
func (c *blockForeverConn) SetWriteDeadline(t time.Time) error    { return errors.New("unimplemented") }
func (c *blockForeverConn) SyscallConn() (syscall.RawConn, error) { return nil, errUnsupportedConnType }
