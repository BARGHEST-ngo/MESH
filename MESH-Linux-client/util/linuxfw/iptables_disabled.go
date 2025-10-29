// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
