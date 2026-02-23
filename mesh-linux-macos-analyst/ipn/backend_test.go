// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package ipn

import (
	"testing"

	"tailscale.com/health"
	"tailscale.com/types/empty"
)

func TestNotifyString(t *testing.T) {
	for _, tt := range []struct {
		name     string
		value    Notify
		expected string
	}{
		{
			name:     "notify-empty",
			value:    Notify{},
			expected: "Notify{}",
		},
		{
			name:     "notify-with-login-finished",
			value:    Notify{LoginFinished: &empty.Message{}},
			expected: "Notify{LoginFinished}",
		},
		{
			name:     "notify-with-multiple-fields",
			value:    Notify{LoginFinished: &empty.Message{}, Health: &health.State{}},
			expected: "Notify{LoginFinished Health{...}}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.value.String()
			if actual != tt.expected {
				t.Fatalf("expected=%q, actual=%q", tt.expected, actual)
			}
		})
	}
}
