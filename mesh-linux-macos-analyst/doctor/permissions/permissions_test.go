// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package permissions

import "testing"

func TestPermissionsImpl(t *testing.T) {
	if err := permissionsImpl(t.Logf); err != nil {
		t.Error(err)
	}
}
