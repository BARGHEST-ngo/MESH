// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
