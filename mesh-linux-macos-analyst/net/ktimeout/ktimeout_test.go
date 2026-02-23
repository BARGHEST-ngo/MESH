// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
