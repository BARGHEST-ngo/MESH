// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package ptr contains the ptr.To function.
package ptr

// To returns a pointer to a shallow copy of v.
func To[T any](v T) *T {
	return &v
}
