// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
