// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ios || android || ts_omit_debugeventbus

package eventbus

type tswebDebugHandler = any // actually *tsweb.DebugHandler; any to avoid import tsweb with ts_omit_debugeventbus

func (*Debugger) RegisterHTTP(td tswebDebugHandler) {}
