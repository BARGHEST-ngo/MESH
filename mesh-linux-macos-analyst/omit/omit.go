// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package omit provides consts to access Tailscale ts_omit_FOO build tags.
// They're often more convenient to eliminate some away locally with a const
// rather than using build tags.
package omit

import "errors"

// Err is an error that can be returned by functions in this package.
var Err = errors.New("feature not linked into binary per ts_omit build tag")
