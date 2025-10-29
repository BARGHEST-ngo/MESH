// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
