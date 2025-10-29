package signingblock

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/avast/apkverifier/internal/x509andr"

	"github.com/avast/apkverifier/apilevel"

	"crypto/x509"
)

type SignatureAlgorithm int32

const (
	SigRsaPssWithSha256            SignatureAlgorithm = 0x0101
	SigRsaPssWithSha512                               = 0x0102
	SigRsaPkcs1V15WithSha256                          = 0x0103
	SigRsaPkcs1V15WithSha512                          = 0x0104
	SigEcdsaWithSha256                                = 0x0201
	SigEcdsaWithSha512                                = 0x0202
	SigDsaWithSha256                                  = 0x0301
	SigVerityRsaPkcs1V15WithSha256                    = 0x0421
	SigVerityEcdsaWithSha256                          = 0x0423
	SigVerityDsaWithSha256                            = 0x425
)

type contentDigest int

const (
	digestChunkedSha256    = 1
	digestChunkedSha512    = 2
	digestVeritySha256     = 3
	digestNonChunkedSha256 = 4
)

func (c contentDigest) Hash() crypto.Hash {
	switch c {
	case digestChunkedSha256:
		return crypto.SHA256
	case digestChunkedSha512:
		return crypto.SHA512
	case digestVeritySha256:
		return crypto.SHA256
	case digestNonChunkedSha256:
		return crypto.SHA256
	default:
		panic(fmt.Sprintf("Unknown contentDigest: %d", c))
	}
}

func (algo SignatureAlgorithm) String() string {
	switch algo {
	case SigRsaPssWithSha256:
		return "SigRsaPssWithSha256"
	case SigRsaPssWithSha512:
		return "SigRsaPssWithSha512"
	case SigRsaPkcs1V15WithSha256:
		return "SigRsaPkcs1V15WithSha256"
	case SigRsaPkcs1V15WithSha512:
		return "SigRsaPkcs1V15WithSha512"
	case SigEcdsaWithSha256:
		return "SigEcdsaWithSha256"
	case SigEcdsaWithSha512:
		return "SigEcdsaWithSha512"
	case SigDsaWithSha256:
		return "SigDsaWithSha256"
	case SigVerityRsaPkcs1V15WithSha256:
		return "SigVerityRsaPkcs1V15WithSha256"
	case SigVerityEcdsaWithSha256:
		return "SigVerityEcdsaWithSha256"
	case SigVerityDsaWithSha256:
		return "SigVerityDsaWithSha256"
	}
	return fmt.Sprintf("0x%04x", uint32(algo))
}

func (algo SignatureAlgorithm) isSupported() bool {
	switch algo {
	case SigRsaPssWithSha256, SigRsaPssWithSha512,
		SigRsaPkcs1V15WithSha256, SigRsaPkcs1V15WithSha512,
		SigEcdsaWithSha256, SigEcdsaWithSha512,
		SigDsaWithSha256,
		SigVerityRsaPkcs1V15WithSha256, SigVerityEcdsaWithSha256, SigVerityDsaWithSha256:
		return true
	default:
		return false
	}
}

func (algo SignatureAlgorithm) getDigestType() contentDigest {
	switch algo {
	case SigRsaPssWithSha256, SigRsaPkcs1V15WithSha256, SigEcdsaWithSha256, SigDsaWithSha256:
		return digestChunkedSha256
	case SigVerityRsaPkcs1V15WithSha256, SigVerityEcdsaWithSha256, SigVerityDsaWithSha256:
		return digestVeritySha256
	case SigRsaPssWithSha512, SigRsaPkcs1V15WithSha512, SigEcdsaWithSha512:
		return digestChunkedSha512
	default:
		panic(fmt.Sprintf("Unknown signature algorithm 0x%x", algo))
	}
}

func (algo SignatureAlgorithm) getMinSdkVersion() int32 {
	switch algo {
	case SigRsaPssWithSha512, SigRsaPkcs1V15WithSha256, SigRsaPkcs1V15WithSha512,
		SigEcdsaWithSha256, SigEcdsaWithSha512, SigDsaWithSha256:
		return apilevel.V7_0_Nougat
	case SigVerityRsaPkcs1V15WithSha256, SigVerityEcdsaWithSha256, SigVerityDsaWithSha256:
		return apilevel.V9_0_Pie
	default:
		return math.MaxInt32
	}
}

/*
const (

	SigRsaPssWithSha256            SignatureAlgorithm = 0x0101
	SigRsaPssWithSha512                               = 0x0102
	SigRsaPkcs1V15WithSha256                          = 0x0103
	SigRsaPkcs1V15WithSha512                          = 0x0104
	SigEcdsaWithSha256                                = 0x0201
	SigEcdsaWithSha512                                = 0x0202
	SigDsaWithSha256                                  = 0x0301
	SigVerityRsaPkcs1V15WithSha256                    = 0x0421
	SigVerityEcdsaWithSha256                          = 0x0423
	SigVerityDsaWithSha256                            = 0x425

)
*/
func (algo SignatureAlgorithm) getMinSdkVersionJca() int32 {
	switch algo {
	case SigRsaPssWithSha256, SigRsaPssWithSha512:
		return apilevel.V6_0_Marshmallow
	case SigRsaPkcs1V15WithSha256, SigRsaPkcs1V15WithSha512, SigDsaWithSha256,
		SigVerityRsaPkcs1V15WithSha256, SigVerityDsaWithSha256:
		return apilevel.V1_0_InitialRelease
	case SigEcdsaWithSha256, SigEcdsaWithSha512, SigVerityEcdsaWithSha256:
		return apilevel.V3_0_Honeycomb
	default:
		return math.MaxInt32
	}
}

type signerContext struct {
	result *VerificationResult

	bestAlgo          SignatureAlgorithm
	bestAlgoSignature []byte
	signaturesAlgos   []SignatureAlgorithm

	publicKeyBytes []byte
}

func (s *signerContext) parseSignatures(signaturesSlice *bytes.Buffer) (success bool) {
	s.bestAlgo = -1
	signatureCount := 0
	for signaturesSlice.Len() > 0 {
		signatureCount++

		signature, err := getLenghtPrefixedSlice(signaturesSlice)
		if err != nil {
			s.result.addError("failed to parse signature record #%d: %s", signatureCount, err.Error())
			return
		}

		if signature.Len() < 8 {
			s.result.addError("signature record %d is too short", signatureCount)
			return
		}

		var algo SignatureAlgorithm
		if err := binary.Read(signature, binary.LittleEndian, &algo); err != nil {
			s.result.addError("failed to parse signature record #%d: %s", signatureCount, err.Error())
			return
		}

		s.signaturesAlgos = append(s.signaturesAlgos, algo)
		if !algo.isSupported() {
			s.result.addWarning("signature %d is using unsupported algorithm %s", signatureCount, algo.String())
			continue
		}

		if s.bestAlgo == -1 || compareAlgos(algo, s.bestAlgo) > 0 {
			s.bestAlgo = algo
			sigBytes, err := getLenghtPrefixedSlice(signature)
			if err != nil {
				s.result.addError("failed to read signature bytes from signature record #%d: %s", signatureCount, err.Error())
				return
			}
			s.bestAlgoSignature = sigBytes.Bytes()
		}
	}

	if s.bestAlgo == -1 {
		if signatureCount == 0 {
			s.result.addError("no signatures found")
		} else {
			s.result.addError("no supported signatures found")
		}
		return
	}

	return true
}

func (s *signerContext) parsePublicKey(publicKeySlice *bytes.Buffer, signedDataBytes []byte) (success bool) {
	s.publicKeyBytes = publicKeySlice.Bytes()

	publicKey, err := x509andr.ParsePKIXPublicKey(s.publicKeyBytes)
	if err != nil {
		s.result.addError("failed to parse public key: %s", err.Error())
		return
	}

	err = verifySignature(publicKey, s.bestAlgo, signedDataBytes, s.bestAlgoSignature)
	if err != nil {
		s.result.addError("failed to verify signature of type 0x%x: %s", uint32(s.bestAlgo), err.Error())
		return
	}

	return true
}

func (s *signerContext) parseCertificates(certificatesSlice *bytes.Buffer) (mainCert *x509.Certificate, success bool) {
	certAdder := s.result.getCertAdder()
	certificateCount := 0
	for certificatesSlice.Len() > 0 {
		certificateCount++
		encodedCert, err := getLenghtPrefixedSlice(certificatesSlice)
		if err != nil {
			s.result.addError("failed to read certificate #%d: %s", certificateCount, err.Error())
			return
		}

		cert, err := x509andr.ParseCertificateForGo(encodedCert.Bytes())
		if err != nil {
			s.result.addError("Failed to parse certificate #%d: %s", certificateCount, err.Error())
			return
		}
		certAdder.append(cert)
	}

	if len(certAdder.Certs) == 0 {
		s.result.addError("No certificates listed.")
		return
	}

	mainCert = certAdder.Certs[0]
	if !bytes.Equal(mainCert.RawSubjectPublicKeyInfo, s.publicKeyBytes) {
		s.result.addError("Public key mismatch between certificate and signature record")
		return
	}

	return mainCert, true
}

func (s *signerContext) parseDigests(digestsSlice *bytes.Buffer, contentDigests map[contentDigest][]byte) (success bool) {
	var contentDigest []byte
	var digestSigAlgorithms []SignatureAlgorithm
	digestCount := 0
	for digestsSlice.Len() > 0 {
		digestCount++

		digest, err := getLenghtPrefixedSlice(digestsSlice)
		if err != nil {
			s.result.addError("failed to parse digest #%d: %s", digestCount, err.Error())
			return
		} else if digest.Len() < 8 {
			s.result.addError("failed to parse digest #%d: record too short", digestCount)
			return
		}

		var sigAlgorithm SignatureAlgorithm
		binary.Read(digest, binary.LittleEndian, &sigAlgorithm)
		digestSigAlgorithms = append(digestSigAlgorithms, sigAlgorithm)
		if sigAlgorithm == s.bestAlgo {
			cd, err := getLenghtPrefixedSlice(digest)
			if err != nil {
				s.result.addError("failed to read content digest for digest #%d: %s", digestCount, err.Error())
				return
			}
			contentDigest = cd.Bytes()
		}
	}

	algosEqual := len(digestSigAlgorithms) == len(s.signaturesAlgos)
	for i := 0; algosEqual && i < len(digestSigAlgorithms); i++ {
		algosEqual = digestSigAlgorithms[i] == s.signaturesAlgos[i]
	}

	if !algosEqual {
		s.result.addError("signature algorithms don't match between digests and signatures records")
		return
	}

	digestAlgorithm := s.bestAlgo.getDigestType()
	previousSignerDigest := contentDigests[digestAlgorithm]
	contentDigests[digestAlgorithm] = contentDigest
	if previousSignerDigest != nil && !bytes.Equal(previousSignerDigest, contentDigest) {
		s.result.addError("0x%x contents digest does not match the digest specified by a preceding signer", digestAlgorithm)
		return
	}

	return true
}

func compareAlgos(a, b SignatureAlgorithm) int {
	digest1 := a.getDigestType()
	digest2 := b.getDigestType()

	if digest1 == digest2 {
		return 0
	}

	switch digest1 {
	case digestChunkedSha256:
		switch digest2 {
		case digestChunkedSha512, digestVeritySha256:
			return -1
		default:
			panic(fmt.Sprintf("Unknown digest2: %d", digest2))
		}
	case digestChunkedSha512:
		switch digest2 {
		case digestChunkedSha256, digestVeritySha256:
			return 1
		default:
			panic(fmt.Sprintf("Unknown digest2: %d", digest2))
		}
	case digestVeritySha256:
		switch digest2 {
		case digestChunkedSha256:
			return 1
		case digestChunkedSha512:
			return -1
		default:
			panic(fmt.Sprintf("Unknown digest2: %d", digest2))
		}
	default:
		panic(fmt.Sprintf("Unknown digest1: %d", digest1))
	}
}
