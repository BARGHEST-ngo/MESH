// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// js: not implemented
// plan9: not implemented
// solaris: currently unsupported by go-smbios:
// https://github.com/digitalocean/go-smbios/pull/21

//go:build solaris || plan9 || js || wasm || tamago || aix || (darwin && !cgo && !ios)

package posture

import (
	"errors"

	"tailscale.com/types/logger"
	"tailscale.com/util/syspolicy/policyclient"
)

// GetSerialNumber returns client machine serial number(s).
func GetSerialNumbers(polc policyclient.Client, _ logger.Logf) ([]string, error) {
	return nil, errors.New("not implemented")
}
