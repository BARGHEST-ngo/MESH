package apkverifier

import (
	"errors"
	"fmt"
	"github.com/avast/apkparser"
	"io/ioutil"
	"strings"
)

const (
	attrSignatureVersionCase = "Signature-Version"
	attrManifestVersionCase  = "Manifest-Version"

	attrSignatureVersion     = "signature-version"
	attrName                 = "name"
	attrCreatedBy            = "created-by"
	attrDigestMainAttrSuffix = "-digest-manifest-main-attributes"
	attrDigestSuffix         = "-digest"
	attrDigestSigntoolSuffix = "-digest-manifest"
	attrAndroidApkSigned     = "x-android-apk-signed"
)

type manifest struct {
	rawData            []byte
	main               map[string]string
	entries            map[string]map[string]string
	chunks             map[string][]byte
	mainAttributtesEnd int
}

type manifestParserContext struct {
	man *manifest

	data                  []byte
	pos, mark             int
	consecutiveLineBreaks int
}

func parseManifest(f *apkparser.ZipReaderFile, withChunks bool) (*manifest, error) {
	if err := f.Open(); err != nil {
		return nil, err
	}
	defer f.Close()

	var err error
	for f.Next() {
		ctx := manifestParserContext{
			man: &manifest{
				main:    make(map[string]string),
				entries: make(map[string]map[string]string),
				chunks:  make(map[string][]byte),
			},
		}

		ctx.data, err = ioutil.ReadAll(f)
		if err != nil {
			continue
		}
		ctx.man.rawData = ctx.data

		if err = ctx.parse(withChunks); err != nil {
			continue
		}

		return ctx.man, nil
	}

	if err == nil {
		err = fmt.Errorf("Failed to open %s", f.Name)
	}

	return nil, err
}

func (ctx *manifestParserContext) parse(withChunks bool) error {
	for {
		n, v, singleBlock, err := ctx.readHeader()
		if err != nil {
			return err
		} else if !singleBlock {
			break
		}
		ctx.man.main[n] = v
	}

	ctx.man.mainAttributtesEnd = ctx.pos

	mark := ctx.pos
	for {
		n, v, singleBlock, err := ctx.readHeader()
		if err != nil {
			return err
		} else if !singleBlock {
			break
		}

		if n != attrName {
			return fmt.Errorf("Entry is not named")
		}

		entryNameValue := v
		entry := ctx.man.entries[entryNameValue]
		if entry == nil {
			entry = make(map[string]string)
			ctx.man.entries[entryNameValue] = entry
		}

		for {
			n, v, singleBlock, err = ctx.readHeader()
			if err != nil {
				return err
			} else if !singleBlock {
				break
			}
			entry[n] = v
		}

		if withChunks {
			if _, prs := ctx.man.chunks[entryNameValue]; prs {
				return fmt.Errorf("More than one entry with the same name!")
			}
			ctx.man.chunks[entryNameValue] = ctx.data[mark:ctx.pos]
			mark = ctx.pos
		}
	}
	return nil
}

func (ctx *manifestParserContext) readHeader() (name, value string, singleBlock bool, err error) {
	if ctx.consecutiveLineBreaks > 1 {
		ctx.consecutiveLineBreaks = 0
		return
	}

	mark := ctx.pos
	for ; ctx.pos < len(ctx.data); ctx.pos++ {
		if ctx.data[ctx.pos] == ':' {
			name = string(ctx.data[mark:ctx.pos])

			if ctx.pos+1 > len(ctx.data) || ctx.data[ctx.pos+1] != ' ' {
				err = fmt.Errorf("Invalid header structure '%s'", name)
				return
			}

			if !ctx.isValidName(name) {
				err = fmt.Errorf("Invalid attribute name in manifest: '%s'", name)
				return
			}

			// Attributes except for Manifest-Version and Signature-Version are case-insensitive
			// https://docs.oracle.com/javase/7/docs/technotes/guides/jar/jar.html#Notes_on_Manifest_and_Signature_Files
			// Android however does not care and has case insensitive everything.
			if true /*name != attrManifestVersionCase && name != attrSignatureVersionCase*/ {
				name = strings.ToLower(name)
			}

			ctx.pos += 2 // For the ': ' separator
			break
		}
	}

	ctx.consecutiveLineBreaks = 0
	mark = ctx.pos
	last := ctx.pos
	lastCr := false

	for ctx.pos < len(ctx.data) {
		next := ctx.data[ctx.pos]
		ctx.pos++

		switch next {
		case 0:
			err = errors.New("NUL character in manifest")
			return
		case '\n':
			if lastCr {
				lastCr = false
			} else {
				ctx.consecutiveLineBreaks++
			}
			continue
		case '\r':
			lastCr = true
			ctx.consecutiveLineBreaks++
			continue
		case ' ':
			if ctx.consecutiveLineBreaks == 1 {
				value += string(ctx.data[mark:last])
				mark = ctx.pos
				ctx.consecutiveLineBreaks = 0
				continue
			}
		}

		if ctx.consecutiveLineBreaks >= 1 {
			ctx.pos--
			break
		}
		last = ctx.pos
	}

	if last > mark {
		value += string(ctx.data[mark:last])
	}
	singleBlock = ctx.consecutiveLineBreaks > 0
	return
}

func (ctx *manifestParserContext) isValidName(name string) bool {
	if len(name) == 0 || len(name) > 70 { // LINE_LENGTH_LIMIT = 72; - 2 for separator
		return false
	}

	for i := 0; i < len(name); i++ {
		c := name[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			c == '_' || c == '-' ||
			(c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
