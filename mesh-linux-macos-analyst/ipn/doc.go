// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:generate go run tailscale.com/cmd/viewer -type=LoginProfile,Prefs,ServeConfig,ServiceConfig,TCPPortHandler,HTTPHandler,WebServerConfig

// Package ipn implements the interactions between the Tailscale cloud
// control plane and the local network stack.
//
// IPN is the abbreviated name for a Tailscale network. What's less
// clear is what it's an abbreviation for: Identified Private Network?
// IP Network? Internet Private Network? I Privately Network?
package ipn
