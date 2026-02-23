// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build plan9

package safesocket

import (
	"context"
	"net"
)

func connect(_ context.Context, path string) (net.Conn, error) {
	return net.Dial("tcp", "localhost:5252")
}

func listen(path string) (net.Listener, error) {
	return net.Listen("tcp", "localhost:5252")
}
