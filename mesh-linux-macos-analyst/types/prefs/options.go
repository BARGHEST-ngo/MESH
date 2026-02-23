// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package prefs

// Options are used to configure additional parameters of a preference.
type Options func(s *metadata)

var (
	// ReadOnly is an option that marks preference as read-only.
	ReadOnly Options = markReadOnly
	// Managed is an option that marks preference as managed.
	Managed Options = markManaged
)

func markReadOnly(s *metadata) {
	s.ReadOnly = true
}

func markManaged(s *metadata) {
	s.Managed = true
}
