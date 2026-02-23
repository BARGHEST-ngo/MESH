// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ios || ts_omit_gro

package gro

import (
	"runtime"

	"tailscale.com/net/packet"
)

type GRO struct{}

func NewGRO() *GRO {
	if runtime.GOOS == "ios" {
		panic("unsupported on iOS")
	}
	panic("GRO disabled in build")

}

func (g *GRO) SetDispatcher(any) {}

func (g *GRO) Enqueue(_ *packet.Parsed) {}

func (g *GRO) Flush() {}
