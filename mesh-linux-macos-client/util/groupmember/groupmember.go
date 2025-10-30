// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
