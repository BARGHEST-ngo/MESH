// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ios

package version

import (
	"os"
)

func CmdName() string {
	e, err := os.Executable()
	if err != nil {
		return "cmd"
	}
	return e
}
