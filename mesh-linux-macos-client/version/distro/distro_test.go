// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package distro

import "testing"

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	var d Distro
	for range b.N {
		d = Get()
	}
	_ = d
}
