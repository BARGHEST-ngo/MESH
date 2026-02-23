// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !linux || android

package netkernelconf

// CheckUDPGROForwarding is unimplemented for non-Linux platforms. Refer to the
// docstring in _linux.go.
func CheckUDPGROForwarding(tunInterface, defaultRouteInterface string) (warn, err error) {
	return nil, nil
}

// SetUDPGROForwarding is unimplemented for non-Linux platforms. Refer to the
// docstring in _linux.go.
func SetUDPGROForwarding(tunInterface, defaultRouteInterface string) error {
	return nil
}
