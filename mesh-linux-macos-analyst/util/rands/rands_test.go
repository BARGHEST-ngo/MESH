// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package rands

import "testing"

func TestHexString(t *testing.T) {
	for i := 0; i <= 8; i++ {
		s := HexString(i)
		if len(s) != i {
			t.Errorf("HexString(%v) = %q; want len %v, not %v", i, s, i, len(s))
		}
	}
}
