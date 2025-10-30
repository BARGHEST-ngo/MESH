// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
