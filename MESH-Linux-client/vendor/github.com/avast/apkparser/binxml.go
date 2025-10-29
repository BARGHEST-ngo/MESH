package apkparser

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"unsafe"
)

type binxmlParseInfo struct {
	strings     stringTable
	resourceIds []uint32
	openTags    []xml.Name

	encoder ManifestEncoder
	res     *ResourceTable
}

// Some samples have manifest in plaintext, this is an error.
// 2c882a2376034ed401be082a42a21f0ac837689e7d3ab6be0afb82f44ca0b859
var ErrPlainTextManifest = errors.New("xml is in plaintext, binary form expected")

// Deprecated: just calls ParseXML
func ParseManifest(r io.Reader, enc ManifestEncoder, resources *ResourceTable) error {
	return ParseXml(r, enc, resources)
}

// Parse the binary Xml format. The resources are optional and can be nil.
func ParseXml(r io.Reader, enc ManifestEncoder, resources *ResourceTable) error {
	x := binxmlParseInfo{
		encoder: enc,
		res:     resources,
	}

	id, headerLen, totalLen, err := parseChunkHeader(r)
	if err != nil {
		return err
	}

	if (id & 0xFF) == '<' {
		buf := bytes.NewBuffer(make([]byte, 0, 8))
		binary.Write(buf, binary.LittleEndian, &id)
		binary.Write(buf, binary.LittleEndian, &headerLen)
		binary.Write(buf, binary.LittleEndian, &totalLen)

		if s := buf.String(); strings.HasPrefix(s, "<?xml ") || strings.HasPrefix(s, "<manif") {
			return ErrPlainTextManifest
		}
	}

	// Android doesn't care.
	/*if id != chunkAxmlFile {
	    return fmt.Errorf("Invalid top chunk id: 0x%08x", id)
	}*/

	defer x.encoder.Flush()

	totalLen -= uint32(headerLen)
	// 29f82928c630897576aa0c9c3d2f36d722a94ec574bd5ca4aaf4ba5f9d014a32
	io.CopyN(ioutil.Discard, r, int64(headerLen)-chunkHeaderSize)

	var len uint32
	var lastId uint16
	for i := uint32(0); i < totalLen; i += len {
		id, headerLen, len, err = parseChunkHeader(r)
		if err != nil {
			return fmt.Errorf("Error parsing header at 0x%08x of 0x%08x %08x: %s", i, totalLen, lastId, err.Error())
		}

		lastId = id

		lm := &io.LimitedReader{R: r, N: int64(len) - 2*4}

		switch id {
		case chunkStringTable:
			x.strings, err = parseStringTable(lm)
		case chunkResourceIds:
			err = x.parseResourceIds(lm)
		default:
			if (id & chunkMaskXml) == 0 {
				// Android ignores it
				// A9842ACCCE9D21549479EEB84F5E2033DA5E479EE1E183D48445BC99324ED983
				//err = fmt.Errorf("Unknown chunk id 0x%x %d", id, len)
				break
			}

			// skip line number and unknown 0xFFFFFFFF
			if _, err = io.CopyN(ioutil.Discard, lm, int64(headerLen)-8); err != nil {
				break
			}

			switch id {
			case chunkXmlNsStart:
				err = x.parseNsStart(lm)
			case chunkXmlNsEnd:
				err = x.parseNsEnd(lm)
			case chunkXmlTagStart:
				err = x.parseTagStart(lm)
			case chunkXmlTagEnd:
				err = x.parseTagEnd(lm)
			case chunkXmlText:
				err = x.parseText(lm)
			case chunkXmlLastChunk: // unimplemented
			case chunkXmlResourceMap: // unimplemented
			default:
				err = fmt.Errorf("Unknown chunk id 0x%x", id)
			}
		}

		if err == ErrEndParsing {
			break
		} else if err != nil {
			return fmt.Errorf("Chunk: 0x%08x: %s", id, err.Error())
		} else if lm.N != 0 {
			// da62a1edc4d9826c8bf2ed8d5be857614f7908163269d80f9d4ad9ee4d12405e
			io.CopyN(ioutil.Discard, lm, lm.N)
			//return fmt.Errorf("Chunk: 0x%08x: was not fully read (%d remaining)", id, lm.N)
		}
	}

	return x.encoder.Flush()
}

func (x *binxmlParseInfo) parseResourceIds(r *io.LimitedReader) error {
	if (r.N % 4) != 0 {
		return fmt.Errorf("Invalid chunk size!")
	}

	count := uint32(r.N / 4)
	var id uint32
	for i := uint32(0); i < count; i++ {
		if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
			return err
		}
		x.resourceIds = append(x.resourceIds, id)
	}
	return nil
}

func (x *binxmlParseInfo) parseNsStart(r *io.LimitedReader) error {
	var err error
	ns := &xml.Name{}

	var idx uint32
	if err = binary.Read(r, binary.LittleEndian, &idx); err != nil {
		return err
	}

	if ns.Local, err = x.strings.get(idx); err != nil {
		return err
	}

	if err = binary.Read(r, binary.LittleEndian, &idx); err != nil {
		return err
	}

	if ns.Space, err = x.strings.get(idx); err != nil {
		return err
	}

	// TODO: what to do with this?
	_ = ns
	return nil
}

func (x *binxmlParseInfo) parseNsEnd(r *io.LimitedReader) error {
	if _, err := io.CopyN(ioutil.Discard, r, 2*4); err != nil {
		return fmt.Errorf("error skipping: %s", err.Error())
	}

	// TODO: what to do with this?
	return nil
}

func mapDisallowedNameRunes(r rune) rune {
	// https://www.w3.org/TR/REC-xml/#NT-Name
	// ":" | [A-Z] | "_" | [a-z] | [#xC0-#xD6] | [#xD8-#xF6] | [#xF8-#x2FF] | [#x370-#x37D] | [#x37F-#x1FFF] | [#x200C-#x200D] | [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF] | [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
	//NameStartChar | "-" | "." | [0-9] | #xB7 | [#x0300-#x036F] | [#x203F-#x2040]

	switch r {
	case ':' | '_' | '-' | '.' | '\u00B7':
		return r
	}

	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
		return r
	}

	if (r >= '\u00C0' && r <= '\u00D6') || (r >= '\u00D8' && r <= '\u00F6') || (r >= '\u00F8' && r <= '\u02FF') {
		return r
	}

	if (r >= '\u0300' && r <= '\u037D') || (r >= '\u037F' && r <= '\u1FFF') || (r >= '\u200C' && r <= '\u200D') {
		return r
	}

	if (r >= '\u203F' && r <= '\u2040') || (r >= '\u2070' && r <= '\u218F') || (r >= '\u2C00' && r <= '\u2FEF') {
		return r
	}

	if (r >= '\u3001' && r <= '\uD7FF') || (r >= '\uF900' && r <= '\uFDCF') || (r >= '\uFDF0' && r <= '\uFFFD') {
		return r
	}

	if r >= '\U00010000' && r <= '\U000EFFFF' {
		return r
	}

	return '_'
}

func (x *binxmlParseInfo) parseTagStart(r *io.LimitedReader) error {
	var namespaceIdx, nameIdx uint32
	var attrStart, attrSize, attrCount uint16

	if err := binary.Read(r, binary.LittleEndian, &namespaceIdx); err != nil {
		return fmt.Errorf("error reading namespace idx: %s", err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &nameIdx); err != nil {
		return fmt.Errorf("error reading name idx: %s", err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &attrStart); err != nil {
		return fmt.Errorf("error reading attrStart: %s", err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &attrSize); err != nil {
		return fmt.Errorf("error reading attrSize: %s", err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &attrCount); err != nil {
		return fmt.Errorf("error reading attrCount: %s", err.Error())
	}

	// attrStart - Byte offset from the start of this structure where the attributes start.
	// A9842ACCCE9D21549479EEB84F5E2033DA5E479EE1E183D48445BC99324ED983
	io.CopyN(io.Discard, r, int64(attrStart)-14)

	namespace, err := x.strings.get(namespaceIdx)
	if err != nil {
		return fmt.Errorf("error decoding namespace: %s", err.Error())
	}

	name, err := x.strings.get(nameIdx)
	if err != nil {
		return fmt.Errorf("error decoding name: %s", err.Error())
	}

	tok := xml.StartElement{
		Name: xml.Name{Local: name, Space: namespace},
	}

	var attr ResAttr
	for i := uint16(0); i < attrCount; i++ {
		if err := binary.Read(r, binary.LittleEndian, &attr); err != nil {
			return fmt.Errorf("error reading attrData: %s", err.Error())
		}

		if uintptr(attrSize) > unsafe.Sizeof(attr) {
			io.CopyN(io.Discard, r, int64(uintptr(attrSize)-unsafe.Sizeof(attr)))
		}

		// Android actually reads attributes purely by their IDs (see frameworks/base/core/res/res/values/attrs_manifest.xml
		// and its generated R class, that's where the indexes come from, namely the AndroidManifestActivity array)
		// but good guy android actually puts the strings into the string table on the same indexes anyway, most of the time.
		// This is for the samples that don't have it, mostly due to obfuscators/minimizers.
		// The ID can't change, because it would break current APKs.
		// Sample: 98d2e837b8f3ac41e74b86b2d532972955e5352197a893206ecd9650f678ae31
		//
		// The exception to this rule is the "package" attribute in the root manifest tag. That one MUST NOT use
		// resource ids, instead, it needs to use the string table. The meta attrs 'platformBuildVersion*'
		// are the same, except Android never parses them so it's just for manual analysis.
		// Sample: a3ee88cf1492237a1be846df824f9de30a6f779973fe3c41c7d7ed0be644ba37
		//
		// In general, android doesn't care about namespaces, but if a resource ID is used, it has to have been
		// in the android: namespace, so we fix that up.

		// frameworks/base/core/jni/android_util_AssetManager.cpp android_content_AssetManager_retrieveAttributes
		// frameworks/base/core/java/android/content/pm/PackageParser.java parsePackageSplitNames
		var attrName string
		if attr.NameIdx < uint32(len(x.resourceIds)) {
			attrName = getAttributteName(x.resourceIds[attr.NameIdx])
		}

		var attrNameFromStrings string
		if attrName == "" || name == "manifest" {
			attrNameFromStrings, err = x.strings.get(attr.NameIdx)
			if err != nil {
				if attrName == "" {
					return fmt.Errorf("error decoding attrNameIdx: %s", err.Error())
				}
			} else if attrName != "" && attrNameFromStrings != "package" && !strings.HasPrefix(attrNameFromStrings, "platformBuildVersion") {
				attrNameFromStrings = ""
			}
		}

		attrNameSpace, err := x.strings.get(attr.NamespaceId)
		if err != nil {
			return fmt.Errorf("error decoding attrNamespaceIdx: %s", err.Error())
		}

		if attrNameFromStrings != "" {
			attrName = strings.Map(mapDisallowedNameRunes, attrNameFromStrings)
		} else if attrNameSpace == "" {
			attrNameSpace = "http://schemas.android.com/apk/res/android"
		}

		resultAttr := xml.Attr{
			Name: xml.Name{Local: attrName, Space: attrNameSpace},
		}

		switch attr.Res.Type {
		case AttrTypeString:
			resultAttr.Value, err = x.strings.get(attr.RawValueIdx)
			if err != nil {
				// da62a1edc4d9826c8bf2ed8d5be857614f7908163269d80f9d4ad9ee4d12405e
				resultAttr.Value = fmt.Sprintf("#%d", attr.RawValueIdx)
				err = nil
				//return fmt.Errorf("error decoding attrStringIdx: %s", err.Error())
			}
		case AttrTypeIntBool:
			resultAttr.Value = strconv.FormatBool(attr.Res.Data != 0)
		case AttrTypeIntHex:
			resultAttr.Value = fmt.Sprintf("0x%x", attr.Res.Data)
		case AttrTypeFloat:
			val := (*float32)(unsafe.Pointer(&attr.Res.Data))
			resultAttr.Value = fmt.Sprintf("%g", *val)
		case AttrTypeReference:
			isValidString := false
			if x.res != nil {
				var e *ResourceEntry
				if resultAttr.Name.Local == "icon" || resultAttr.Name.Local == "roundIcon" {
					e, err = x.res.GetIconPng(attr.Res.Data)
				} else {
					e, err = x.res.GetResourceEntry(attr.Res.Data)
				}

				if err == nil {
					resultAttr.Value, err = e.value.String()
					isValidString = err == nil
				}
			}

			if !isValidString && resultAttr.Value == "" {
				resultAttr.Value = fmt.Sprintf("@%x", attr.Res.Data)
			}
		default:
			resultAttr.Value = strconv.FormatInt(int64(int32(attr.Res.Data)), 10)
		}

		tok.Attr = append(tok.Attr, resultAttr)
	}

	x.openTags = append(x.openTags, tok.Name)

	return x.encoder.EncodeToken(tok)
}

func (x *binxmlParseInfo) parseTagEnd(r *io.LimitedReader) error {
	var namespaceIdx, nameIdx uint32
	if err := binary.Read(r, binary.LittleEndian, &namespaceIdx); err != nil {
		return fmt.Errorf("error reading namespace idx: %s", err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &nameIdx); err != nil {
		return fmt.Errorf("error reading name idx: %s", err.Error())
	}

	namespace, err := x.strings.get(namespaceIdx)
	if err != nil {
		return fmt.Errorf("error decoding namespace: %s", err.Error())
	}

	name, err := x.strings.get(nameIdx)
	if err != nil {
		// 4D8029A256A7FC3571BC497F9B6D1D734A5F2D4D95E032A47AE86F2C6812DCEB
		if len(x.openTags) != 0 {
			name = x.openTags[len(x.openTags)-1].Local
		} else {
			return fmt.Errorf("error decoding name: %s", err.Error())
		}
	}

	if len(x.openTags) != 0 {
		x.openTags = x.openTags[:len(x.openTags)-1]
	}

	return x.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: name, Space: namespace}})
}

func (x *binxmlParseInfo) parseText(r *io.LimitedReader) error {
	var idx uint32
	if err := binary.Read(r, binary.LittleEndian, &idx); err != nil {
		return fmt.Errorf("error reading idx: %s", err.Error())
	}

	text, err := x.strings.get(idx)
	if err != nil {
		return fmt.Errorf("error decoding idx: %s", err.Error())
	}

	if _, err := io.CopyN(ioutil.Discard, r, 2*4); err != nil {
		return fmt.Errorf("error skipping: %s", err.Error())
	}

	return x.encoder.EncodeToken(xml.CharData(text))
}
