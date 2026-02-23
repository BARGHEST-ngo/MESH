// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build ios || android || ts_omit_debugeventbus

package eventbus

type tswebDebugHandler = any // actually *tsweb.DebugHandler; any to avoid import tsweb with ts_omit_debugeventbus

func (*Debugger) RegisterHTTP(td tswebDebugHandler) {}
