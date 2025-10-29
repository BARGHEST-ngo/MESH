// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package ktimeout

import (
	"time"

	"golang.org/x/sys/unix"
)

// SetUserTimeout sets the TCP_USER_TIMEOUT option on the given file descriptor.
func SetUserTimeout(fd uintptr, timeout time.Duration) error {
	return unix.SetsockoptInt(int(fd), unix.SOL_TCP, unix.TCP_USER_TIMEOUT, int(timeout/time.Millisecond))
}
