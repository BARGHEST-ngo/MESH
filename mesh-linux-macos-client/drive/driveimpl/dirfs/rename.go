// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package dirfs

import (
	"context"
	"os"
)

// Rename implements interface webdav.FileSystem. No renaming is supported and
// this always returns os.ErrPermission.
func (dfs *FS) Rename(ctx context.Context, oldName, newName string) error {
	return &os.PathError{Op: "mv", Path: oldName, Err: os.ErrPermission}
}
