// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build wasm || plan9 || tamago

package filch

import (
	"os"
)

func saveStderr() (*os.File, error) {
	return os.Stderr, nil
}

func unsaveStderr(f *os.File) error {
	os.Stderr = f
	return nil
}

func dup2Stderr(f *os.File) error {
	return nil
}
