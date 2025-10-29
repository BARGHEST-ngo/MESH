package awg

import (
	"sync"

	"github.com/tevino/abool"
)

type Protocol struct {
	IsASecOn abool.AtomicBool
	// TODO: revision the need of the mutex
	ASecMux     sync.RWMutex
	ASecCfg     aSecCfgType
	JunkCreator junkCreator

	HandshakeHandler SpecialHandshakeHandler
}

type aSecCfgType struct {
	IsSet                      bool
	JunkPacketCount            int
	JunkPacketMinSize          int
	JunkPacketMaxSize          int
	InitPacketJunkSize         int
	ResponsePacketJunkSize     int
	InitPacketMagicHeader      uint32
	ResponsePacketMagicHeader  uint32
	UnderloadPacketMagicHeader uint32
	TransportPacketMagicHeader uint32
}
