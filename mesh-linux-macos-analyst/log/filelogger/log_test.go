// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
