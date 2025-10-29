package signingblock

import (
	"bytes"
	"encoding/binary"
	"github.com/avast/apkverifier/apilevel"
)

const (
	attrV2StrippingProtection = 0xbeeff00d
)

type schemeV2 struct {
	minSdkVersion, maxSdkVersion int32
}

func (s *schemeV2) parseSigners(block *bytes.Buffer, contentDigests map[contentDigest][]byte, result *VerificationResult) {
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

	if len(result.Certs) > maxApkSigners {
		result.addError("APK Signature Scheme V2 only supports a maximum of %d signers, found %d",
			maxApkSigners, len(result.Certs))
	}
}

func (s *schemeV2) finalizeResult(minSdkVersion, maxSdkVersion int32, result *VerificationResult) {

}

func (s *schemeV2) verifySigner(signerBlock *bytes.Buffer, contentDigests map[contentDigest][]byte, result *VerificationResult) {
	signedData, err := getLenghtPrefixedSlice(signerBlock)
	if err != nil {
		result.addError("failed to read signed data: %s", err.Error())
		return
	}
	signedDataBytes := signedData.Bytes()

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

	// Parse certificates
	if _, success := ctx.parseCertificates(certificatesSlice); !success {
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
		case attrV2StrippingProtection:
			if !apilevel.SupportsSigV3(s.maxSdkVersion) {
				break
			}

			var strippedSchemeId int32
			if err := binary.Read(attribute, binary.LittleEndian, &strippedSchemeId); err != nil {
				result.addError("failed to read additional attribute %d's strippedSchemeId: %s", additionalAttributeCount, err.Error())
				return
			}

			strippedBlockId, prs := schemeIdToBlockId[strippedSchemeId]
			if !prs {
				result.addError("unknown stripped scheme id: %d", strippedSchemeId)
			} else if result.ExtraBlocks[strippedBlockId] == nil {
				result.addError("this apk was signed with %d signing scheme, but it was stripped, downgrade attack?", strippedSchemeId)
				return
			}
		default:
			result.addWarning("unknown additional attribute id 0x%x", uint32(id))
		}
	}
}
