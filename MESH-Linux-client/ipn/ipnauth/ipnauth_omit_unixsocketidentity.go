// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !windows && ts_omit_unixsocketidentity

package ipnauth

import (
	"net"

	"tailscale.com/types/logger"
)

// GetConnIdentity extracts the identity information from the connection
// based on the user who owns the other end of the connection.
// and couldn't. The returned connIdentity has NotWindows set to true.
func GetConnIdentity(_ logger.Logf, c net.Conn) (ci *ConnIdentity, err error) {
	return &ConnIdentity{conn: c, notWindows: true}, nil
}

// WindowsToken is unsupported when GOOS != windows and always returns
// ErrNotImplemented.
func (ci *ConnIdentity) WindowsToken() (WindowsToken, error) {
	return nil, ErrNotImplemented
}
