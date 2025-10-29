// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package ktimeout

import (
	"context"
	"fmt"
	"net"
	"time"
)

func ExampleUserTimeout() {
	lc := net.ListenConfig{
		Control: UserTimeout(30 * time.Second),
	}
	l, err := lc.Listen(context.TODO(), "tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	l.Close()
	// Output:
}
