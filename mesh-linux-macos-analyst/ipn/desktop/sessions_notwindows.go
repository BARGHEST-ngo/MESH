// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !windows

package desktop

import "tailscale.com/types/logger"

// NewSessionManager returns a new [SessionManager] for the current platform,
// [ErrNotImplemented] if the platform is not supported, or an error if the
// session manager could not be created.
func NewSessionManager(logger.Logf) (SessionManager, error) {
	return nil, ErrNotImplemented
}
