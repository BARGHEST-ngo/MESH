// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
