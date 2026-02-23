// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build js

package portmapper

import (
	"context"
	"net/netip"
)

type upnpClient any

type uPnPDiscoResponse struct{}

func parseUPnPDiscoResponse([]byte) (uPnPDiscoResponse, error) {
	return uPnPDiscoResponse{}, nil
}

func processUPnPResponses(metas []uPnPDiscoResponse) []uPnPDiscoResponse {
	return metas
}

func (c *Client) getUPnPPortMapping(
	ctx context.Context,
	gw netip.Addr,
	internal netip.AddrPort,
	prevPort uint16,
) (external netip.AddrPort, ok bool) {
	return netip.AddrPort{}, false
}
