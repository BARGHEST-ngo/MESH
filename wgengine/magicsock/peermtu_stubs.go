// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build (!linux && !darwin) || android || ios

package magicsock

import "tailscale.com/disco"

func (c *Conn) DontFragSetting() (bool, error) {
	return false, nil
}

func (c *Conn) ShouldPMTUD() bool {
	return false
}

func (c *Conn) PeerMTUEnabled() bool {
	return false
}

func (c *Conn) UpdatePMTUD() {
}

func pmtuShouldLogDiscoTxErr(m disco.Message, err error) bool {
	return true
}
