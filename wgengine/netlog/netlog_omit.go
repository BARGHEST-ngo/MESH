// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ts_omit_netlog || ts_omit_logtail

package netlog

type Logger struct{}

func (*Logger) Startup(...any) error { return nil }
func (*Logger) Running() bool        { return false }
func (*Logger) Shutdown(any) error   { return nil }
func (*Logger) ReconfigRoutes(any)   {}
