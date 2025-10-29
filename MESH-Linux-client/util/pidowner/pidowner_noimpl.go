// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !windows && !linux

package pidowner

func ownerOfPID(pid int) (userID string, err error) { return "", ErrNotImplemented }
