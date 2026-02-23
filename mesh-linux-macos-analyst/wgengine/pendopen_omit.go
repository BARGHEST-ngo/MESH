// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_omit_debug

package wgengine

import (
	"tailscale.com/net/packet"
	"tailscale.com/net/tstun"
	"tailscale.com/wgengine/filter"
)

type flowtrackTuple = struct{}

type pendingOpenFlow struct{}

func (*userspaceEngine) trackOpenPreFilterIn(pp *packet.Parsed, t *tstun.Wrapper) (res filter.Response) {
	panic("unreachable")
}

func (*userspaceEngine) trackOpenPostFilterOut(pp *packet.Parsed, t *tstun.Wrapper) (res filter.Response) {
	panic("unreachable")
}
