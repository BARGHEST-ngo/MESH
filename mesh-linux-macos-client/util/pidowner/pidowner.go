// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
