// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
