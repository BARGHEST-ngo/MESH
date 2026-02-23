// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package execqueue

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestExecQueue(t *testing.T) {
	ctx := context.Background()
	var n atomic.Int32
	q := &ExecQueue{}
	defer q.Shutdown()
	q.Add(func() { n.Add(1) })
	q.Wait(ctx)
	if got := n.Load(); got != 1 {
		t.Errorf("n=%d; want 1", got)
	}
}
