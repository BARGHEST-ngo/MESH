// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !(ts_omit_debug && ts_omit_clientmetrics && ts_omit_usermetrics)

// Package expvar contains type aliases for expvar types, to allow conditionally
// excluding the package from builds.
package expvar

import "expvar"

type Int = expvar.Int
