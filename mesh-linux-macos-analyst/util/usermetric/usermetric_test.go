// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
