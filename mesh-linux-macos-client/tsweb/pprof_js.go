// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build js && wasm

package tsweb

func addProfilingHandlers(d *DebugHandler) {
	// No pprof in js builds, pprof doesn't work and bloats the build.
}
