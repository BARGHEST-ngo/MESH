// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
