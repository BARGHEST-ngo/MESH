// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
