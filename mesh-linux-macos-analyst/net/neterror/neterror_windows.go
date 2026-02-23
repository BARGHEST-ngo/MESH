// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package neterror

import (
	"errors"

	"golang.org/x/sys/windows"
)

func init() {
	packetWasTruncated = func(err error) bool {
		return errors.Is(err, windows.WSAEMSGSIZE)
	}
}
