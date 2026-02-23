// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
