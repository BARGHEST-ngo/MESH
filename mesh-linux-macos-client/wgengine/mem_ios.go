// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package wgengine

import (
	"github.com/amnezia-vpn/amneziawg-go/device"
)

// iOS has a very restrictive memory limit on network extensions.
// Reduce the maximum amount of memory that wireguard-go can allocate
// to avoid getting killed.

func init() {
	device.QueueStagedSize = 64
	device.QueueOutboundSize = 64
	device.QueueInboundSize = 64
	device.QueueHandshakeSize = 64
	device.PreallocatedBuffersPerPool = 64
}
