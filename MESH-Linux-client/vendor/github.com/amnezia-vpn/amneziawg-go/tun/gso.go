/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2025 WireGuard LLC. All Rights Reserved.
 * Additional GSO/GRO types added for Tailscale compatibility.
 */

package tun

import (
	"errors"
)

// GSOType is the type of GSO (Generic Segmentation Offload) to perform.
type GSOType uint8

const (
	// GSONone means no GSO.
	GSONone GSOType = iota
	// GSOTCPv4 means TCPv4 GSO.
	GSOTCPv4
	// GSOTCPv6 means TCPv6 GSO.
	GSOTCPv6
)

// GSOOptions contains options for GSO (Generic Segmentation Offload).
type GSOOptions struct {
	// GSOType is the type of GSO to perform.
	GSOType GSOType
	// HdrLen is the length of the headers (L3 + L4).
	HdrLen uint16
	// GSOSize is the maximum segment size.
	GSOSize uint16
	// CsumStart is the offset to start checksumming from.
	CsumStart uint16
	// CsumOffset is the offset within the checksum field.
	CsumOffset uint16
	// NeedsCsum indicates whether checksumming is needed.
	NeedsCsum bool
}

// GSOSplit splits a GSO packet into multiple packets.
// It takes a packet with GSO options and splits it into outBuffs,
// setting the sizes in the sizes slice.
// Returns the number of packets created.
func GSOSplit(pkt []byte, options GSOOptions, outBuffs [][]byte, sizes []int, offset int) (int, error) {
	if options.GSOType == GSONone {
		// No GSO, just copy the packet
		if len(outBuffs) < 1 {
			return 0, errors.New("insufficient output buffers")
		}
		n := copy(outBuffs[0][offset:], pkt)
		sizes[0] = n
		return 1, nil
	}

	// For GSO packets, we need to split them
	// This is a simplified implementation
	hdrLen := int(options.HdrLen)
	gsoSize := int(options.GSOSize)

	if hdrLen > len(pkt) {
		return 0, errors.New("header length exceeds packet length")
	}

	payload := pkt[hdrLen:]
	payloadLen := len(payload)

	if payloadLen == 0 {
		// No payload, just copy headers
		if len(outBuffs) < 1 {
			return 0, errors.New("insufficient output buffers")
		}
		n := copy(outBuffs[0][offset:], pkt[:hdrLen])
		sizes[0] = n
		return 1, nil
	}

	// Calculate number of segments needed
	numSegments := (payloadLen + gsoSize - 1) / gsoSize
	if numSegments > len(outBuffs) {
		return 0, errors.New("insufficient output buffers for GSO split")
	}

	for i := 0; i < numSegments; i++ {
		start := i * gsoSize
		end := start + gsoSize
		if end > payloadLen {
			end = payloadLen
		}

		// Copy headers
		n := copy(outBuffs[i][offset:], pkt[:hdrLen])
		// Copy payload segment
		n += copy(outBuffs[i][offset+n:], payload[start:end])
		sizes[i] = n
	}

	return numSegments, nil
}

// GRODevice is an interface for devices that support GRO (Generic Receive Offload).
type GRODevice interface {
	Device
	// GRO performs Generic Receive Offload on the provided packets.
	GRO(bufs [][]byte, sizes []int, offset int) (int, error)
	// DisableTCPGRO disables TCP GRO.
	DisableTCPGRO()
	// DisableUDPGRO disables UDP GRO.
	DisableUDPGRO()
}
