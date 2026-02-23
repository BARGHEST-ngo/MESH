// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build android || (!linux && !darwin && !freebsd)

package routetable

import (
	"errors"
	"runtime"
)

var errUnsupported = errors.New("cannot get route table on platform " + runtime.GOOS)

func Get(max int) ([]RouteEntry, error) {
	return nil, errUnsupported
}
