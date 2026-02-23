// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package goroutines

import "testing"

func TestScrubbedGoroutineDump(t *testing.T) {
	t.Logf("Got:\n%s\n", ScrubbedGoroutineDump(true))
}

func TestScrubHex(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"foo", "foo"},
		{"", ""},
		{"0x", "?_"},
		{"0x001 and same 0x001", "v1%1_ and same v1%1_"},
		{"0x008 and same 0x008", "v1%0_ and same v1%0_"},
		{"0x001 and diff 0x002", "v1%1_ and diff v2%2_"},
	}
	for _, tt := range tests {
		got := scrubHex([]byte(tt.in))
		if string(got) != tt.want {
			t.Errorf("for input:\n%s\n\ngot:\n%s\n\nwant:\n%s\n", tt.in, got, tt.want)
		}
	}
}
