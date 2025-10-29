/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2025 WireGuard LLC. All Rights Reserved.
 * Additional types added for Tailscale compatibility.
 */

package conn

import (
	"net/netip"
)

// InitiationAwareEndpoint is an Endpoint that can provide information about
// the public key used in a WireGuard handshake initiation message.
type InitiationAwareEndpoint interface {
	Endpoint
	// InitiationMessagePublicKey returns the public key from the most recent
	// handshake initiation message received from this endpoint.
	InitiationMessagePublicKey() (key [32]byte, ok bool)
}

// PeerAwareEndpoint is an Endpoint that can provide information about
// the peer's IP address.
type PeerAwareEndpoint interface {
	Endpoint
	// PeerAddr returns the peer's IP address and port.
	PeerAddr() netip.AddrPort
}

