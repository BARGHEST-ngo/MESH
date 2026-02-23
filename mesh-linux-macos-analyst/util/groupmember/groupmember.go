// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package groupmember verifies group membership of the provided user on the
// local system.
package groupmember

import (
	"os/user"
	"slices"
)

// IsMemberOfGroup reports whether the provided user is a member of
// the provided system group.
func IsMemberOfGroup(group, userName string) (bool, error) {
	u, err := user.Lookup(userName)
	if err != nil {
		return false, err
	}
	g, err := user.LookupGroup(group)
	if err != nil {
		return false, err
	}
	ugids, err := u.GroupIds()
	if err != nil {
		return false, err
	}
	return slices.Contains(ugids, g.Gid), nil
}
