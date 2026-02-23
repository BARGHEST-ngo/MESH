// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package ace registers support for Alternate Connectivity Endpoints (ACE).
package ace

import (
	"net/netip"

	"tailscale.com/control/controlhttp"
	"tailscale.com/net/ace"
	"tailscale.com/net/netx"
)

func init() {
	controlhttp.HookMakeACEDialer.Set(mkDialer)
}

func mkDialer(dialer netx.DialFunc, aceHost string, optIP netip.Addr) netx.DialFunc {
	return (&ace.Dialer{
		ACEHost:   aceHost,
		ACEHostIP: optIP, // may be zero
		NetDialer: dialer,
	}).Dial
}
