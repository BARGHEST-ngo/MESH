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
)

// RemoveAll implements webdav.File. No removal is supported and this always
// returns os.ErrPermission.
func (dfs *FS) RemoveAll(ctx context.Context, name string) error {
	return &os.PathError{Op: "rm", Path: name, Err: os.ErrPermission}
}
