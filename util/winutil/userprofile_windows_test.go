// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package winutil

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestGetRoamingProfilePath(t *testing.T) {
	token := windows.GetCurrentProcessToken()
	computerName, userName, err := getComputerAndUserName(token, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := getRoamingProfilePath(t.Logf, token, computerName, userName); err != nil {
		t.Error(err)
	}

	// TODO(aaron): Flesh out better once can run tests under domain accounts.
}
