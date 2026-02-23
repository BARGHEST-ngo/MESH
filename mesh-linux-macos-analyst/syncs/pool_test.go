// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package syncs

import "testing"

func TestPool(t *testing.T) {
	var pool Pool[string]
	s := pool.Get() // should not panic
	if s != "" {
		t.Fatalf("got %q, want %q", s, "")
	}
	pool.New = func() string { return "new" }
	s = pool.Get()
	if s != "new" {
		t.Fatalf("got %q, want %q", s, "new")
	}
	var found bool
	for range 1000 {
		pool.Put("something")
		found = pool.Get() == "something"
		if found {
			break
		}
	}
	if !found {
		t.Fatalf("unable to get any value put in the pool")
	}
}
