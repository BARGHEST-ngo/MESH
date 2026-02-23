// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package set

// HandleSet is a set of T.
//
// It is not safe for concurrent use.
type HandleSet[T any] map[Handle]T

// Handle is an opaque comparable value that's used as the map key in a
// HandleSet.
type Handle struct {
	v *byte
}

// NewHandle returns a new handle value.
func NewHandle() Handle {
	return Handle{new(byte)}
}

// Add adds the element (map value) e to the set.
//
// It returns a new handle (map key) with which e can be removed, using a map
// delete or the [HandleSet.Delete] method.
func (s *HandleSet[T]) Add(e T) Handle {
	h := NewHandle()
	if *s == nil {
		*s = make(HandleSet[T])
	}
	(*s)[h] = e
	return h
}

// Delete removes the element with handle h from the set.
func (s HandleSet[T]) Delete(h Handle) { delete(s, h) }
