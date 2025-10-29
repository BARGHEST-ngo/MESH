// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
