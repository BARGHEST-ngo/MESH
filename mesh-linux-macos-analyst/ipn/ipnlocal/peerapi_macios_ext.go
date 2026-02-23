// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_macext && (darwin || ios)

package ipnlocal

import (
	"fmt"
	"net"
	"net/netip"

	"tailscale.com/net/netmon"
	"tailscale.com/net/netns"
)

func init() {
	initListenConfig = initListenConfigNetworkExtension
}

// initListenConfigNetworkExtension configures nc for listening on IP
// through the iOS/macOS Network/System Extension (Packet Tunnel
// Provider) sandbox.
func initListenConfigNetworkExtension(nc *net.ListenConfig, ip netip.Addr, st *netmon.State, tunIfName string) error {
	tunIf, ok := st.Interface[tunIfName]
	if !ok {
		return fmt.Errorf("no interface with name %q", tunIfName)
	}
	return netns.SetListenConfigInterfaceIndex(nc, tunIf.Index)
}
