// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package ktimeout

import (
	"time"

	"golang.org/x/sys/unix"
)

// SetUserTimeout sets the TCP_USER_TIMEOUT option on the given file descriptor.
func SetUserTimeout(fd uintptr, timeout time.Duration) error {
	return unix.SetsockoptInt(int(fd), unix.SOL_TCP, unix.TCP_USER_TIMEOUT, int(timeout/time.Millisecond))
}
