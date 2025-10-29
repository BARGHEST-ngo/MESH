// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package dirfs

import (
	"context"
	"os"
)

// RemoveAll implements webdav.File. No removal is supported and this always
// returns os.ErrPermission.
func (dfs *FS) RemoveAll(ctx context.Context, name string) error {
	return &os.PathError{Op: "rm", Path: name, Err: os.ErrPermission}
}
