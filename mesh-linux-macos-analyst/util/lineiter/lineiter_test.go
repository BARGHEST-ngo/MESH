// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package lineiter

import (
	"slices"
	"strings"
	"testing"
)

func TestBytesLines(t *testing.T) {
	var got []string
	for line := range Bytes([]byte("foo\n\nbar\nbaz")) {
		got = append(got, string(line))
	}
	want := []string{"foo", "", "bar", "baz"}
	if !slices.Equal(got, want) {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestReader(t *testing.T) {
	var got []string
	for line := range Reader(strings.NewReader("foo\n\nbar\nbaz")) {
		got = append(got, string(line.MustValue()))
	}
	want := []string{"foo", "", "bar", "baz"}
	if !slices.Equal(got, want) {
		t.Errorf("got %q; want %q", got, want)
	}
}
