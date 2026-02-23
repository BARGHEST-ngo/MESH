// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_not_in_tests

package envknob

import "runtime"

func GOOS() string {
	// When the "ts_not_in_tests" build tag is used, we define this func to just
	// return a simple constant so callers optimize just as if the knob were not
	// present. We can then build production/optimized builds with the
	// "ts_not_in_tests" build tag.
	return runtime.GOOS
}
