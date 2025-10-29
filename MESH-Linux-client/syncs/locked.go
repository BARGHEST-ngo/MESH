// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
