// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package filelogger

import "testing"

func TestRemoveDatePrefix(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"\n", "\n"},
		{"2009/01/23 01:23:23", "2009/01/23 01:23:23"},
		{"2009/01/23 01:23:23 \n", "\n"},
		{"2009/01/23 01:23:23 foo\n", "foo\n"},
		{"9999/01/23 01:23:23 foo\n", "foo\n"},
		{"2009_01/23 01:23:23 had an underscore\n", "2009_01/23 01:23:23 had an underscore\n"},
	}
	for i, tt := range tests {
		got := removeDatePrefix([]byte(tt.in))
		if string(got) != tt.want {
			t.Logf("[%d] removeDatePrefix(%q) = %q; want %q", i, tt.in, got, tt.want)
		}
	}

}
