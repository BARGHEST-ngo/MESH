// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build darwin

package routetable

import "golang.org/x/sys/unix"

const (
	ribType        = unix.NET_RT_DUMP2
	parseType      = unix.NET_RT_IFLIST2
	rmExpectedType = unix.RTM_GET2

	// Skip routes that were cloned from a parent
	skipFlags = unix.RTF_WASCLONED
)

var flags = map[int]string{
	unix.RTF_BLACKHOLE: "blackhole",
	unix.RTF_BROADCAST: "broadcast",
	unix.RTF_GATEWAY:   "gateway",
	unix.RTF_GLOBAL:    "global",
	unix.RTF_HOST:      "host",
	unix.RTF_IFSCOPE:   "ifscope",
	unix.RTF_LOCAL:     "local",
	unix.RTF_MULTICAST: "multicast",
	unix.RTF_REJECT:    "reject",
	unix.RTF_ROUTER:    "router",
	unix.RTF_STATIC:    "static",
	unix.RTF_UP:        "up",
	// More obscure flags, just to have full coverage.
	unix.RTF_LLINFO:    "{RTF_LLINFO}",
	unix.RTF_PRCLONING: "{RTF_PRCLONING}",
	unix.RTF_CLONING:   "{RTF_CLONING}",
}
