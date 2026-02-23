// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package portmapper registers support for NAT-PMP, PCP, and UPnP port
// mapping protocols to help get direction connections through NATs.
package portmapper

import (
	"tailscale.com/net/netmon"
	"tailscale.com/net/portmapper"
	"tailscale.com/net/portmapper/portmappertype"
	"tailscale.com/types/logger"
	"tailscale.com/util/eventbus"
)

func init() {
	portmappertype.HookNewPortMapper.Set(newPortMapper)
}

func newPortMapper(
	logf logger.Logf,
	bus *eventbus.Bus,
	netMon *netmon.Monitor,
	disableUPnPOrNil func() bool,
	onlyTCP443OrNil func() bool) portmappertype.Client {

	pm := portmapper.NewClient(portmapper.Config{
		EventBus: bus,
		Logf:     logf,
		NetMon:   netMon,
		DebugKnobs: &portmapper.DebugKnobs{
			DisableAll:      onlyTCP443OrNil,
			DisableUPnPFunc: disableUPnPOrNil,
		},
	})
	pm.SetGatewayLookupFunc(netMon.GatewayAndSelfIP)
	return pm
}
