//
// Copyright 2024 Stacklok, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package store provides utilities to work with a cache store.
package store

import (
	"github.com/puzpuzpuz/xsync"
)

// RefCacher is an interface for caching references.
type RefCacher interface {
	Store(key, value string)
	Load(key string) (string, bool)
}

type refCacher struct {
	cache *xsync.MapOf[string, string]
}

// NewRefCacher returns a new RefCacher. The default implementation is
// thread-safe.
func NewRefCacher() RefCacher {
	return &refCacher{
		cache: xsync.NewMapOf[string](),
	}
}

// Store stores a key-value pair.
func (r *refCacher) Store(key, value string) {
	r.cache.Store(key, value)
}

// Load loads a value for a given key.
func (r *refCacher) Load(key string) (string, bool) {
	return r.cache.Load(key)
}

type unsafeCacher struct {
	cache map[string]string
}

// NewUnsafeCacher returns a new RefCacher that's not thread-safe.
func NewUnsafeCacher() RefCacher {
	return &unsafeCacher{
		cache: map[string]string{},
	}
}

// Store stores a key-value pair.
func (r *unsafeCacher) Store(key, value string) {
	r.cache[key] = value
}

// Load loads a value for a given key.
func (r *unsafeCacher) Load(key string) (string, bool) {
	v, ok := r.cache[key]
	return v, ok
}
