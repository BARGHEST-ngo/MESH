// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build cgo && darwin && !ios

package posture

import (
	"fmt"
	"testing"

	"tailscale.com/types/logger"
	"tailscale.com/util/cibuild"
	"tailscale.com/util/syspolicy/policyclient"
)

func TestGetSerialNumberMac(t *testing.T) {
	// Do not run this test on CI, it can only be ran on macOS
	// and we currently only use Linux runners.
	if cibuild.On() {
		t.Skip()
	}

	sns, err := GetSerialNumbers(policyclient.NoPolicyClient{}, logger.Discard)
	if err != nil {
		t.Fatalf("failed to get serial number: %s", err)
	}

	if len(sns) != 1 {
		t.Errorf("expected list of one serial number, got %v", sns)
	}

	if len(sns[0]) <= 0 {
		t.Errorf("expected a serial number with more than zero characters, got %s", sns[0])
	}

	fmt.Printf("serials: %v\n", sns)
}
