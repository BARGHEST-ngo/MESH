// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !ts_omit_netstack

// Package gro implements GRO for the receive (write) path into gVisor.
package gro

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	"gvisor.dev/gvisor/pkg/buffer"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/header/parse"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"tailscale.com/net/packet"
	"tailscale.com/types/ipproto"
)

// Checksum computes the Internet checksum (RFC 1071) of the provided bytes.
// The initial checksum value can be used to chain checksums together.
func Checksum(b []byte, initial uint64) uint16 {
	ac := checksumNoFold(b, initial)
	ac = (ac >> 16) + (ac & 0xffff)
	ac = (ac >> 16) + (ac & 0xffff)
	ac = (ac >> 16) + (ac & 0xffff)
	ac = (ac >> 16) + (ac & 0xffff)
	return uint16(ac)
}

// PseudoHeaderChecksum calculates the pseudo-header checksum for TCP/UDP.
func PseudoHeaderChecksum(protocol uint8, srcAddr, dstAddr []byte, totalLen uint16) uint16 {
	sum := checksumNoFold(srcAddr, 0)
	sum = checksumNoFold(dstAddr, sum)
	sum = checksumNoFold([]byte{0, protocol}, sum)
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, totalLen)
	sum = checksumNoFold(tmp, sum)
	sum = (sum >> 16) + (sum & 0xffff)
	sum = (sum >> 16) + (sum & 0xffff)
	sum = (sum >> 16) + (sum & 0xffff)
	sum = (sum >> 16) + (sum & 0xffff)
	return uint16(sum)
}

func checksumNoFold(b []byte, initial uint64) uint64 {
	tmp := make([]byte, 8)
	binary.NativeEndian.PutUint64(tmp, initial)
	ac := binary.BigEndian.Uint64(tmp)
	var carry uint64

	for len(b) >= 128 {
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[:8]), 0)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[8:16]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[16:24]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[24:32]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[32:40]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[40:48]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[48:56]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[56:64]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[64:72]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[72:80]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[80:88]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[88:96]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[96:104]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[104:112]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[112:120]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[120:128]), carry)
		ac += carry
		b = b[128:]
	}
	if len(b) >= 64 {
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[:8]), 0)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[8:16]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[16:24]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[24:32]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[32:40]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[40:48]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[48:56]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[56:64]), carry)
		ac += carry
		b = b[64:]
	}
	if len(b) >= 32 {
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[:8]), 0)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[8:16]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[16:24]), carry)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[24:32]), carry)
		ac += carry
		b = b[32:]
	}
	if len(b) >= 16 {
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[:8]), 0)
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[8:16]), carry)
		ac += carry
		b = b[16:]
	}
	if len(b) >= 8 {
		ac, carry = bits.Add64(ac, binary.NativeEndian.Uint64(b[:8]), 0)
		ac += carry
		b = b[8:]
	}
	if len(b) >= 4 {
		ac, carry = bits.Add64(ac, uint64(binary.NativeEndian.Uint32(b[:4])), 0)
		ac += carry
		b = b[4:]
	}
	if len(b) >= 2 {
		ac, carry = bits.Add64(ac, uint64(binary.NativeEndian.Uint16(b[:2])), 0)
		ac += carry
		b = b[2:]
	}
	if len(b) == 1 {
		tmp := binary.NativeEndian.Uint16([]byte{b[0], 0})
		ac, carry = bits.Add64(ac, uint64(tmp), 0)
		ac += carry
	}

	binary.NativeEndian.PutUint64(tmp, ac)
	return binary.BigEndian.Uint64(tmp)
}

// RXChecksumOffload validates IPv4, TCP, and UDP header checksums in p,
// returning an equivalent *stack.PacketBuffer if they are valid, otherwise nil.
// The set of headers validated covers where gVisor would perform validation if
// !stack.PacketBuffer.RXChecksumValidated, i.e. it satisfies
// stack.CapabilityRXChecksumOffload. Other protocols with checksum fields,
// e.g. ICMP{v6}, are still validated by gVisor regardless of rx checksum
// offloading capabilities.
func RXChecksumOffload(p *packet.Parsed) *stack.PacketBuffer {
	var (
		pn        tcpip.NetworkProtocolNumber
		csumStart int
	)
	buf := p.Buffer()

	switch p.IPVersion {
	case 4:
		if len(buf) < header.IPv4MinimumSize {
			return nil
		}
		csumStart = int((buf[0] & 0x0F) * 4)
		if csumStart < header.IPv4MinimumSize || csumStart > header.IPv4MaximumHeaderSize || len(buf) < csumStart {
			return nil
		}
		if ^Checksum(buf[:csumStart], 0) != 0 {
			return nil
		}
		pn = header.IPv4ProtocolNumber
	case 6:
		if len(buf) < header.IPv6FixedHeaderSize {
			return nil
		}
		csumStart = header.IPv6FixedHeaderSize
		pn = header.IPv6ProtocolNumber
		if p.IPProto != ipproto.ICMPv6 && p.IPProto != ipproto.TCP && p.IPProto != ipproto.UDP {
			// buf could have extension headers before a UDP or TCP header, but
			// packet.Parsed.IPProto will be set to the ext header type, so we
			// have to look deeper. We are still responsible for validating the
			// L4 checksum in this case. So, make use of gVisor's existing
			// extension header parsing via parse.IPv6() in order to unpack the
			// L4 csumStart index. This is not particularly efficient as we have
			// to allocate a short-lived stack.PacketBuffer that cannot be
			// re-used. parse.IPv6() "consumes" the IPv6 headers, so we can't
			// inject this stack.PacketBuffer into the stack at a later point.
			packetBuf := stack.NewPacketBuffer(stack.PacketBufferOptions{
				Payload: buffer.MakeWithData(bytes.Clone(buf)),
			})
			defer packetBuf.DecRef()
			// The rightmost bool returns false only if packetBuf is too short,
			// which we've already accounted for above.
			transportProto, _, _, _, _ := parse.IPv6(packetBuf)
			if transportProto == header.TCPProtocolNumber || transportProto == header.UDPProtocolNumber {
				csumLen := packetBuf.Data().Size()
				if len(buf) < csumLen {
					return nil
				}
				csumStart = len(buf) - csumLen
				p.IPProto = ipproto.Proto(transportProto)
			}
		}
	}

	if p.IPProto == ipproto.TCP || p.IPProto == ipproto.UDP {
		lenForPseudo := len(buf) - csumStart
		csum := PseudoHeaderChecksum(
			uint8(p.IPProto),
			p.Src.Addr().AsSlice(),
			p.Dst.Addr().AsSlice(),
			uint16(lenForPseudo))
		csum = Checksum(buf[csumStart:], uint64(csum))
		if ^csum != 0 {
			return nil
		}
	}

	packetBuf := stack.NewPacketBuffer(stack.PacketBufferOptions{
		Payload: buffer.MakeWithData(bytes.Clone(buf)),
	})
	packetBuf.NetworkProtocolNumber = pn
	// Setting this is not technically required. gVisor overrides where
	// stack.CapabilityRXChecksumOffload is advertised from Capabilities().
	// https://github.com/google/gvisor/blob/64c016c92987cc04dfd4c7b091ddd21bdad875f8/pkg/tcpip/stack/nic.go#L763
	// This is also why we offload for all packets since we cannot signal this
	// per-packet.
	packetBuf.RXChecksumValidated = true
	return packetBuf
}
