// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
