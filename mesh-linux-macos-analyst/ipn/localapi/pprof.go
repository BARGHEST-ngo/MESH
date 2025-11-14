// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !ios && !android && !js && !ts_omit_debug

// We don't include it on mobile where we're more memory constrained and
// there's no CLI to get at the results anyway.

package localapi

import (
	"net/http"
	"net/http/pprof"
)

func init() {
	servePprofFunc = servePprof
}

func servePprof(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	switch name {
	case "profile":
		pprof.Profile(w, r)
	default:
		pprof.Handler(name).ServeHTTP(w, r)
	}
}
