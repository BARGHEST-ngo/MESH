// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !windows

package tstun

import "github.com/amnezia-vpn/amneziawg-go/tun"

func interfaceName(dev tun.Device) (string, error) {
	return dev.Name()
}
