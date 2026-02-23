// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package derphttp

func SetTestHookWatchLookConnectResult(f func(connectError error, wasSelfConnect bool) (keepRunning bool)) {
	testHookWatchLookConnectResult = f
}

// breakConnection breaks the connection, which should trigger a reconnect.
func (c *Client) BreakConnection(brokenClient *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client != brokenClient.client {
		return
	}
	if c.netConn != nil {
		c.netConn.Close()
		c.netConn = nil
	}
	c.client = nil
}

var RetryInterval = &retryInterval
