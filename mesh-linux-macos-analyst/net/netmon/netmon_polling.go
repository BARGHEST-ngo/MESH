// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build (!linux && !freebsd && !windows && !darwin) || android

package netmon

import (
	"tailscale.com/types/logger"
	"tailscale.com/util/eventbus"
)

func newOSMon(_ *eventbus.Bus, logf logger.Logf, m *Monitor) (osMon, error) {
	return newPollingMon(logf, m)
}

// unspecifiedMessage is a minimal message implementation that should not
// be ignored. In general, OS-specific implementations should use better
// types and avoid this if they can.
type unspecifiedMessage struct{}

func (unspecifiedMessage) ignore() bool { return false }
