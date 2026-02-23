// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !ios && !android && !ts_omit_hujsonconf

package conffile

import "github.com/tailscale/hujson"

// Only link the hujson package on platforms that use it, to reduce binary size
// & memory a bit.
//
// (iOS and Android don't have config files)

// While the linker's deadcode mostly handles the hujson package today, this
// keeps us honest for the future.

func init() {
	hujsonStandardize = hujson.Standardize
}
