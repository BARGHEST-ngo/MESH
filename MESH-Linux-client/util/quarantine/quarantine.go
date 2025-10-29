// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package quarantine sets platform specific "quarantine" attributes on files
// that are received from other hosts.
package quarantine

import "os"

// SetOnFile sets the platform-specific quarantine attribute (if any) on the
// provided file.
func SetOnFile(f *os.File) error {
	return setQuarantineAttr(f)
}
