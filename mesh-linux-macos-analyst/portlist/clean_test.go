// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package portlist

import "testing"

func TestArgvSubject(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{
			in:   nil,
			want: "",
		},
		{
			in:   []string{"/usr/bin/sshd"},
			want: "sshd",
		},
		{
			in:   []string{"/bin/mono"},
			want: "mono",
		},
		{
			in:   []string{"/nix/store/x2cw2xjw98zdysf56bdlfzsr7cyxv0jf-mono-5.20.1.27/bin/mono", "/bin/exampleProgram.exe"},
			want: "exampleProgram",
		},
		{
			in:   []string{"/bin/mono", "/sbin/exampleProgram.bin"},
			want: "exampleProgram.bin",
		},
		{
			in:   []string{"/usr/bin/sshd_config [listener] 1 of 10-100 startups"},
			want: "sshd_config",
		},
		{
			in:   []string{"/usr/bin/sshd [listener] 0 of 10-100 startups"},
			want: "sshd",
		},
		{
			in:   []string{"/opt/aws/bin/eic_run_authorized_keys %u %f -o AuthorizedKeysCommandUser ec2-instance-connect [listener] 0 of 10-100 startups"},
			want: "eic_run_authorized_keys",
		},
		{
			in:   []string{"/usr/bin/nginx worker"},
			want: "nginx",
		},
	}

	for _, test := range tests {
		got := argvSubject(test.in...)
		if got != test.want {
			t.Errorf("argvSubject(%v) = %q, want %q", test.in, got, test.want)
		}
	}
}
