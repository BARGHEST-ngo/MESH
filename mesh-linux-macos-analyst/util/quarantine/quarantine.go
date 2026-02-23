// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package quarantine sets platform specific "quarantine" attributes on files
// that are received from other hosts.
package quarantine

import "os"

// SetOnFile sets the platform-specific quarantine attribute (if any) on the
// provided file.
func SetOnFile(f *os.File) error {
	return setQuarantineAttr(f)
}
