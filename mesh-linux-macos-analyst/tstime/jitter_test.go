// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package tstime

import (
	"testing"
	"time"
)

func TestRandomDurationBetween(t *testing.T) {
	if got := RandomDurationBetween(1, 1); got != 1 {
		t.Errorf("between 1 and 1 = %v; want 1", int64(got))
	}
	const min = 1 * time.Second
	const max = 10 * time.Second
	for range 500 {
		if got := RandomDurationBetween(min, max); got < min || got >= max {
			t.Fatalf("%v (%d) out of range", got, got)
		}
	}
}
