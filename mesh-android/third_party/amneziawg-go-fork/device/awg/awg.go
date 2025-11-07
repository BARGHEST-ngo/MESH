package awg

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/tevino/abool"
)

type Cfg struct {
	IsSet                     bool
	JunkPacketCount           int
	JunkPacketMinSize         int
	JunkPacketMaxSize         int
	InitHeaderJunkSize        int
	ResponseHeaderJunkSize    int
	CookieReplyHeaderJunkSize int
	TransportHeaderJunkSize   int

	MagicHeaders MagicHeaders
}

type Protocol struct {
	IsOn abool.AtomicBool
	// TODO: revision the need of the mutex
	Mux         sync.RWMutex
	Cfg         Cfg
	JunkCreator JunkCreator

	HandshakeHandler SpecialHandshakeHandler
}

func (protocol *Protocol) CreateInitHeaderJunk() ([]byte, error) {
	// DEBUG: Nil check
	if protocol == nil {
		panic("CreateInitHeaderJunk: protocol is nil")
	}

	protocol.Mux.RLock()
	defer protocol.Mux.RUnlock()

	return protocol.createHeaderJunk(protocol.Cfg.InitHeaderJunkSize, 0)
}

func (protocol *Protocol) CreateResponseHeaderJunk() ([]byte, error) {
	protocol.Mux.RLock()
	defer protocol.Mux.RUnlock()

	return protocol.createHeaderJunk(protocol.Cfg.ResponseHeaderJunkSize, 0)
}

func (protocol *Protocol) CreateCookieReplyHeaderJunk() ([]byte, error) {
	protocol.Mux.RLock()
	defer protocol.Mux.RUnlock()

	return protocol.createHeaderJunk(protocol.Cfg.CookieReplyHeaderJunkSize, 0)
}

func (protocol *Protocol) CreateTransportHeaderJunk(packetSize int) ([]byte, error) {
	protocol.Mux.RLock()
	defer protocol.Mux.RUnlock()

	return protocol.createHeaderJunk(protocol.Cfg.TransportHeaderJunkSize, packetSize)
}

func (protocol *Protocol) createHeaderJunk(junkSize int, extraSize int) ([]byte, error) {
	// DEBUG: Nil check
	if protocol == nil {
		panic("createHeaderJunk: protocol is nil")
	}

	if junkSize == 0 {
		return nil, nil
	}

	buf := make([]byte, 0, junkSize+extraSize)
	writer := bytes.NewBuffer(buf[:0])

	err := protocol.JunkCreator.AppendJunk(writer, junkSize)
	if err != nil {
		return nil, fmt.Errorf("append junk: %w", err)
	}

	return writer.Bytes(), nil
}

func (protocol *Protocol) GetMagicHeaderMinFor(msgType uint32) (uint32, error) {
	// Defensive: If MagicHeaders not initialized, return the msgType as-is (standard WireGuard behavior)
	if protocol == nil || protocol.Cfg.MagicHeaders.Values == nil {
		return msgType, nil
	}

	for _, magicHeader := range protocol.Cfg.MagicHeaders.Values {
		if magicHeader.Min <= msgType && msgType <= magicHeader.Max {
			return magicHeader.Min, nil
		}
	}

	return 0, fmt.Errorf("no header for value: %d", msgType)
}

func (protocol *Protocol) GetMsgType(defaultMsgType uint32) (uint32, error) {
	// DEBUG: Defensive nil checks
	if protocol == nil {
		panic("GetMsgType: protocol is nil")
	}
	if protocol.Cfg.MagicHeaders.Values == nil {
		panic(fmt.Sprintf("GetMsgType: MagicHeaders.Values is nil (defaultMsgType=%d)", defaultMsgType))
	}
	return protocol.Cfg.MagicHeaders.Get(defaultMsgType)
}
