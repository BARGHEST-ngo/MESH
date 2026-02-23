// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build darwin || freebsd || openbsd

package permissions

import (
	"golang.org/x/sys/unix"
	"tailscale.com/types/logger"
)

func permissionsImpl(logf logger.Logf) error {
	groups, _ := unix.Getgroups()
	logf("uid=%s euid=%s gid=%s egid=%s groups=%s",
		formatUserID(unix.Getuid()),
		formatUserID(unix.Geteuid()),
		formatGroupID(unix.Getgid()),
		formatGroupID(unix.Getegid()),
		formatGroups(groups),
	)
	return nil
}
