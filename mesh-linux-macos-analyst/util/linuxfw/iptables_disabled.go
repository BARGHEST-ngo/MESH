// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build (linux && !(arm64 || amd64)) || ts_omit_iptables

package linuxfw

import (
	"errors"

	"tailscale.com/types/logger"
)

func detectIptables() (int, error) {
	return 0, nil
}

func newIPTablesRunner(logf logger.Logf) (*iptablesRunner, error) {
	return nil, errors.New("iptables disabled in build")
}
