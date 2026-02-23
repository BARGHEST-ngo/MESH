// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package dirfs

import (
	"context"
	"os"

	"tailscale.com/drive/driveimpl/shared"
)

// Mkdir implements webdav.FileSystem. All attempts to Mkdir a directory that
// already exists will succeed. All other attempts will fail with
// os.ErrPermission.
func (dfs *FS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	nameWithoutStaticRoot, isStaticRoot := dfs.trimStaticRoot(name)
	if isStaticRoot || shared.IsRoot(name) {
		// root directory already exists, consider this okay
		return nil
	}

	child := dfs.childFor(nameWithoutStaticRoot)
	if child != nil {
		// child already exists, consider this okay
		return nil
	}

	return &os.PathError{Op: "mkdir", Path: name, Err: os.ErrPermission}
}
