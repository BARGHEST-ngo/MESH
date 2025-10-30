// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package healthmsg contains some constants for health messages.
//
// It's a leaf so both the server and CLI can depend on it without bringing too
// much in to the CLI binary.
package healthmsg

const (
	WarnAcceptRoutesOff = "Some peers are advertising routes but --accept-routes is false"
	TailscaleSSHOnBut   = "Tailscale SSH enabled, but " // + ... something from caller
	LockedOut           = "this node is locked out; it will not have connectivity until it is signed. For more info, see https://tailscale.com/s/locked-out"
	WarnExitNodeUsage   = "The following issues on your machine will likely make usage of exit nodes impossible"
	DisableRPFilter     = "Please set rp_filter=2 instead of rp_filter=1; see https://github.com/tailscale/tailscale/issues/3310"
)
