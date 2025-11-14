// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package routetable provides a doctor.Check that dumps the current system's
// route table to the log.
package routetable

import (
	"context"

	"tailscale.com/net/routetable"
	"tailscale.com/types/logger"
)

// MaxRoutes is the maximum number of routes that will be displayed.
const MaxRoutes = 1000

// Check implements the doctor.Check interface.
type Check struct{}

func (Check) Name() string {
	return "routetable"
}

func (Check) Run(_ context.Context, logf logger.Logf) error {
	rs, err := routetable.Get(MaxRoutes)
	if err != nil {
		return err
	}
	for _, r := range rs {
		logf("%s", r)
	}
	return nil
}
