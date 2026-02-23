// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
