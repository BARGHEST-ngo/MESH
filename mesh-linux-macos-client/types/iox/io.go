// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package iox provides types to implement [io] functionality.
package iox

// TODO(https://go.dev/issue/21670): Deprecate or remove this functionality
// once the Go language supports implementing an 1-method interface directly
// using a function value of a matching signature.

// ReaderFunc implements [io.Reader] using the underlying function value.
type ReaderFunc func([]byte) (int, error)

func (f ReaderFunc) Read(b []byte) (int, error) {
	return f(b)
}

// WriterFunc implements [io.Writer] using the underlying function value.
type WriterFunc func([]byte) (int, error)

func (f WriterFunc) Write(b []byte) (int, error) {
	return f(b)
}
