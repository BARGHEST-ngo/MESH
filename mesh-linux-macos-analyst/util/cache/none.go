// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package cache

// None provides no caching and always calls the provided FillFunc.
//
// It is safe for concurrent use if the underlying FillFunc is.
type None[K comparable, V any] struct{}

var _ Cache[int, int] = None[int, int]{}

// Get always calls the provided FillFunc and returns what it does.
func (c None[K, V]) Get(_ K, f FillFunc[V]) (V, error) {
	v, _, e := f()
	return v, e
}

// Forget implements Cache.
func (None[K, V]) Forget(K) {}

// Empty implements Cache.
func (None[K, V]) Empty() {}
