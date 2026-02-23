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

import "net/netip"

//lint:ignore U1000 used in Windows/Linux tests only
func mustCIDRs(ss ...string) []netip.Prefix {
	var ret []netip.Prefix
	for _, s := range ss {
		ret = append(ret, netip.MustParsePrefix(s))
	}
	return ret
}
