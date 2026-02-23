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

package resolver

import (
	"errors"
	"net"

	"tailscale.com/net/netmon"
	"tailscale.com/net/netns"
)

func init() {
	initListenConfig = initListenConfigNetworkExtension
}

func initListenConfigNetworkExtension(nc *net.ListenConfig, netMon *netmon.Monitor, tunName string) error {
	nif, ok := netMon.InterfaceState().Interface[tunName]
	if !ok {
		return errors.New("utun not found")
	}
	return netns.SetListenConfigInterfaceIndex(nc, nif.Interface.Index)
}
