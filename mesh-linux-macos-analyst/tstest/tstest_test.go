// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package tstest

import "testing"

func TestReplace(t *testing.T) {
	before := "before"
	done := false
	t.Run("replace", func(t *testing.T) {
		Replace(t, &before, "after")
		if before != "after" {
			t.Errorf("before = %q; want %q", before, "after")
		}
		done = true
	})
	if !done {
		t.Fatal("subtest didn't run")
	}
	if before != "before" {
		t.Errorf("before = %q; want %q", before, "before")
	}
}
