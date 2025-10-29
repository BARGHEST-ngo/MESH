// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package netmon

import "testing"

func BenchmarkGetPACWindows(b *testing.B) {
	b.ReportAllocs()
	for i := range b.N {
		v := getPACWindows()
		if i == 0 {
			b.Logf("Got: %q", v)
		}
	}
}
