// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package osrouter

import "net/netip"

//lint:ignore U1000 used in Windows/Linux tests only
func mustCIDRs(ss ...string) []netip.Prefix {
	var ret []netip.Prefix
	for _, s := range ss {
		ret = append(ret, netip.MustParsePrefix(s))
	}
	return ret
}
