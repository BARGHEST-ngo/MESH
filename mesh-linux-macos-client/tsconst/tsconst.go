// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package tsconst exports some constants used elsewhere in the
// codebase.
package tsconst

// WintunInterfaceDesc is the description attached to Tailscale
// interfaces on Windows. This is set by the WinTun driver.
const WintunInterfaceDesc = "Tailscale Tunnel"
const WintunInterfaceDesc0_14 = "Wintun Userspace Tunnel"

// TailnetLockNotTrustedMsg is the error message used by network lock
// and sniffed (via substring) out of an error sent over the network.
const TailnetLockNotTrustedMsg = "this node is not trusted by network lock"
