// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
