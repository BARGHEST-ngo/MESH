// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux

package batching

import (
	"tailscale.com/types/nettype"
)

// TryUpgradeToConn is no-op on all platforms except linux.
func TryUpgradeToConn(pconn nettype.PacketConn, _ string, _ int) nettype.PacketConn {
	return pconn
}

var controlMessageSize = 0

func MinControlMessageSize() int {
	return controlMessageSize
}

const IdealBatchSize = 1
