// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package packet contains packet parsing and marshaling utilities.
//
// Parsed provides allocation-free minimal packet header decoding, for
// use in packet filtering. The other types in the package are for
// constructing and marshaling packets into []bytes.
//
// To support allocation-free parsing, this package defines IPv4 and
// IPv6 address types. You should prefer to use netaddr's types,
// except where you absolutely need allocation-free IP handling
// (i.e. in the tunnel datapath) and are willing to implement all
// codepaths and data structures twice, once per IP family.
package packet
