// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package dissector contains the Lua dissector for Tailscale packets.
package dissector

import (
	_ "embed"
)

//go:embed ts-dissector.lua
var Lua string
