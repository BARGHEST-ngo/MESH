// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package main

//go:generate go run tailscale.com/cmd/mkmanifest amd64 windows-manifest.xml manifest_windows_amd64.syso
//go:generate go run tailscale.com/cmd/mkmanifest 386 windows-manifest.xml manifest_windows_386.syso
//go:generate go run tailscale.com/cmd/mkmanifest arm64 windows-manifest.xml manifest_windows_arm64.syso
