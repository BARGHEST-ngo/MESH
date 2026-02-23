// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !windows

package netstat

// OSMetadata includes any additional OS-specific information that may be
// obtained during the retrieval of a given Entry.
type OSMetadata struct{}

func get() (*Table, error) {
	return nil, ErrNotImplemented
}
