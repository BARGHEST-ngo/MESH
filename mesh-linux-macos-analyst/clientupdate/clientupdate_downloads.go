// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build (linux && !android) || windows

package clientupdate

import (
	"context"

	"tailscale.com/clientupdate/distsign"
)

func (up *Updater) downloadURLToFile(pathSrc, fileDst string) (ret error) {
	c, err := distsign.NewClient(up.Logf, up.PkgsAddr)
	if err != nil {
		return err
	}
	return c.Download(context.Background(), pathSrc, fileDst)
}
