// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
