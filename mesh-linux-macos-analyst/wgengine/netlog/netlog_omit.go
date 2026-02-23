// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ts_omit_netlog || ts_omit_logtail

package netlog

type Logger struct{}

func (*Logger) Startup(...any) error { return nil }
func (*Logger) Running() bool        { return false }
func (*Logger) Shutdown(any) error   { return nil }
func (*Logger) ReconfigRoutes(any)   {}
