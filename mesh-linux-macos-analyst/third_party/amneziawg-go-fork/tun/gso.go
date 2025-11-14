// Copyright (c) 2024 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package tun

// GSOType is the type of generic segmentation offload.
type GSOType uint8

const (
	GSONone GSOType = iota
	GSOTCPv4
	GSOTCPv6
	GSOGvisor
)

// GSOOptions contains generic segmentation offload options for a packet.
type GSOOptions struct {
	GSOType    GSOType
	GSOSize    uint16
	CsumStart  uint16
	CsumOffset uint16
	NeedsCsum  bool
	HdrLen     uint16
}

// GSOSplit splits a packet into smaller packets for GSO.
// This is a compatibility shim for amneziawg-go which doesn't have GSO support.
// The actual implementation would split the packet based on GSOSize.
// Returns the number of packets and any error.
func GSOSplit(pkt []byte, opts GSOOptions, outBuffs [][]byte, sizes []int, offset int) (int, error) {
	if opts.GSOType == GSONone || opts.GSOSize == 0 {
		if len(outBuffs) > 0 {
			outBuffs[0] = pkt
			sizes[0] = len(pkt)
			return 1, nil
		}
		return 0, nil
	}
	// For now, return the packet as-is since amneziawg-go doesn't support GSO
	if len(outBuffs) > 0 {
		outBuffs[0] = pkt
		sizes[0] = len(pkt)
		return 1, nil
	}
	return 0, nil
}

// GRODevice is a generic receive offload device interface.
// This is a compatibility shim for amneziawg-go which doesn't have GRO support.
type GRODevice interface {
	// DisableUDPGRO disables UDP GRO.
	DisableUDPGRO()
	// DisableTCPGRO disables TCP GRO.
	DisableTCPGRO()
	// Write writes packets to the device.
	Write(bufs [][]byte, offset int) (int, error)
	// Enqueue enqueues a packet for GRO processing.
	Enqueue(pkt interface{}) error
	// Flush flushes any pending packets.
	Flush() error
	// Close closes the GRO device.
	Close() error
}

