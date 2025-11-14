// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
