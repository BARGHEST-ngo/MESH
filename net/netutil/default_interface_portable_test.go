// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
