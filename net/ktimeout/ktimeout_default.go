// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux

package ktimeout

import (
	"time"
)

// SetUserTimeout is a no-op on this platform.
func SetUserTimeout(fd uintptr, timeout time.Duration) error {
	return nil
}
