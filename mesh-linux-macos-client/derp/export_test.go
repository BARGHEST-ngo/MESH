// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package derp

import "time"

func (c *Client) RecvTimeoutForTest(timeout time.Duration) (m ReceivedMessage, err error) {
	return c.recvTimeout(timeout)
}
