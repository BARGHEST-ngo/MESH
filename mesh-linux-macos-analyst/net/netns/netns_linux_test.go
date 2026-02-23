// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package netns

import (
	"testing"
)

func TestSocketMarkWorks(t *testing.T) {
	_ = socketMarkWorks()
	// we cannot actually assert whether the test runner has SO_MARK available
	// or not, as we don't know. We're just checking that it doesn't panic.
}
