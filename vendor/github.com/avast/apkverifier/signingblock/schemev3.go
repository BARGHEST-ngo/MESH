package signingblock

// https://android.googlesource.com/platform/tools/apksig/+/master/src/main/java/com/android/apksig

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/avast/apkverifier/apilevel"
)

const (
	attrV3ProofOfRotation = 0x3ba06f8c
	attrV3MinSdkVersion   = 0x559f8b02
	attrV3RotationOnDev   = 0xc2a6b3ba

	lineageVersionFirst = 1
)

type schemeV3 struct {
	actingAsV31                  bool
	minSdkVersion, maxSdkVersion int32
	sourceStampBlock             []byte
	signers                      []*schemeV3SignerInfo
	rotationMinSdkVersion        int32
}

type schemeV3SignerInfo struct {
	lineage                      *V3SigningLineage
	minSdkVersion, maxSdkVersion int32
}

func (s *schemeV3) parseSigners(block *bytes.Buffer, contentDigests map[contentDigest][]byte, result *VerificationResult) {
	signers, err := getLenghtPrefixedSlice(block)
	if err != nil {
		result.addError("failed to read list of signers: %s", err.Error())
		return
	}

	signerCount := 0
	for signers.Len() > 0 {
		signerCount++

		signer, err := getLenghtPrefixedSlice(signers)
		if err != nil {
			result.addError("failed to parse/verify signer #%d block: %s", signerCount, err.Error())
			return
		}

		s.verifySigner(signer, contentDigests, result)
	}
}

func (s *schemeV3) finalizeResult(requestedMinSdkVersion, requestedMaxSdkVersion int32, result *VerificationResult) {
	// v3 didn't exist prior to P, so make sure that we're only judging v3 on its supported
	// platforms
	if requestedMinSdkVersion < apilevel.V9_0_Pie {
		requestedMinSdkVersion = apilevel.V9_0_Pie
	}

	if s.actingAsV31 {
		// V3.1 supports targeting an SDK version later than that of the initial release
		// in which it is supported; allow any range for V3.1 as long as V3.0 covers the
		// rest of the range.
		requestedMinSdkVersion = requestedMaxSdkVersion
	}

	sort.Slice(s.signers, func(i, j int) bool {
		return s.signers[i].minSdkVersion < s.signers[j].minSdkVersion
	})

	var firstMin, lastMax int32
	lastLineageSize := 0
	lineages := make([]*V3SigningLineage, 0, len(s.signers))
	for _, signer := range s.signers {
		currentMin, currentMax := signer.minSdkVersion, signer.maxSdkVersion
		if firstMin == 0 {
			firstMin = currentMin
		} else if currentMin != lastMax+1 {
			result.addError("inconsistent signer sdkversions")
			break
		}
		lastMax = currentMax

		if signer.lineage != nil {
			if size := len(signer.lineage.Nodes); size < lastLineageSize {
				result.addError("inconsistent signer lineage sizes")
				break
			} else {
				lastLineageSize = size
				lineages = append(lineages, signer.lineage)
			}
		}
	}

	if s.rotationMinSdkVersion != 0 {
		requestedMaxSdkVersion = s.rotationMinSdkVersion - 1
	}

	if firstMin > requestedMinSdkVersion || lastMax < requestedMaxSdkVersion {
		result.addError("missing sdk versions, supports only <%d;%d>, got range (%d;%d)", firstMin, lastMax, requestedMinSdkVersion, requestedMaxSdkVersion)
	}

	if result.SigningLineage == nil {
		var err error
		result.SigningLineage, err = s.consolidateLineages(lineages)
		if err != nil {
			result.addError(err.Error())
		}
	}
}

func (s *schemeV3) verifySigner(signerBlock *bytes.Buffer, contentDigests map[contentDigest][]byte, result *VerificationResult) {
	signedData, err := getLenghtPrefixedSlice(signerBlock)
	if err != nil {
		result.addError("failed to read signed data: %s", err.Error())
		return
	}
	signedDataBytes := signedData.Bytes()

	// Parse min/max sdk
	var parsedMinSdkVersion, parsedMaxSdkVersion int32
	if err := binary.Read(signerBlock, binary.LittleEndian, &parsedMinSdkVersion); err != nil {
		result.addError("failed to read parsedMinSdkVersion: %s", err.Error())
		return
	}
	if err := binary.Read(signerBlock, binary.LittleEndian, &parsedMaxSdkVersion); err != nil {
		result.addError("failed to read parsedMaxSdkVersion: %s", err.Error())
		return
	}

	if parsedMinSdkVersion < 1 || (parsedMinSdkVersion > parsedMaxSdkVersion) {
		result.addError("invalid min/max sdk versions: <%d,%d>", parsedMinSdkVersion, parsedMaxSdkVersion)
	}

	signerInfo := &schemeV3SignerInfo{
		minSdkVersion: parsedMinSdkVersion,
		maxSdkVersion: parsedMaxSdkVersion,
	}
	s.signers = append(s.signers, signerInfo)

	ctx := signerContext{result: result}

	// Parse signatures
	signaturesSlice, err := getLenghtPrefixedSlice(signerBlock)
	if err != nil {
		result.addError("failed to read signatures: %s", err.Error())
		return
	}

	if !ctx.parseSignatures(signaturesSlice) {
		return
	}

	// Parse & verify public key
	publicKeySlice, err := getLenghtPrefixedSlice(signerBlock)
	if err != nil {
		result.addError("failed to read publicKeySlice: %s", err.Error())
		return
	}

	if !ctx.parsePublicKey(publicKeySlice, signedDataBytes) {
		return
	}

	// Get digests
	digestsSlice, err := getLenghtPrefixedSlice(signedData)
	if err != nil {
		result.addError("failed to read digests from signedData: %s", err.Error())
		return
	}

	// Get certs
	certificatesSlice, err := getLenghtPrefixedSlice(signedData)
	if err != nil {
		result.addError("failed to read certificates from signedData: %s", err.Error())
		return
	}

	// Parse/verify signedMin/MaxSdkVersion
	var signedMinSdkVersion, signedMaxSdkVersion int32
	if err := binary.Read(signedData, binary.LittleEndian, &signedMinSdkVersion); err != nil {
		result.addError("failed to read signedMinSdkVersion: %s", err.Error())
		return
	}
	if err := binary.Read(signedData, binary.LittleEndian, &signedMaxSdkVersion); err != nil {
		result.addError("Failed to read signedMaxSdkVersion: %s", err.Error())
		return
	}

	if parsedMinSdkVersion != signedMinSdkVersion {
		result.addError("mismatch between parsed and signed minSdkVersion: %d != %d", parsedMinSdkVersion, signedMinSdkVersion)
	}
	if parsedMaxSdkVersion != parsedMaxSdkVersion {
		result.addError("mismatch between parsed and signed minSdkVersion: %d != %d", parsedMinSdkVersion, signedMinSdkVersion)
	}

	// Parse certificates
	mainCert, success := ctx.parseCertificates(certificatesSlice)
	if !success {
		return
	}

	// Parse digests
	if !ctx.parseDigests(digestsSlice, contentDigests) {
		return
	}

	// Parse additional attributes
	additionalAttributes, err := getLenghtPrefixedSlice(signedData)
	if err != nil {
		result.addError("failed to read additionalAttributes from signedData: %s", err.Error())
		return
	}

	attrV3MinSdkVersionFound := false

	additionalAttributeCount := 0
	for additionalAttributes.Len() > 0 {
		additionalAttributeCount++

		attribute, err := getLenghtPrefixedSlice(additionalAttributes)
		if err != nil {
			result.addError("failed to read additional attribute %d: %s", additionalAttributeCount, err.Error())
			return
		}

		var id uint32
		if err := binary.Read(attribute, binary.LittleEndian, &id); err != nil {
			result.addError("failed to read additional attribute %d's id: %s", additionalAttributeCount, err.Error())
			return
		}

		switch id {
		case attrV3ProofOfRotation:
			nodes, err := s.readSigningCertificateLineage(attribute)
			if err != nil {
				result.addError("failed to read signing certificate lineage attribute: %s", err.Error())
				return
			}

			if len(nodes) != 0 {
				signerInfo.lineage = &V3SigningLineage{
					MinSdkVersion: lineageCalculateMinSdkVersion(nodes),
					Nodes:         nodes,
				}

				subLineage, err := signerInfo.lineage.getSubLineage(mainCert)
				if err != nil {
					result.addError(err.Error())
					break
				}

				if len(signerInfo.lineage.Nodes) != len(subLineage.Nodes) {
					result.addError("sub lineage cert mismatch")
				}
			}
		case attrV3MinSdkVersion:
			attrV3MinSdkVersionFound = true
			if apilevel.SupportsSigV31(s.maxSdkVersion) {
				var currentAttrMinSdk int32
				if err := binary.Read(attribute, binary.LittleEndian, &currentAttrMinSdk); err != nil {
					result.addError("failed to read attrRotationMinSdkVersion: %s", err.Error())
					return
				}

				if s.rotationMinSdkVersion != 0 {
					if s.rotationMinSdkVersion != currentAttrMinSdk {
						result.addError("The v3 signer indicates key rotation should be supported starting from SDK version %d, but the v3.1 block targets %d for rotation",
							s.rotationMinSdkVersion, currentAttrMinSdk)
					}
				} else {
					result.addError("The v3 signer indicates key rotation should be supported starting from SDK version %d, but a v3.1 block was not found", currentAttrMinSdk)
				}
			} else {
				result.addWarning("unknown additional attribute id 0x%x", id)
			}
		case attrV3RotationOnDev:
			// This attribute should only be used by a v3.1 signer to indicate rotation
			// is targeting the development release that is using the SDK version of the
			// previously released platform version.
			if !s.actingAsV31 {
				result.addWarning("The rotation-targets-dev-release attribute is only supported on v3.1 signers; this attribute will be ignored by the platform in a v3.0 signer")
			}
		default:
			result.addWarning("unknown additional attribute id 0x%x", id)
		}
	}

	if s.rotationMinSdkVersion != 0 && !attrV3MinSdkVersionFound {
		result.addWarning("APK supports key rotation starting from SDK version %d, but the v3 signer does not "+
			"contain the attribute to detect if this signature is stripped", s.rotationMinSdkVersion)
	}
}

func (s *schemeV3) readSigningCertificateLineage(lineageSlice *bytes.Buffer) (V3LineageSigningCertificateNodeList, error) {
	if lineageSlice.Len() == 0 {
		return nil, nil
	}

	var version int32
	if err := binary.Read(lineageSlice, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %s", err.Error())
	}

	if version != lineageVersionFirst {
		return nil, fmt.Errorf("unknown lineage version %d", version)
	}

	var nodeCount int
	var lastCert *x509.Certificate
	var lastSigAlgorithmId SignatureAlgorithm
	var certHistory []*x509.Certificate
	var result V3LineageSigningCertificateNodeList
	for lineageSlice.Len() != 0 {
		nodeCount++

		nodeBytes, err := getLenghtPrefixedSlice(lineageSlice)
		if err != nil {
			return nil, fmt.Errorf("failed to read nodeBytes for node %d: %s", nodeCount, err.Error())
		}

		signedData, err := getLenghtPrefixedSlice(nodeBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read signedData for node %d: %s", nodeCount, err.Error())
		}

		var flags int32
		if err := binary.Read(nodeBytes, binary.LittleEndian, &flags); err != nil {
			return nil, fmt.Errorf("failed to read flags for node %d: %s", nodeCount, err.Error())
		}

		var sigAlgorithmId SignatureAlgorithm
		if err := binary.Read(nodeBytes, binary.LittleEndian, &sigAlgorithmId); err != nil {
			return nil, fmt.Errorf("failed to read sigAlgorithmId for node %d: %s", nodeCount, err.Error())
		}

		signature, err := getLenghtPrefixedSlice(nodeBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read signature for node %d: %s", nodeCount, err.Error())
		}

		if lastCert != nil {
			if err := verifySignature(lastCert.PublicKey, lastSigAlgorithmId, signedData.Bytes(), signature.Bytes()); err != nil {
				return nil, fmt.Errorf("unable to verify signature of certificate #%d using algo %d: %s",
					nodeCount, lastSigAlgorithmId, err.Error())
			}
		}

		encodedCert, err := getLenghtPrefixedSlice(signedData)
		if err != nil {
			return nil, fmt.Errorf("failed to read encodedCert for node %d: %s", nodeCount, err.Error())
		}

		var signedSigAlgorithm SignatureAlgorithm
		if err := binary.Read(signedData, binary.LittleEndian, &signedSigAlgorithm); err != nil {
			return nil, fmt.Errorf("failed to read signedSigAlgorithm for node %d: %s", nodeCount, err.Error())
		}

		if lastCert != nil && lastSigAlgorithmId != signedSigAlgorithm {
			return nil, fmt.Errorf("signing algorithm ID mismatch for certificate #%d", nodeCount)
		}

		lastCert, err = x509.ParseCertificate(encodedCert.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to parse cert in node #%d", nodeCount)
		}

		for _, histCert := range certHistory {
			if histCert.Equal(lastCert) {
				return nil, fmt.Errorf("encountered duplicate entries at node #%d, signing certificates should be unique.", nodeCount)
			}
		}

		certHistory = append(certHistory, lastCert)
		lastSigAlgorithmId = sigAlgorithmId

		result = append(result, &V3LineageSigningCertificateNode{
			SigningCert:        lastCert,
			ParentSigAlgorithm: signedSigAlgorithm,
			SigAlgorithm:       sigAlgorithmId,
			Signature:          signature.Bytes(),
			Flags:              LineageCertCaps(flags),
		})
	}

	return result, nil
}

func (s *schemeV3) consolidateLineages(lineages []*V3SigningLineage) (*V3SigningLineage, error) {
	if len(lineages) == 0 {
		return nil, nil
	}

	largestIndex := 0
	maxSize := 0
	for i := range lineages {
		if s := len(lineages[i].Nodes); s > maxSize {
			maxSize = s
			largestIndex = i
		}
	}

	largestList := lineages[largestIndex].Nodes
	for i := range lineages {
		if i == largestIndex {
			continue
		}

		underTest := lineages[i].Nodes
		if !underTest.Equal(largestList[:len(underTest)]) {
			return nil, errors.New("inconsistent certificate lineage, not all lineages are subsets of each other")
		}
	}
	return lineages[largestIndex], nil
}
