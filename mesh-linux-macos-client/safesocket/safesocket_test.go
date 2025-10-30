// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package safesocket

import "testing"

func TestLocalTCPPortAndToken(t *testing.T) {
	// Just test that it compiles for now (is available on all platforms).
	port, token, err := LocalTCPPortAndToken()
	t.Logf("got %v, %s, %v", port, token, err)
}
