// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build gokrazy

package dns

const (
	resolvConf = "/tmp/resolv.conf"
	backupConf = "/tmp/resolv.pre-tailscale-backup.conf"
)
