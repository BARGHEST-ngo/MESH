// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package syncs

import "sync"

// Pool is the generic version of [sync.Pool].
type Pool[T any] struct {
	pool sync.Pool

	// New optionally specifies a function to generate
	// a value when Get would otherwise return the zero value of T.
	// It may not be changed concurrently with calls to Get.
	New func() T
}

// Get selects an arbitrary item from the Pool, removes it from the Pool,
// and returns it to the caller. See [sync.Pool.Get].
func (p *Pool[T]) Get() T {
	x, ok := p.pool.Get().(T)
	if !ok && p.New != nil {
		x = p.New()
	}
	return x
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
