// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package netutil

import (
	"testing"
)

func TestDefaultInterfacePortable(t *testing.T) {
	ifName, addr, err := DefaultInterfacePortable()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Default interface: %s", ifName)
	t.Logf("Default address: %s", addr)

	if ifName == "" {
		t.Fatal("Default interface name is empty")
	}
	if !addr.IsValid() {
		t.Fatal("Default address is invalid")
	}
}
