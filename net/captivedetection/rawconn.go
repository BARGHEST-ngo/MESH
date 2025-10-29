// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !(ios || darwin)

package captivedetection

import (
	"syscall"

	"tailscale.com/types/logger"
)

// setSocketInterfaceIndex sets the IP_BOUND_IF socket option on the given RawConn.
// This forces the socket to use the given interface.
func setSocketInterfaceIndex(c syscall.RawConn, ifIndex int, logf logger.Logf) error {
	// No-op on non-Darwin platforms.
	return nil
}
