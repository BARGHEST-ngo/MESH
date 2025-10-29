// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build linux && !android && ts_package_container

package hostinfo

import (
	"testing"
)

func TestInContainer(t *testing.T) {
	if got := inContainer(); !got.EqualBool(true) {
		t.Errorf("inContainer = %v; want true due to ts_package_container build tag", got)
	}
}
