// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
