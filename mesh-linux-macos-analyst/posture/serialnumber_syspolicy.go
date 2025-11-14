// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build android || ios

package posture

import (
	"fmt"

	"tailscale.com/types/logger"
	"tailscale.com/util/syspolicy/pkey"
	"tailscale.com/util/syspolicy/policyclient"
)

// GetSerialNumbers returns the serial number of the device as reported by an
// MDM solution. It requires configuration via the DeviceSerialNumber system policy.
// This is the only way to gather serial numbers on iOS, tvOS and Android.
func GetSerialNumbers(polc policyclient.Client, _ logger.Logf) ([]string, error) {
	s, err := polc.GetString(pkey.DeviceSerialNumber, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get serial number from MDM: %v", err)
	}
	if s != "" {
		return []string{s}, nil
	}
	return nil, nil
}
