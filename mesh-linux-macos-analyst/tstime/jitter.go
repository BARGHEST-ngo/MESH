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
