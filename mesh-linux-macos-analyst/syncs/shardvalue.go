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

// TODO(raggi): this implementation is still imperfect as it will still result
// in cross CPU sharing periodically, we instead really want a per-CPU shard
// key, but the limitations of calling platform code make reaching for even the
// getcpu vdso very painful. See https://github.com/golang/go/issues/18802, and
// hopefully one day we can replace with a primitive that falls out of that
// work.

// ShardValue contains a value sharded over a set of shards.
// In order to be useful, T should be aligned to cache lines.
// Users must organize that usage in One and All is concurrency safe.
// The zero value is not safe for use; use [NewShardValue].
type ShardValue[T any] struct {
	shards []T

	//lint:ignore U1000 unused under tailscale_go builds.
	pool shardValuePool
}

// Len returns the number of shards.
func (sp *ShardValue[T]) Len() int {
	return len(sp.shards)
}

// All yields a pointer to the value in each shard.
func (sp *ShardValue[T]) All(yield func(*T) bool) {
	for i := range sp.shards {
		if !yield(&sp.shards[i]) {
			return
		}
	}
}
