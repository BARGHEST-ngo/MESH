// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
