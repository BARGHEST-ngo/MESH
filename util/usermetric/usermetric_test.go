// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package usermetric

import (
	"bytes"
	"testing"
)

func TestGauge(t *testing.T) {
	var reg Registry
	g := reg.NewGauge("test_gauge", "This is a test gauge")
	g.Set(15)

	var buf bytes.Buffer
	g.WritePrometheus(&buf, "test_gauge")
	const want = `# TYPE test_gauge gauge
# HELP test_gauge This is a test gauge
test_gauge 15
`
	if got := buf.String(); got != want {
		t.Errorf("got %q; want %q", got, want)
	}

}
