// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package sysresources

// TotalMemory returns the total accessible system memory, in bytes. If the
// value cannot be determined, then 0 will be returned.
func TotalMemory() uint64 {
	return totalMemoryImpl()
}
