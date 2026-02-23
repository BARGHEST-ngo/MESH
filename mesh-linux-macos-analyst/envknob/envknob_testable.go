// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !ts_not_in_tests

package envknob

import "runtime"

// GOOS reports the effective runtime.GOOS to run as.
//
// In practice this returns just runtime.GOOS, unless overridden by
// test TS_DEBUG_FAKE_GOOS.
//
// This allows changing OS-specific stuff like the IPN server behavior
// for tests so we can e.g. test Windows-specific behaviors on Linux.
// This isn't universally used.
func GOOS() string {
	if v := String("TS_DEBUG_FAKE_GOOS"); v != "" {
		return v
	}
	return runtime.GOOS
}
