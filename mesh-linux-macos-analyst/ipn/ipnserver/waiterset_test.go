// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package ipnserver

import (
	"context"
	"sync"
	"testing"
)

func TestWaiterSet(t *testing.T) {
	var s waiterSet

	wantLen := func(want int, when string) {
		t.Helper()
		if got := len(s); got != want {
			t.Errorf("%s: len = %v; want %v", when, got, want)
		}
	}
	wantLen(0, "initial")
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())

	ready, cleanup := s.add(&mu, ctx)
	wantLen(1, "after add")

	select {
	case <-ready:
		t.Fatal("should not be ready")
	default:
	}
	s.wakeAll()
	<-ready

	wantLen(1, "after fire")
	cleanup()
	wantLen(0, "after cleanup")

	// And again but on an already-expired ctx.
	cancel()
	ready, cleanup = s.add(&mu, ctx)
	<-ready // shouldn't block
	cleanup()
	wantLen(0, "at end")
}
