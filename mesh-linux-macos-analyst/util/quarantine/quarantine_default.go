// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !darwin && !windows

package quarantine

import (
	"os"
)

func setQuarantineAttr(f *os.File) error {
	return nil
}
