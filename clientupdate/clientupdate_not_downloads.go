// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !((linux && !android) || windows)

package clientupdate

func (up *Updater) downloadURLToFile(pathSrc, fileDst string) (ret error) {
	panic("unreachable")
}
