// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !race

package version

// IsRace reports whether the current binary was built with the Go
// race detector enabled.
func IsRace() bool { return false }
