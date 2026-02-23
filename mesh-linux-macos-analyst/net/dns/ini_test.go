// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build windows

package dns

import (
	"reflect"
	"testing"
)

func TestParseIni(t *testing.T) {
	var tests = []struct {
		src  string
		want map[string]map[string]string
	}{
		{
			src: `# appended wsl.conf file
[automount]
	enabled = true
	root=/mnt/
# added by tailscale
[network] # trailing comment
generateResolvConf = false  # trailing comment`,
			want: map[string]map[string]string{
				"automount": {"enabled": "true", "root": "/mnt/"},
				"network":   {"generateResolvConf": "false"},
			},
		},
	}
	for _, test := range tests {
		got := parseIni(test.src)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("for:\n%s\ngot:  %v\nwant: %v", test.src, got, test.want)
		}
	}
}
