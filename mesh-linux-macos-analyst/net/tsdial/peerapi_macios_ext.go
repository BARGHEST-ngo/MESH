// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// This file's built on iOS and on two of three macOS build variants:
// the two GUI variants that both use Extensions (Network Extension
// and System Extension). It's not used on tailscaled-on-macOS.

//go:build ts_macext && (darwin || ios)

package tsdial

import (
	"errors"
	"net"
	"syscall"

	"tailscale.com/net/netns"
)

func init() {
	peerDialControlFunc = peerDialControlFuncNetworkExtension
}

func peerDialControlFuncNetworkExtension(d *Dialer) func(network, address string, c syscall.RawConn) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := -1
	if x, ok := interfaceIndexLocked(d); ok {
		index = x
	}
	var lc net.ListenConfig
	netns.SetListenConfigInterfaceIndex(&lc, index)
	return func(network, address string, c syscall.RawConn) error {
		if index == -1 {
			return errors.New("failed to find TUN interface to bind to")
		}
		return lc.Control(network, address, c)
	}
}

func interfaceIndexLocked(d *Dialer) (index int, ok bool) {
	if d.netMon == nil {
		return 0, false
	}
	st := d.netMon.InterfaceState()
	iface, ok := st.Interface[d.tunName]
	if !ok {
		return 0, false
	}
	return iface.Index, true
}
