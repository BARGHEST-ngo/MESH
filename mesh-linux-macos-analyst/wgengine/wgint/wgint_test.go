// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package wgint

import (
	"testing"

	"github.com/amnezia-vpn/amneziawg-go/device"
)

func TestInternalOffsets(t *testing.T) {
	peer := new(device.Peer)
	if got := peerLastHandshakeNano(peer); got != 0 {
		t.Errorf("PeerLastHandshakeNano = %v, want 0", got)
	}
	if got := peerRxBytes(peer); got != 0 {
		t.Errorf("PeerRxBytes = %v, want 0", got)
	}
	if got := peerTxBytes(peer); got != 0 {
		t.Errorf("PeerTxBytes = %v, want 0", got)
	}
	if got := peerHandshakeAttempts(peer); got != 0 {
		t.Errorf("PeerHandshakeAttempts = %v, want 0", got)
	}
}
