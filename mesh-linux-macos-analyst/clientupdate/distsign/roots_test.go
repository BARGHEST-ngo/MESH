// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
