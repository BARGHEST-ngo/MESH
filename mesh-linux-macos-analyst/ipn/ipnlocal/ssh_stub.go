// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_omit_ssh || ios || android || (!linux && !darwin && !freebsd && !openbsd && !plan9)

package ipnlocal

import (
	"errors"

	"tailscale.com/tailcfg"
)

func (b *LocalBackend) getSSHHostKeyPublicStrings() ([]string, error) {
	return nil, nil
}

func (b *LocalBackend) getSSHUsernames(*tailcfg.C2NSSHUsernamesRequest) (*tailcfg.C2NSSHUsernamesResponse, error) {
	return nil, errors.New("not implemented")
}
