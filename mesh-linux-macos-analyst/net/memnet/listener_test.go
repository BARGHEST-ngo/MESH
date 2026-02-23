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

import (
	"context"
	"testing"
)

func TestListener(t *testing.T) {
	l := Listen("srv.local")
	defer l.Close()
	go func() {
		c, err := l.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
	}()

	if c, err := l.Dial(context.Background(), "tcp", "invalid"); err == nil {
		c.Close()
		t.Fatalf("dial to invalid address succeeded")
	}
	c, err := l.Dial(context.Background(), "tcp", "srv.local")
	if err != nil {
		t.Fatalf("dial failed: %v", err)
		return
	}
	c.Close()
}
