// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !js && !wasm && !plan9

package magicsock

import (
	"errors"
	"syscall"
)

// errHOSTUNREACH wraps unix.EHOSTUNREACH in an interface type to pass to
// errors.Is while avoiding an allocation per call.
var errHOSTUNREACH error = syscall.EHOSTUNREACH

// isBadEndpointErr checks if err is one which is known to report that an
// endpoint can no longer be sent to. It is not exhaustive, and for unknown
// errors always reports false.
func isBadEndpointErr(err error) bool {
	return errors.Is(err, errHOSTUNREACH)
}
