// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !(linux || darwin || freebsd || openbsd)

package permissions

import (
	"runtime"

	"tailscale.com/types/logger"
)

func permissionsImpl(logf logger.Logf) error {
	logf("unsupported on %s/%s", runtime.GOOS, runtime.GOARCH)
	return nil
}
