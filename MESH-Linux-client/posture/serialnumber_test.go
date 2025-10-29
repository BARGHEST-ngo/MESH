// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package posture

import (
	"testing"

	"tailscale.com/types/logger"
	"tailscale.com/util/syspolicy/policyclient"
)

func TestGetSerialNumber(t *testing.T) {
	// ensure GetSerialNumbers is implemented
	// or covered by a stub on a given platform.
	_, _ = GetSerialNumbers(policyclient.NoPolicyClient{}, logger.Discard)
}
