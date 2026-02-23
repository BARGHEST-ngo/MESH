// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package vms

import (
	"net/netip"
	"runtime"
	"testing"

	"tailscale.com/net/netmon"
)

func deriveBindhost(t *testing.T) string {
	t.Helper()

	ifName, err := netmon.DefaultRouteInterface()
	if err != nil {
		t.Fatal(err)
	}

	var ret string
	err = netmon.ForeachInterfaceAddress(func(i netmon.Interface, prefix netip.Prefix) {
		if ret != "" || i.Name != ifName {
			return
		}
		ret = prefix.Addr().String()
	})
	if ret != "" {
		return ret
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal("can't find a bindhost")
	return "unreachable"
}

func TestDeriveBindhost(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("requires GOOS=linux")
	}
	t.Log(deriveBindhost(t))
}
