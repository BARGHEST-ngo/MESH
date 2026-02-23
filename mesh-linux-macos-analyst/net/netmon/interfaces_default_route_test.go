// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build linux || (darwin && !ts_macext)

package netmon

import (
	"testing"
)

func TestDefaultRouteInterface(t *testing.T) {
	// tests /proc/net/route on the local system, cannot make an assertion about
	// the correct interface name, but good as a sanity check.
	v, err := DefaultRouteInterface()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got %q", v)
}
