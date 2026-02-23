// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package cibuild reports runtime CI information.
package cibuild

import "os"

// On reports whether the current binary is executing on a CI system.
func On() bool {
	// CI env variable is set by GitHub.
	// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	return os.Getenv("GITHUB_ACTIONS") != "" || os.Getenv("CI") == "true"
}
