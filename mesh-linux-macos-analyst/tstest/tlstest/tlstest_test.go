// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
