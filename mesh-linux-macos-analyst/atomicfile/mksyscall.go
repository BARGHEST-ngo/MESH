// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package atomicfile

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go mksyscall.go

//sys replaceFileW(replaced *uint16, replacement *uint16, backup *uint16, flags uint32, exclude unsafe.Pointer, reserved unsafe.Pointer) (err error) [int32(failretval)==0] = kernel32.ReplaceFileW
