// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !windows

package sockopts

import (
	"tailscale.com/types/nettype"
)

// SetICMPErrImmunity is no-op on non-Windows.
func SetICMPErrImmunity(pconn nettype.PacketConn) error {
	return nil
}
