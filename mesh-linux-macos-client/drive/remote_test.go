// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package drive

import (
	"fmt"
	"testing"
)

func TestNormalizeShareName(t *testing.T) {
	tests := []struct {
		name string
		want string
		err  error
	}{
		{
			name: "  (_this is A 5 nAme )_ ",
			want: "(_this is a 5 name )_",
		},
		{
			name: "",
			err:  ErrInvalidShareName,
		},
		{
			name: "generally good except for .",
			err:  ErrInvalidShareName,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("name %q", tt.name), func(t *testing.T) {
			got, err := NormalizeShareName(tt.name)
			if tt.err != nil && err != tt.err {
				t.Errorf("wanted error %v, got %v", tt.err, err)
			} else if got != tt.want {
				t.Errorf("wanted %q, got %q", tt.want, got)
			}
		})
	}
}
