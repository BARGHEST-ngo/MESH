// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ios || android || js

package localapi

import (
	"net/http"
	"runtime"
)

func (h *Handler) serveCert(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "disabled on "+runtime.GOOS, http.StatusNotFound)
}
