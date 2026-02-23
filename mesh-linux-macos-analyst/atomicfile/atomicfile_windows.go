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

import (
	"os"

	"golang.org/x/sys/windows"
)

func rename(srcFile, destFile string) error {
	// Use replaceFile when possible to preserve the original file's attributes and ACLs.
	if err := replaceFile(destFile, srcFile); err == nil || err != windows.ERROR_FILE_NOT_FOUND {
		return err
	}
	// destFile doesn't exist. Just do a normal rename.
	return os.Rename(srcFile, destFile)
}

func replaceFile(destFile, srcFile string) error {
	destFile16, err := windows.UTF16PtrFromString(destFile)
	if err != nil {
		return err
	}

	srcFile16, err := windows.UTF16PtrFromString(srcFile)
	if err != nil {
		return err
	}

	return replaceFileW(destFile16, srcFile16, nil, 0, nil, nil)
}
