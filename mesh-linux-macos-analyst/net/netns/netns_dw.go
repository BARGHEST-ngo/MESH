// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build darwin || windows

package netns

import (
	"errors"
	"net"
	"net/netip"
)

var errUnspecifiedHost = errors.New("unspecified host")

func parseAddress(address string) (addr netip.Addr, err error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		// error means the string didn't contain a port number, so use the string directly
		host = address
	}
	if host == "" {
		return addr, errUnspecifiedHost
	}

	return netip.ParseAddr(host)
}
