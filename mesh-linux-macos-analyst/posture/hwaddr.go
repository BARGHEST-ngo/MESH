// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package posture

import (
	"net/netip"
	"slices"

	"tailscale.com/net/netmon"
)

// GetHardwareAddrs returns the hardware addresses of all non-loopback
// network interfaces.
func GetHardwareAddrs() (hwaddrs []string, err error) {
	err = netmon.ForeachInterface(func(i netmon.Interface, _ []netip.Prefix) {
		if i.IsLoopback() {
			return
		}
		if a := i.HardwareAddr.String(); a != "" {
			hwaddrs = append(hwaddrs, a)
		}
	})
	slices.Sort(hwaddrs)
	return slices.Compact(hwaddrs), err
}
