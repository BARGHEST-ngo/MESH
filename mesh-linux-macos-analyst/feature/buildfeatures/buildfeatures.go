// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:generate go run gen.go

// The buildfeatures package contains boolean constants indicating which
// features were included in the binary (via build tags), for use in dead code
// elimination when using separate build tag protected files is impractical
// or undesirable.
package buildfeatures
