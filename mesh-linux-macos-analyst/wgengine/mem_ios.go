// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
