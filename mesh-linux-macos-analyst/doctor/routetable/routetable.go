// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
