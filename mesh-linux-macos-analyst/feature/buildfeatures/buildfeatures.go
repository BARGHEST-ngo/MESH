// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:generate go run gen.go

// The buildfeatures package contains boolean constants indicating which
// features were included in the binary (via build tags), for use in dead code
// elimination when using separate build tag protected files is impractical
// or undesirable.
package buildfeatures
