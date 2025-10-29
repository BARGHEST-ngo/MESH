package signingblock

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"strings"
)

/*
 *  *** FROSTING BLOCK STRUCTURE ***
 *
 *  ------------------------------------------------------------
 *  |- size of the protobuf + size of the metadata array (varint)  (AKA relative offset to the signatures array)
 *   ---------------------------  SIGNED DATA START ---------------------------
 *  |- protobuf size (varint)
 *  |- protobuf frosting info (byte array)
 *  |
 *  |- metadata array size (varint)
 *  |- metadata array, pick first non-disabled
 *      |- entry size (varint)
 *      |- disabled flag (varint), checked != 0 in Play Store
 *      |- public key index (varint), from finsky.peer_app_sharing_api.frosting_public_keys comma separated array
 *      |- fileSha256 (digest over the apk before signing block, schemev2 signing block and eocd, see verifyApk()
 *        --------------------------- SIGNED DATA END ---------------------------
 *  |
 *  |- size of the signatures array (varint)
 *  |- signatures array, indexing matches the metadata array
 *      |- entry size (varint)
 *      |- signature (byte array)
 *  -------------------------------------------------------------
 *
 * All the size fields are in bytes and counted excluding the size varint itself.
 *
 * The SIGNED DATA end after the metadata entry you're checking againts, even if it is not the last entry.
 * They are verified against signature created by frostingPublicKeys key. It uses ECDSAWithSHA256 algorithm.
 * The fileSha256 hash must match the APK for the frosting to be valid.
 *
 * The 'protobuf frosting info' is a rather complex protobuf structure. It contains some APK's metadata,
 * nothing really useful it seems. It does not seem to be relevant to the frosting signature's validity,
 * so I have not examined it further.
 * Interestingly, Google Photos apk from apkmirror.com contains 'com.google.android.apps.photos.PIXEL_2018_PRELOAD'
 * string inside this protobuf, perhaps it is possible to distinguish preloaded /system apps from the Play Store apks?
 * Some apk's have very short protobuf infos (Netflix, 200 bytes), some have much longer (Facebook, 2200 bytes).
 * Here's an example of parsed protobuf info from the Netflix app:
 *
 *   1 <varint> = 1                                                        // frosting versions?
 *   2 <varint> = 0
 *   3 <varint> = 1
 *   4 <varint> = 1541545744578                                            // Timestamp of the frosting creation?
 *   5 <chunk> = message:
 *       8 <chunk> = message:
 *           1 <chunk> = message(1 <varint> = 22)                          // minSdkLevel?
 *           6 <varint> = 2
 *       9 <chunk> = message:
 *           1 <chunk> = message(1 <varint> = 2266, 4 <varint> = 2)        // versionCode
 *           2 <chunk> = message(1 <varint> = 50003, 4 <varint> = 4)
 *       10 <chunk> = message:
 *           1 <chunk> = bytes (30)                                        // ?? only last byte changes across apks
 *               0000   FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FE FF FF FF FF FF FF  ........................
 *               0018   FF FF FF FF FF 3F                                                        .....?
 *           3 <chunk> = message:
 *               1 <chunk> = bytes (32)                                   // sha256 of something?
 *                   0000   16 F8 22 A6 93 26 89 34 D8 2A 88 BB 8C AD B6 68 2C EB 77 A8 AA E4 5F AA  .."..&.4.*.....h,.w..._.
 *                   0018   F9 3C CA 63 44 2A A4 B9                                                  .<.cD*..
 *               2 <varint> = 20
 */

const (
	// Value 'finsky.peer_app_sharing_api.frosting_public_keys' from com.android.vending 11.2.14-all [0] [PR] 207207134
	frostingPublicKeys = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEZH2+1+E07dnErAD3L6BbTnaohU0bbXriNlJI7VxJU+LjdSwPyXR5pomARAMoyPkMksLz/gitUPtFuJoPL2ziEw=="
)

var (
	ErrFrostingInvalidSignature = errors.New("invalid frosting signature")
	ErrFrostingDigestMismatch   = errors.New("frosting apk file digest mismatch")
)

type frostingInfo struct {
	fileSha256   []byte
	protobufInfo []byte

	usedKeySha256 string
}

func (f *frostingInfo) readInt8(r io.ByteReader) (int32, error) {
	x, err := r.ReadByte()
	return int32(int8(x)), err
}

func (f *frostingInfo) readInt32(r io.ByteReader) (int32, error) {
	v0, err := f.readInt8(r)
	if err != nil {
		return 0, err
	}

	if v0 < 0 {
		v0 &= 0x7f

		v1, err := f.readInt8(r)
		if err != nil {
			return 0, err
		}

		if v1 < 0 {
			v0 |= (v1 & 0x7F) << 7
			if v1, err = f.readInt8(r); err != nil {
				return 0, err
			}

			if v1 < 0 {
				v0 |= (v1 & 0x7F) << 14
				if v1, err = f.readInt8(r); err != nil {
					return 0, err
				}

				if v1 < 0 {
					v2, err := f.readInt8(r)
					if err != nil {
						return 0, err
					}

					v0 = v0 | (v1&0x7F)<<21 | v2<<28
					if v2 < 0 {
						return 0, fmt.Errorf("varint overflow")
					}
				} else {
					v0 |= v1 << 21
				}
			} else {
				v0 |= v1 << 14
			}
		} else {
			v0 |= v1 << 7
		}
	}
	return v0, nil
}

func (f *frostingInfo) readArrayItem(bf *bytes.Reader, index int) ([]byte, error) {
	totalSize, err := f.readInt32(bf)
	if err != nil {
		return nil, err
	}

	if totalSize <= 0 || int(totalSize) > bf.Len() {
		return nil, fmt.Errorf("invalid base for array")
	}

	startPos, _ := bf.Seek(0, io.SeekCurrent)
	for i := 0; i < index; i++ {
		entrySize, err := f.readInt32(bf)
		if err != nil {
			return nil, err
		}

		if pos, _ := bf.Seek(0, io.SeekCurrent); entrySize < 0 || pos+int64(entrySize) >= startPos+int64(totalSize) {
			return nil, fmt.Errorf("invalid next value")
		}

		if _, err := bf.Seek(int64(entrySize), io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek to next")
		}
	}

	size, err := f.readInt32(bf)
	if err != nil {
		return nil, fmt.Errorf("failed to read size")
	}

	if pos, _ := bf.Seek(0, io.SeekCurrent); size < 0 || pos+int64(size) > int64(totalSize)+startPos {
		return nil, fmt.Errorf("invalid size value")
	}

	res := make([]byte, size)
	_, err = bf.Read(res)
	return res, err
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (f *frostingInfo) verifySignature(signed, signature []byte, keyBase64 string) error {
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(keyBase64))
	pubKeyAsn, err := ioutil.ReadAll(dec)
	if err != nil {
		return fmt.Errorf("failed to parse key base64: %s", err.Error())
	}

	keyDigest := sha256.Sum256(pubKeyAsn)
	f.usedKeySha256 = hex.EncodeToString(keyDigest[:])

	pkGen, err := x509.ParsePKIXPublicKey(pubKeyAsn)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pk: %s", err.Error())
	}

	pk, ok := pkGen.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid key type: %T", pkGen)
	}

	digest := sha256.Sum256(signed)

	ecdsaSig := new(ecdsaSignature)
	if rest, err := asn1.Unmarshal(signature, ecdsaSig); err != nil {
		return err
	} else if len(rest) != 0 {
		return errors.New("x509: trailing data after ECDSA signature")
	}

	if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
		return errors.New("x509: ECDSA signature contained zero or negative values")
	}

	if !ecdsa.Verify(pk, digest[:], ecdsaSig.R, ecdsaSig.S) {
		return ErrFrostingInvalidSignature
	}

	return nil
}

func (f *frostingInfo) parseFrostingBlock(block []byte) error {
	bf := bytes.NewReader(block)

	protoAndMetaSize, err := f.readInt32(bf)
	if err != nil || protoAndMetaSize <= 0 || int(protoAndMetaSize) > bf.Len() {
		return fmt.Errorf("invalid protoAndMetaSize value")
	}
	offAfterSignaturesEnd, _ := bf.Seek(0, io.SeekCurrent)

	protobufSize, err := f.readInt32(bf)
	frostingProtobufStart, _ := bf.Seek(0, io.SeekCurrent)

	if err != nil || protobufSize <= 0 || int(protobufSize) > bf.Len() ||
		protobufSize > protoAndMetaSize-int32(frostingProtobufStart-offAfterSignaturesEnd) {
		return fmt.Errorf("invalid protobufSize value")
	}

	if _, err := bf.Seek(int64(protobufSize), io.SeekCurrent); err != nil {
		return fmt.Errorf("seek to metadata start failed %s", err.Error())
	}

	metadataSize, err := f.readInt32(bf)
	if err != nil || metadataSize <= 0 || metadataSize > (protoAndMetaSize-protobufSize) {
		return fmt.Errorf("invalid metadataSize")
	}

	limit, _ := bf.Seek(0, io.SeekCurrent)
	limit += int64(metadataSize)
	for signatureIdx := 0; signatureIdx < math.MaxInt16; signatureIdx++ {
		pos, _ := bf.Seek(0, io.SeekCurrent)
		if pos >= limit {
			break
		}

		nextKey, err := f.readInt32(bf)
		pos, _ = bf.Seek(0, io.SeekCurrent)
		nextKey += int32(pos)
		if err != nil || int64(nextKey) > limit || nextKey <= 0 {
			return fmt.Errorf("invalid nextKey offset")
		}

		disabled, err := f.readInt32(bf)
		if pos, _ := bf.Seek(0, io.SeekCurrent); err != nil || pos > int64(nextKey) {
			return fmt.Errorf("invalid 'disabled' read")
		}

		if disabled != 0 {
			if _, err := bf.Seek(int64(nextKey), io.SeekStart); err != nil {
				return fmt.Errorf("can't seek after disabled: %s", err.Error())
			}
			continue
		}

		keyIndex, err := f.readInt32(bf)
		keys := strings.Split(frostingPublicKeys, ",")
		if err != nil || int(keyIndex) >= len(keys) {
			return fmt.Errorf("invalid/unknown key used in this frosting (idx %d)!", keyIndex)
		}

		pos, _ = bf.Seek(0, io.SeekCurrent)
		if int64(nextKey)-pos != sha256.Size {
			return fmt.Errorf("invalind key index length")
		}

		fileSha256 := make([]byte, sha256.Size)
		if _, err := bf.Read(fileSha256); err != nil {
			return fmt.Errorf("failed to read signature: %s", err.Error())
		}

		if _, err := bf.Seek(offAfterSignaturesEnd+int64(protoAndMetaSize), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek after signature read: %s", err.Error())
		}

		signature, err := f.readArrayItem(bf, signatureIdx)
		if err != nil {
			return fmt.Errorf("failed to read signature: %s", err.Error())
		}

		if _, err := bf.Seek(offAfterSignaturesEnd, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to signed data: %s", err.Error())
		}

		signedData := make([]byte, int64(nextKey)-offAfterSignaturesEnd)
		if _, err := bf.Read(signedData); err != nil {
			return fmt.Errorf("failed to read signed data: %s", err.Error())
		}

		if err := f.verifySignature(signedData, signature, keys[keyIndex]); err != nil {
			if err == ErrFrostingInvalidSignature {
				return err
			}
			return fmt.Errorf("failed to verify signature: %s", err.Error())
		}

		if _, err := bf.Seek(frostingProtobufStart, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to signed data: %s", err.Error())
		}

		protobufInfo := make([]byte, protobufSize)
		if _, err := bf.Read(protobufInfo); err != nil {
			return fmt.Errorf("failed to read protobufInfo data: %s", err.Error())
		}

		f.fileSha256 = fileSha256
		f.protobufInfo = protobufInfo
		return nil
	}

	return fmt.Errorf("no enabled signature found")
}

func (f *frostingInfo) parse(block []byte) (string, []byte, error) {
	err := f.parseFrostingBlock(block)
	return f.usedKeySha256, f.protobufInfo, err
}

func (f *frostingInfo) hashFileSection(hasher hash.Hash, r io.ReadSeeker, offset int64, size int) error {
	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	buf := make([]byte, 0x2000)
	var chunk int
	var err error
	for i := 0; i < size; i += chunk {
		chunk = len(buf)
		if chunk > (size - i) {
			chunk = size - i
		}

		chunk, err = io.ReadFull(r, buf[:chunk])
		if err == io.EOF {
			if chunk <= 0 {
				return err
			}
		} else if err != nil {
			return err
		}

		hasher.Write(buf[:chunk])
	}
	return nil
}

func (f *frostingInfo) verifyApk(r io.ReadSeeker, signingBlockOffset int64, schemeV2block []byte, zipCdOffset, zipCdSize int64, eocdOrig []byte) error {
	hasher := sha256.New()

	if err := f.hashFileSection(hasher, r, 0, int(signingBlockOffset)); err != nil {
		return fmt.Errorf("failed to hash apk: %s", err.Error())
	}

	if len(schemeV2block) != 0 {
		binary.Write(hasher, binary.LittleEndian, uint32(blockIdSchemeV2))
		hasher.Write(schemeV2block)
	}

	if err := f.hashFileSection(hasher, r, zipCdOffset, int(zipCdSize)); err != nil {
		return fmt.Errorf("failed to hash zip central directory: %s", err.Error())
	}

	// For the purposes of integrity verification, ZIP End of Central Directory's field Start of
	// Central Directory must be considered to point to the offset of the APK Signing Block.
	eocd := make([]byte, len(eocdOrig))
	copy(eocd, eocdOrig)
	binary.LittleEndian.PutUint32(eocd[eocdCentralDirOffsetOffset:], uint32(signingBlockOffset))
	hasher.Write(eocd)

	if !bytes.Equal(hasher.Sum(nil), f.fileSha256) {
		return ErrFrostingDigestMismatch
	}

	return nil
}
