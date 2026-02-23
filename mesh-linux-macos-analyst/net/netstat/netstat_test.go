// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package netstat

import (
	"testing"
)

func TestGet(t *testing.T) {
	nt, err := Get()
	if err == ErrNotImplemented {
		t.Skip("TODO: not implemented")
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range nt.Entries {
		t.Logf("Entry: %+v", e)
	}
}
