package apkparser

import (
	"encoding/binary"
	"io"
)

// frameworks/base/libs/androidfw/include/androidfw/ResourceTypes.h
const (
	chunkNull          = 0x0000
	chunkStringTable   = 0x0001
	chunkTable         = 0x0002
	chunkAxmlFile      = 0x0003
	chunkResourceIds   = 0x0180
	chunkTablePackage  = 0x0200
	chunkTableType     = 0x0201
	chunkTableTypeSpec = 0x0202
	chunkTableLibrary  = 0x0203

	chunkMaskXml        = 0x0100
	chunkXmlNsStart     = 0x0100
	chunkXmlNsEnd       = 0x0101
	chunkXmlTagStart    = 0x0102
	chunkXmlTagEnd      = 0x0103
	chunkXmlText        = 0x0104
	chunkXmlLastChunk   = 0x017f
	chunkXmlResourceMap = 0x0180

	chunkHeaderSize = (2 + 2 + 4)
)

type ResAttr struct {
	NamespaceId uint32
	NameIdx     uint32
	RawValueIdx uint32
	Res         ResValue
}

type ResValue struct {
	Size uint16
	Res0 uint8 // padding
	Type AttrType
	Data uint32
}

type AttrType uint8

const (
	AttrTypeNull          AttrType = 0x00
	AttrTypeReference              = 0x01
	AttrTypeAttribute              = 0x02
	AttrTypeString                 = 0x03
	AttrTypeFloat                  = 0x04
	AttrTypeIntDec                 = 0x10
	AttrTypeIntHex                 = 0x11
	AttrTypeIntBool                = 0x12
	AttrTypeIntColorArgb8          = 0x1c
	AttrTypeIntColorRgb8           = 0x1d
	AttrTypeIntColorArgb4          = 0x1e
	AttrTypeIntColorRgb4           = 0x1f
)

func parseChunkHeader(r io.Reader) (id, headerLen uint16, len uint32, err error) {
	if err = binary.Read(r, binary.LittleEndian, &id); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &headerLen); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &len); err != nil {
		return
	}

	return
}
