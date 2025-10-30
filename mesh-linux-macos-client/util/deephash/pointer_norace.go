// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !race

package deephash

import "reflect"

type pointer = unsafePointer

// pointerOf returns a pointer from v, which must be a reflect.Pointer.
func pointerOf(v reflect.Value) pointer { return unsafePointerOf(v) }
