// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package ipn

import (
	"bytes"
	"iter"
	"sync"
	"testing"

	"tailscale.com/util/mak"
)

type memStore struct {
	mu     sync.Mutex
	writes int
	m      map[StateKey][]byte
}

func (s *memStore) ReadState(k StateKey) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return bytes.Clone(s.m[k]), nil
}

func (s *memStore) WriteState(k StateKey, v []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	mak.Set(&s.m, k, bytes.Clone(v))
	s.writes++
	return nil
}

func (s *memStore) All() iter.Seq2[StateKey, []byte] {
	return func(yield func(StateKey, []byte) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for k, v := range s.m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func TestWriteState(t *testing.T) {
	var ss StateStore = new(memStore)
	WriteState(ss, "foo", []byte("bar"))
	WriteState(ss, "foo", []byte("bar"))
	got, err := ss.ReadState("foo")
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte("bar"); !bytes.Equal(got, want) {
		t.Errorf("got %q; want %q", got, want)
	}
	if got, want := ss.(*memStore).writes, 1; got != want {
		t.Errorf("got %d writes; want %d", got, want)
	}
}
