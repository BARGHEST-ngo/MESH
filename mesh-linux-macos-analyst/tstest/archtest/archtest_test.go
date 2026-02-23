// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package archtest

import (
	"runtime"
	"testing"

	"gvisor.dev/gvisor/pkg/atomicbitops"
)

// tests netstack's AlignedAtomicInt64.
func TestAlignedAtomicInt64(t *testing.T) {
	type T struct {
		A atomicbitops.Int64
		_ int32
		B atomicbitops.Int64
	}

	t.Logf("I am %v/%v\n", runtime.GOOS, runtime.GOARCH)
	var x T
	x.A.Store(1)
	x.B.Store(2)
	if got, want := x.A.Load(), int64(1); got != want {
		t.Errorf("A = %v; want %v", got, want)
	}
	if got, want := x.B.Load(), int64(2); got != want {
		t.Errorf("A = %v; want %v", got, want)
	}
}
