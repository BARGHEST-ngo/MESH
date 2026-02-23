// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package safesocket

import "testing"

func TestLocalTCPPortAndToken(t *testing.T) {
	// Just test that it compiles for now (is available on all platforms).
	port, token, err := LocalTCPPortAndToken()
	t.Logf("got %v, %s, %v", port, token, err)
}
