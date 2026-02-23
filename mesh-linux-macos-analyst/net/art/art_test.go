// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package art

import (
	"os"
	"testing"

	"tailscale.com/util/cibuild"
)

func TestMain(m *testing.M) {
	if cibuild.On() {
		// Skip CI on GitHub for now
		// TODO: https://github.com/tailscale/tailscale/issues/7866
		os.Exit(0)
	}
	os.Exit(m.Run())
}
