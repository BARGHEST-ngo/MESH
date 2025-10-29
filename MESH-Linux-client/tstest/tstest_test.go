// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
