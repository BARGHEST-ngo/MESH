// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package wgcfg has types and a parser for representing WireGuard config.
package wgcfg

import (
	"net/netip"
	"slices"

	"tailscale.com/tailcfg"
	"tailscale.com/types/key"
	"tailscale.com/types/logid"
)

//go:generate go run tailscale.com/cmd/cloner -type=Config,Peer

// Config is a WireGuard configuration.
// It only supports the set of things Tailscale uses.
type Config struct {
	Name       string
	NodeID     tailcfg.StableNodeID
	PrivateKey key.NodePrivate
	Addresses  []netip.Prefix
	MTU        uint16
	DNS        []netip.Addr
	Peers      []Peer

	Jc   uint8
	Jmin uint16
	Jmax uint16
	S1   uint16
	S2   uint16
	H1   uint32
	H2   uint32
	H3   uint32
	H4   uint32

	NetworkLogging struct {
		NodeID             logid.PrivateID
		DomainID           logid.PrivateID
		LogExitFlowEnabled bool
	}
}

func (c *Config) Equal(o *Config) bool {
	if c == nil || o == nil {
		return c == o
	}
	return c.Name == o.Name &&
		c.NodeID == o.NodeID &&
		c.PrivateKey.Equal(o.PrivateKey) &&
		c.MTU == o.MTU &&
		c.Jc == o.Jc &&
		c.Jmin == o.Jmin &&
		c.Jmax == o.Jmax &&
		c.S1 == o.S1 &&
		c.S2 == o.S2 &&
		c.H1 == o.H1 &&
		c.H2 == o.H2 &&
		c.H3 == o.H3 &&
		c.H4 == o.H4 &&
		c.NetworkLogging == o.NetworkLogging &&
		slices.Equal(c.Addresses, o.Addresses) &&
		slices.Equal(c.DNS, o.DNS) &&
		slices.EqualFunc(c.Peers, o.Peers, Peer.Equal)
}

type Peer struct {
	PublicKey           key.NodePublic
	DiscoKey            key.DiscoPublic // present only so we can handle restarts within wgengine, not passed to WireGuard
	AllowedIPs          []netip.Prefix
	V4MasqAddr          *netip.Addr // if non-nil, masquerade IPv4 traffic to this peer using this address
	V6MasqAddr          *netip.Addr // if non-nil, masquerade IPv6 traffic to this peer using this address
	IsJailed            bool        // if true, this peer is jailed and cannot initiate connections
	PersistentKeepalive uint16      // in seconds between keep-alives; 0 to disable
	// wireguard-go's endpoint for this peer. It should always equal Peer.PublicKey.
	// We represent it explicitly so that we can detect if they diverge and recover.
	// There is no need to set WGEndpoint explicitly when constructing a Peer by hand.
	// It is only populated when reading Peers from wireguard-go.
	WGEndpoint key.NodePublic
}

func addrPtrEq(a, b *netip.Addr) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (p Peer) Equal(o Peer) bool {
	return p.PublicKey == o.PublicKey &&
		p.DiscoKey == o.DiscoKey &&
		slices.Equal(p.AllowedIPs, o.AllowedIPs) &&
		p.IsJailed == o.IsJailed &&
		p.PersistentKeepalive == o.PersistentKeepalive &&
		addrPtrEq(p.V4MasqAddr, o.V4MasqAddr) &&
		addrPtrEq(p.V6MasqAddr, o.V6MasqAddr) &&
		p.WGEndpoint == o.WGEndpoint
}

// PeerWithKey returns the Peer with key k and reports whether it was found.
func (config Config) PeerWithKey(k key.NodePublic) (Peer, bool) {
	for _, p := range config.Peers {
		if p.PublicKey == k {
			return p, true
		}
	}
	return Peer{}, false
}
