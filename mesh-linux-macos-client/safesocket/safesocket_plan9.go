// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
