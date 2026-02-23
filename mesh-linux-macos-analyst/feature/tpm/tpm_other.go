// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !linux && !windows

package tpm

import (
	"errors"

	"github.com/google/go-tpm/tpm2/transport"
)

func open() (transport.TPMCloser, error) {
	return nil, errors.New("TPM not supported on this platform")
}
