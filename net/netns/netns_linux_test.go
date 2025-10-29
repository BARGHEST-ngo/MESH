// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package netns

import (
	"testing"
)

func TestSocketMarkWorks(t *testing.T) {
	_ = socketMarkWorks()
	// we cannot actually assert whether the test runner has SO_MARK available
	// or not, as we don't know. We're just checking that it doesn't panic.
}
