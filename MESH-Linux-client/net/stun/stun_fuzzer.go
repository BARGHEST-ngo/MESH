// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.
//go:build gofuzz

package stun

func FuzzStunParser(data []byte) int {
	_, _, _ = ParseResponse(data)

	_, _ = ParseBindingRequest(data)
	return 1
}
