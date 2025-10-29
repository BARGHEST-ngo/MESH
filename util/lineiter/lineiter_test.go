// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
