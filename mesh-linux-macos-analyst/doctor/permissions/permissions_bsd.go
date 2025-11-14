// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
