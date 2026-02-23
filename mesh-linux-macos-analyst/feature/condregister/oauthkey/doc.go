// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package oauthkey registers support for OAuth key resolution
// if it's not disabled via the ts_omit_oauthkey build tag.
// Currently (2025-09-19), tailscaled does not need OAuth key
// resolution, only the CLI and tsnet do, so this package is
// pulled out separately to avoid linking OAuth packages into
// tailscaled.
package oauthkey
