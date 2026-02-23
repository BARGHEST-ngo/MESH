// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package metrics

import (
	"io/fs"
	"sync"

	"go4.org/mem"
	"tailscale.com/util/dirwalk"
)

// counter is a reusable counter for counting file descriptors.
type counter struct {
	n int

	// cb is the (*counter).count method value. Creating it allocates,
	// so we have to save it away and use a sync.Pool to keep currentFDs
	// amortized alloc-free.
	cb func(name mem.RO, de fs.DirEntry) error
}

var counterPool = &sync.Pool{New: func() any {
	c := new(counter)
	c.cb = c.count
	return c
}}

func (c *counter) count(name mem.RO, de fs.DirEntry) error {
	c.n++
	return nil
}

func currentFDs() int {
	c := counterPool.Get().(*counter)
	defer counterPool.Put(c)
	c.n = 0
	dirwalk.WalkShallow(mem.S("/proc/self/fd"), c.cb)
	return c.n
}
