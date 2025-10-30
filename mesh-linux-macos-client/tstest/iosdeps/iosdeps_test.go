// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package iosdeps

import (
	"testing"

	"tailscale.com/tstest/deptest"
)

func TestDeps(t *testing.T) {
	deptest.DepChecker{
		GOOS:   "ios",
		GOARCH: "arm64",
		BadDeps: map[string]string{
			"testing":                             "do not use testing package in production code",
			"text/template":                       "linker bloat (MethodByName)",
			"html/template":                       "linker bloat (MethodByName)",
			"tailscale.com/net/wsconn":            "https://github.com/tailscale/tailscale/issues/13762",
			"github.com/coder/websocket":          "https://github.com/tailscale/tailscale/issues/13762",
			"github.com/mitchellh/go-ps":          "https://github.com/tailscale/tailscale/pull/13759",
			"database/sql/driver":                 "iOS doesn't use an SQL database",
			"github.com/google/uuid":              "see tailscale/tailscale#13760",
			"tailscale.com/clientupdate/distsign": "downloads via AppStore, not distsign",
			"github.com/tailscale/hujson":         "no config file support on iOS",
			"tailscale.com/feature/capture":       "no debug packet capture on iOS",
		},
	}.Check(t)
}
