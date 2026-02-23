// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package osrouter

import (
	"tailscale.com/net/netmon"
	"tailscale.com/types/logger"
	"tailscale.com/wgengine/router"
)

func init() {
	router.HookCleanUp.Set(func(logf logger.Logf, netMon *netmon.Monitor, ifName string) {
		cleanUp(logf, ifName)
	})
}

func cleanUp(logf logger.Logf, interfaceName string) {
	// If the interface was left behind, ifconfig down will not remove it.
	// In fact, this will leave a system in a tainted state where starting tailscaled
	// will result in "interface tailscale0 already exists"
	// until the defunct interface is ifconfig-destroyed.
	ifup := []string{"ifconfig", interfaceName, "destroy"}
	if out, err := cmd(ifup...).CombinedOutput(); err != nil {
		logf("ifconfig destroy: %v\n%s", err, out)
	}
}
