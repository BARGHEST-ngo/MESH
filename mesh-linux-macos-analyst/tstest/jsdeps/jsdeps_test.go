// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package jsdeps

import (
	"testing"

	"tailscale.com/tstest/deptest"
)

func TestDeps(t *testing.T) {
	deptest.DepChecker{
		GOOS:   "js",
		GOARCH: "wasm",
		BadDeps: map[string]string{
			"testing":                     "do not use testing package in production code",
			"runtime/pprof":               "bloat",
			"golang.org/x/net/http2/h2c":  "bloat",
			"net/http/pprof":              "bloat",
			"golang.org/x/net/proxy":      "bloat",
			"github.com/tailscale/goupnp": "bloat, which can't work anyway in wasm",
		},
	}.Check(t)
}
