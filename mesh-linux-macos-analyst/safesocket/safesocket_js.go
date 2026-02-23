// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package safesocket

import (
	"context"
	"net"

	"github.com/akutz/memconn"
)

const memName = "Tailscale-IPN"

func listen(path string) (net.Listener, error) {
	return memconn.Listen("memu", memName)
}

func connect(ctx context.Context, _ string) (net.Conn, error) {
	return memconn.DialContext(ctx, "memu", memName)
}
