// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// This is the Taildrive stuff that should ideally be registered in init only when
// the ts_omit_drive is not set, but for transition reasons is currently (2025-09-08)
// always defined, as we work to pull it out of LocalBackend.

package ipnlocal

import "tailscale.com/tailcfg"

const (
	// DriveLocalPort is the port on which the Taildrive listens for location
	// connections on quad 100.
	DriveLocalPort = 8080
)

// DriveSharingEnabled reports whether sharing to remote nodes via Taildrive is
// enabled. This is currently based on checking for the drive:share node
// attribute.
func (b *LocalBackend) DriveSharingEnabled() bool {
	return b.currentNode().SelfHasCap(tailcfg.NodeAttrsTaildriveShare)
}

// DriveAccessEnabled reports whether accessing Taildrive shares on remote nodes
// is enabled. This is currently based on checking for the drive:access node
// attribute.
func (b *LocalBackend) DriveAccessEnabled() bool {
	return b.currentNode().SelfHasCap(tailcfg.NodeAttrsTaildriveAccess)
}
