// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package truncate_test

import (
	"testing"

	"tailscale.com/util/truncate"
)

func TestString(t *testing.T) {
	tests := []struct {
		input string
		size  int
		want  string
	}{
		{"", 1000, ""},                 // n > length
		{"abc", 4, "abc"},              // n > length
		{"abc", 3, "abc"},              // n == length
		{"abcdefg", 4, "abcd"},         // n < length, safe
		{"abcdefg", 0, ""},             // n < length, safe
		{"abc\U0001fc2d", 3, "abc"},    // n < length, at boundary
		{"abc\U0001fc2d", 4, "abc"},    // n < length, mid-rune
		{"abc\U0001fc2d", 5, "abc"},    // n < length, mid-rune
		{"abc\U0001fc2d", 6, "abc"},    // n < length, mid-rune
		{"abc\U0001fc2defg", 7, "abc"}, // n < length, cut multibyte
	}

	for _, tc := range tests {
		got := truncate.String(tc.input, tc.size)
		if got != tc.want {
			t.Errorf("truncate(%q, %d): got %q, want %q", tc.input, tc.size, got, tc.want)
		}
	}
}
