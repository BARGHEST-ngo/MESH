// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package distsign

import "testing"

func TestParseRoots(t *testing.T) {
	roots, err := parseRoots()
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) == 0 {
		t.Error("parseRoots returned no root keys")
	}
}
