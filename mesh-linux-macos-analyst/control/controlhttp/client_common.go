// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package controlhttp

import (
	"tailscale.com/control/controlbase"
)

// ClientConn is a Tailscale control client as returned by the Dialer.
//
// It's effectively just a *controlbase.Conn (which it embeds) with
// optional metadata.
type ClientConn struct {
	// Conn is the noise connection.
	*controlbase.Conn
}
