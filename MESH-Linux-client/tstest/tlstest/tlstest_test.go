// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package tlstest

import (
	"testing"
)

func TestPrivateKey(t *testing.T) {
	a := privateKey("a.tstest")
	a2 := privateKey("a.tstest")
	b := privateKey("b.tstest")

	if string(a) != string(a2) {
		t.Errorf("a and a2 should be equal")
	}
	if string(a) == string(b) {
		t.Errorf("a and b should not be equal")
	}
}
