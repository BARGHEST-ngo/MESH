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

import (
	"sync"
)

// AssertLocked panics if m is not locked.
func AssertLocked(m *sync.Mutex) {
	if m.TryLock() {
		m.Unlock()
		panic("mutex is not locked")
	}
}

// AssertRLocked panics if rw is not locked for reading or writing.
func AssertRLocked(rw *sync.RWMutex) {
	if rw.TryLock() {
		rw.Unlock()
		panic("mutex is not locked")
	}
}

// AssertWLocked panics if rw is not locked for writing.
func AssertWLocked(rw *sync.RWMutex) {
	if rw.TryRLock() {
		rw.RUnlock()
		panic("mutex is not rlocked")
	}
}
