// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package main

//go:generate go run tailscale.com/cmd/mkmanifest amd64 windows-manifest.xml manifest_windows_amd64.syso
//go:generate go run tailscale.com/cmd/mkmanifest 386 windows-manifest.xml manifest_windows_386.syso
//go:generate go run tailscale.com/cmd/mkmanifest arm64 windows-manifest.xml manifest_windows_arm64.syso
