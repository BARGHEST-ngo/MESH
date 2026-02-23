// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build darwin

package hostinfo

import (
	"os"
	"path/filepath"
)

func init() {
	packageType = packageTypeDarwin
}

func packageTypeDarwin() string {
	// Using tailscaled or IPNExtension?
	exe, _ := os.Executable()
	return filepath.Base(exe)
}
