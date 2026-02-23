// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build tailscale_go && android

package version

import "fmt"

func init() {
	// For official Android builds using the tailscale_go toolchain,
	// panic if the builder is screwed up and we fail to stamp a valid
	// version string.
	if !isValidLongWithTwoRepos(Long()) {
		panic(fmt.Sprintf("malformed version.Long value %q", Long()))
	}
}
