// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package tstime

import (
	"math/rand/v2"
	"time"
)

// RandomDurationBetween returns a random duration in range [min,max).
// If panics if max < min.
func RandomDurationBetween(min, max time.Duration) time.Duration {
	diff := max - min
	if diff == 0 {
		return min
	}
	return min + rand.N(max-min)
}
