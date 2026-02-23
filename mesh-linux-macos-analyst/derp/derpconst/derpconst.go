// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package derpconst contains constants used by the DERP client and server.
package derpconst

// MetaCertCommonNamePrefix is the prefix that the DERP server
// puts on for the common name of its "metacert". The suffix of
// the common name after "derpkey" is the hex key.NodePublic
// of the DERP server.
const MetaCertCommonNamePrefix = "derpkey"
