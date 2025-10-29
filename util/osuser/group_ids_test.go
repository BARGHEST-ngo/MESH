// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package osuser

import (
	"slices"
	"testing"
)

func TestParseGroupIds(t *testing.T) {
	tests := []struct {
		in       string
		expected []string
	}{
		{"5000\x005001\n", []string{"5000", "5001"}},
		{"5000\n", []string{"5000"}},
		{"\n", []string{""}},
	}
	for _, test := range tests {
		actual := parseGroupIds([]byte(test.in))
		if !slices.Equal(actual, test.expected) {
			t.Errorf("parseGroupIds(%q) = %q, wanted %s", test.in, actual, test.expected)
		}
	}
}
