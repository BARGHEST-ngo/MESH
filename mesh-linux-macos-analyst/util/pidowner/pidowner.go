// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package pidowner handles lookups from process ID to its owning user.
package pidowner

import (
	"errors"
	"runtime"
)

var ErrNotImplemented = errors.New("not implemented for GOOS=" + runtime.GOOS)

var ErrProcessNotFound = errors.New("process not found")

// OwnerOfPID returns the user ID that owns the given process ID.
//
// The returned user ID is suitable to passing to os/user.LookupId.
//
// The returned error will be ErrNotImplemented for operating systems where
// this isn't supported.
func OwnerOfPID(pid int) (userID string, err error) {
	return ownerOfPID(pid)
}
