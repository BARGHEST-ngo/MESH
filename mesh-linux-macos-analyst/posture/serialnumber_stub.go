// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
