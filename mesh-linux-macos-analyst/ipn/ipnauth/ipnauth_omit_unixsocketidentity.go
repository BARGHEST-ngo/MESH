// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
