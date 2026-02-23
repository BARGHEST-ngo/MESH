// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package rands contains utility functions for randomness.
package rands

import (
	crand "crypto/rand"
	"encoding/hex"
)

// HexString returns a string of n cryptographically random lowercase
// hex characters.
//
// That is, HexString(3) returns something like "0fc", containing 12
// bits of randomness.
func HexString(n int) string {
	nb := n / 2
	if n%2 == 1 {
		nb++
	}
	b := make([]byte, nb)
	crand.Read(b)
	return hex.EncodeToString(b)[:n]
}
