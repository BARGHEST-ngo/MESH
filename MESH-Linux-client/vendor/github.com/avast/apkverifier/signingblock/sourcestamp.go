package signingblock

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/avast/apkparser"

	"github.com/avast/apkverifier/internal/x509andr"
)

const (
	sourceStampAttrProofOfRotation = 0x9d6303f7
	sourceStampAttrTime            = 0xe43c5946

	sourceStampZipEntryName  = "stamp-cert-sha256"
	sourceStampHashSizeLimit = 64 * 1024
)

type SourceStampLineageNode struct {
	Cert       *x509.Certificate
	ParentAlgo SignatureAlgorithm
	Algo       SignatureAlgorithm
	Signature  []byte
	Flags      int32
}

type SourceStampResult struct {
	Cert        *x509.Certificate
	SigningTime time.Time
	Lineage     []*SourceStampLineageNode
	Errors      []error
	Warnings    []string
}

type SourceStampCertMismatchError struct {
	CertInApkSha256          string
	CertInSigningBlockSha256 string
}

func (e *SourceStampCertMismatchError) Error() string {
	return fmt.Sprintf("Mismatch between cert in %s and signing block: %s != %s",
		sourceStampZipEntryName, e.CertInApkSha256, e.CertInSigningBlockSha256)
}

type sourceStampVerifier struct {
	verifiedSchemeId             int32
	minSdkVersion, maxSdkVersion int32
	res                          SourceStampResult
}

func (v *sourceStampVerifier) addError(format string, args ...interface{}) {
	v.res.Errors = append(v.res.Errors, fmt.Errorf(format, args...))
}

func (v *sourceStampVerifier) addWarning(format string, args ...interface{}) {
	v.res.Warnings = append(v.res.Warnings, fmt.Sprintf(format, args...))
}

func (v *sourceStampVerifier) VerifySourceV2Stamp(zip *apkparser.ZipReader, block []byte, contentDigests map[contentDigest][]byte) (res *SourceStampResult) {
	res = &v.res
	srcStampHashEntry := zip.File[sourceStampZipEntryName]
	if block == nil {
		if srcStampHashEntry != nil {
			v.addError("File %s is present in the APK, but SourceStampV2 block is missing.", sourceStampZipEntryName)
			return
		}
		return nil
	}

	stampBlock, err := getLenghtPrefixedSlice(bytes.NewBuffer(block))
	if err != nil {
		v.addError("failed to parse top-level block: %s", err.Error())
		return
	}

	v.res.Cert = v.parseCertificate(stampBlock)
	if v.res.Cert == nil {
		return
	}

	if srcStampHashEntry == nil {
		v.addWarning("Block SourceStampV2 is present in the APK, but %s file is missing.", sourceStampZipEntryName)
		return
	}

	if len(contentDigests) == 0 {
		digest := v.getSchemeV1ContentDigest(zip)
		if contentDigests == nil {
			contentDigests = map[contentDigest][]byte{}
		}
		contentDigests[digestNonChunkedSha256] = digest
	}

	var srcStampHash []byte
	srcStampHash, err = srcStampHashEntry.ReadAll(sourceStampHashSizeLimit)
	if err != nil {
		v.addError("failed to read from %s: %s", sourceStampZipEntryName, err.Error())
		return
	} else if len(srcStampHash) >= sourceStampHashSizeLimit {
		v.addError("The %s file is too big", sourceStampZipEntryName)
		return
	}

	if certSha256 := sha256.Sum256(v.res.Cert.Raw); !bytes.Equal(certSha256[:], srcStampHash) {
		v.res.Errors = append(v.res.Errors, &SourceStampCertMismatchError{
			hex.EncodeToString(certSha256[:]),
			hex.EncodeToString(srcStampHash),
		})
		return
	}

	signedSignatureSchemes, err := getLenghtPrefixedSlice(stampBlock)
	if err != nil {
		v.addError("failed to parse signedSignatureSchemes: %s", err.Error())
		return
	}

	neededSchemeId := v.verifiedSchemeId
	if neededSchemeId == schemeIdV31 {
		neededSchemeId = schemeIdV3
	}

	var signaturesBlock *bytes.Buffer
	for signedSignatureSchemes.Len() != 0 {
		schemeBuf, err := getLenghtPrefixedSlice(signedSignatureSchemes)
		if err != nil {
			v.addError("failed to parse signedSignature scheme: %s", err.Error())
			return
		}

		var schemeId int32
		if err := binary.Read(schemeBuf, binary.LittleEndian, &schemeId); err != nil {
			v.addError("failed to parse signedSignature scheme int: %s", err.Error())
			return
		}

		apkDigestSignatures, err := getLenghtPrefixedSlice(schemeBuf)
		if err != nil {
			v.addError("failed to parse apkDigestSignatures: %s", err.Error())
			return
		}

		if schemeId == neededSchemeId {
			signaturesBlock = apkDigestSignatures
		}
	}

	if signaturesBlock == nil {
		v.addError("No source stamp signature for scheme %d", neededSchemeId)
		return
	}

	digest := v.getSignatureSchemeDigest(contentDigests)

	if !v.verifySignature(digest, v.res.Cert, signaturesBlock) {
		return
	}

	if stampBlock.Len() != 0 {
		attributeData, err := getLenghtPrefixedSlice(stampBlock)
		if err != nil {
			v.addError("failed to parse attributeData: %s", err.Error())
			return
		}

		attributeSignatures, err := getLenghtPrefixedSlice(stampBlock)
		if err != nil {
			v.addError("failed to parse attributeSignatures: %s", err.Error())
			return
		}

		if !v.verifySignature(attributeData.Bytes(), v.res.Cert, attributeSignatures) {
			return
		}

		if !v.parseStampAttributes(attributeData) {
			return
		}
	}
	return
}

func (v *sourceStampVerifier) getSignatureSchemeDigest(contentDigests map[contentDigest][]byte) []byte {
	type digestPair struct {
		algo contentDigest
		data []byte
	}

	pairs := make([]digestPair, 0, len(contentDigests))
	for algo, data := range contentDigests {
		pairs = append(pairs, digestPair{algo, data})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].algo < pairs[j].algo
	})

	// encodeAsSequenceOfLengthPrefixedPairsOfIntAndLengthPrefixedBytes
	var buf bytes.Buffer
	for _, p := range pairs {
		binary.Write(&buf, binary.LittleEndian, int32(8+len(p.data)))
		binary.Write(&buf, binary.LittleEndian, int32(p.algo))
		binary.Write(&buf, binary.LittleEndian, int32(len(p.data)))
		buf.Write(p.data)
	}
	return buf.Bytes()
}

func (v *sourceStampVerifier) verifySignature(data []byte, cert *x509.Certificate, signatures *bytes.Buffer) bool {
	supportedSignatures := map[SignatureAlgorithm][]byte{}
	for count := 0; signatures.Len() != 0; count++ {
		signature, err := getLenghtPrefixedSlice(signatures)
		if err != nil {
			v.addError("failed to parse signature %d: %s", count, err.Error())
			return false
		}

		var sigAlgoId SignatureAlgorithm
		if err := binary.Read(signature, binary.LittleEndian, &sigAlgoId); err != nil {
			v.addError("failed to parse signature algo id %d: %s", count, err.Error())
			return false
		}

		if !sigAlgoId.isSupported() {
			v.addWarning("unsupported signature algo %d: 0x%04x", count, int32(sigAlgoId))
			continue
		}

		sigBytes, err := getLenghtPrefixedSlice(signature)
		if err != nil {
			v.addError("failed to parse signature bytes %d: %s", count, err.Error())
			return false
		}
		supportedSignatures[sigAlgoId] = sigBytes.Bytes()
	}

	if len(supportedSignatures) == 0 {
		v.addError("no supported signatures")
		return false
	}

	toVerify := v.filterSignaturesToVerify(supportedSignatures)
	if len(toVerify) == 0 {
		v.addError("no supported signatures for verification")
		return false
	}

	for algo, sigBytes := range toVerify {
		if err := verifySignature(cert.PublicKey, algo, data, sigBytes); err != nil {
			v.addError("failed to verify signature %s: %s", algo.String(), err.Error())
			return false
		}
	}

	return true
}

func (v *sourceStampVerifier) filterSignaturesToVerify(signatures map[SignatureAlgorithm][]byte) map[SignatureAlgorithm][]byte {
	type signature struct {
		algo SignatureAlgorithm
		data []byte
	}

	candidates := map[int32]*signature{}

	var minProvidedSignaturesVersion int32 = math.MaxInt32
	for algo, content := range signatures {
		sigMin := algo.getMinSdkVersionJca()
		if sigMin > v.maxSdkVersion {
			continue
		}
		if sigMin < minProvidedSignaturesVersion {
			minProvidedSignaturesVersion = sigMin
		}

		if c := candidates[sigMin]; c == nil || compareAlgos(algo, c.algo) > 0 {
			candidates[sigMin] = &signature{algo, content}
		}
	}

	if v.minSdkVersion < minProvidedSignaturesVersion || len(candidates) == 0 {
		var present []string
		for algo := range signatures {
			present = append(present, algo.String())
		}

		if len(candidates) == 0 {
			v.addError("no supported signature, present: %s", strings.Join(present, ","))
		} else {
			v.addError("minimum provided signature version %d > minimum sdk version %d, present: %s",
				minProvidedSignaturesVersion, v.minSdkVersion, strings.Join(present, ","))
		}
		return nil
	}

	res := make(map[SignatureAlgorithm][]byte, len(candidates))
	for _, c := range candidates {
		res[c.algo] = c.data
	}
	return res
}

func (v *sourceStampVerifier) parseCertificate(block *bytes.Buffer) *x509.Certificate {
	encodedCert, err := getLenghtPrefixedSlice(block)
	if err != nil {
		v.addError("failed to parse encoded cert: %s", err.Error())
		return nil
	}

	cert, err := x509andr.ParseCertificateForGo(encodedCert.Bytes())
	if err != nil {
		v.addError("failed to parse certificate: %s", err.Error())
		return nil
	}

	return cert
}

func (v *sourceStampVerifier) parseStampAttributes(attributesData *bytes.Buffer) bool {
	attributes, err := getLenghtPrefixedSlice(attributesData)
	if err != nil {
		v.addError("failed to parse attributes from: %s", err.Error())
		return false
	}

	for count := 0; attributes.Len() != 0; count++ {
		attr, err := getLenghtPrefixedSlice(attributes)
		if err != nil {
			v.addError("failed to parse single attribute %d: %s", count, err.Error())
			return false
		}

		var id uint32
		if err := binary.Read(attr, binary.LittleEndian, &id); err != nil {
			v.addError("failed to parse attribute ID %d: %s", count, err.Error())
			return false
		}

		switch id {
		case sourceStampAttrProofOfRotation:
			v.res.Lineage = v.parseStampCertLineage(attr)
			if v.res.Lineage == nil {
				return false
			}
			if !v.res.Cert.Equal(v.res.Lineage[len(v.res.Lineage)-1].Cert) {
				v.addError("lineage certificate mismatch")
				return false
			}
		case sourceStampAttrTime:
			var timestamp int64
			if err := binary.Read(attr, binary.LittleEndian, &timestamp); err != nil {
				v.addError("failed to parse attribute timestamp: %s", err.Error())
				return false
			}

			if timestamp <= 0 {
				v.addError("Invalid timestamp %d in source stamp", timestamp)
			} else {
				v.res.SigningTime = time.Unix(timestamp, 0)
			}
		default:
			v.addWarning("Unknown attribute 0x%08x is present.", id)
		}
	}
	return true
}

func (v *sourceStampVerifier) parseStampCertLineage(lineage *bytes.Buffer) []*SourceStampLineageNode {
	var version int32
	if err := binary.Read(lineage, binary.LittleEndian, &version); err != nil {
		v.addError("failed to parse lineage version: %s", err.Error())
		return nil
	}

	if version != 1 {
		v.addError("unsupported lineage version: %d", version)
		return nil
	}

	var lastCert *x509.Certificate
	var lastSigAlgo SignatureAlgorithm
	var certHistory []*x509.Certificate
	var res []*SourceStampLineageNode
	for count := 0; lineage.Len() != 0; count++ {
		node, err := getLenghtPrefixedSlice(lineage)
		if err != nil {
			v.addError("failed to parse lineage node %d: %s", count, err.Error())
			return nil
		}

		signedData, err := getLenghtPrefixedSlice(node)
		if err != nil {
			v.addError("failed to parse lineage signedData %d: %s", count, err.Error())
			return nil
		}

		var flags int32
		if err := binary.Read(node, binary.LittleEndian, &flags); err != nil {
			v.addError("failed to parse lineage node %d flags: %s", count, err.Error())
			return nil
		}

		var sigAlgo SignatureAlgorithm
		if err := binary.Read(node, binary.LittleEndian, &sigAlgo); err != nil {
			v.addError("failed to parse lineage node %d sigAlgo: %s", count, err.Error())
			return nil
		}

		signature, err := getLenghtPrefixedSlice(node)
		if err != nil {
			v.addError("failed to parse lineage signature %d: %s", count, err.Error())
			return nil
		}

		if lastCert != nil {
			if err := verifySignature(lastCert.PublicKey, lastSigAlgo, signedData.Bytes(), signature.Bytes()); err != nil {
				v.addError("failed to verify lineage signature %d: %s", count, err.Error())
				return nil
			}
		}

		encodedCert, err := getLenghtPrefixedSlice(signedData)
		if err != nil {
			v.addError("failed to parse lineage encodedCert %d: %s", count, err.Error())
			return nil
		}

		var signedSigAlgorithm SignatureAlgorithm
		if err := binary.Read(signedData, binary.LittleEndian, &signedSigAlgorithm); err != nil {
			v.addError("failed to parse lineage node %d signedSigAlgorithm: %s", count, err.Error())
			return nil
		}

		if lastCert != nil && lastSigAlgo != signedSigAlgorithm {
			v.addError("signing algo mismatch for node %d: %s vs %s", count, lastSigAlgo, signedSigAlgorithm)
			return nil
		}

		lastCert, err = x509andr.ParseCertificateForGo(encodedCert.Bytes())
		if err != nil {
			v.addError("failed to parse lineage cert at node %d: %s", count, err.Error())
			return nil
		}

		for i := range certHistory {
			if certHistory[i].Equal(lastCert) {
				v.addError("encountered duplicate entries in lineage")
				return nil
			}
		}
		certHistory = append(certHistory, lastCert)
		lastSigAlgo = sigAlgo

		res = append(res, &SourceStampLineageNode{
			Cert:       lastCert,
			ParentAlgo: signedSigAlgorithm,
			Algo:       sigAlgo,
			Signature:  signature.Bytes(),
			Flags:      flags,
		})
	}
	return res
}

func (v *sourceStampVerifier) getSchemeV1ContentDigest(zip *apkparser.ZipReader) []byte {
	manifestFile := zip.File["META-INF/MANIFEST.MF"]
	if manifestFile == nil {
		return []byte{}
	}

	manBytes, err := manifestFile.ReadAll(128 * 1024 * 1024)
	if err != nil {
		v.addWarning("failed to read %s: %s", manifestFile.Name, err.Error())
		return []byte{}
	}

	hash := sha256.Sum256(manBytes)
	return hash[:]
}
