// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
