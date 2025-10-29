// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package memnet

import "testing"

func TestListenAddressReuse(t *testing.T) {
	var nw Network
	ln1, err := nw.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	if _, err := nw.Listen("tcp", "127.0.0.1:80"); err == nil {
		t.Errorf("listen on in-use address succeeded")
	}
	if err := ln1.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if _, err := nw.Listen("tcp", "127.0.0.1:80"); err != nil {
		t.Errorf("listen on same address after close failed: %v", err)
	}
}
