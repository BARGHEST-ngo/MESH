// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !plan9

package magicsock

import (
	"errors"
	"syscall"
)

// shouldRebind returns if the error is one that is known to be healed by a
// rebind, and if so also returns a resason string for the rebind.
func shouldRebind(err error) (ok bool, reason string) {
	switch {
	// EPIPE/ENOTCONN are common errors when a send fails due to a closed
	// socket. There is some platform and version inconsistency in which
	// error is returned, but the meaning is the same.
	case errors.Is(err, syscall.EPIPE), errors.Is(err, syscall.ENOTCONN):
		return true, "broken-pipe"

	// EPERM is typically caused by EDR software, and has been observed to be
	// transient, it seems that some versions of some EDR lose track of sockets
	// at times, and return EPERM, but reconnects will establish appropriate
	// rights associated with a new socket.
	case errors.Is(err, syscall.EPERM):
		return true, "operation-not-permitted"
	}
	return false, ""
}
