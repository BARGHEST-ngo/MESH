package signingblock

import (
	"bytes"
	"github.com/avast/apkverifier/apilevel"
	"math"
)

type schemeV31 struct {
	backendV3  schemeV3
	backendV31 schemeV3
}

func (s *schemeV31) parseSigners(block *bytes.Buffer, contentDigests map[contentDigest][]byte, result *VerificationResult) {
	v31ContentDigest := make(map[contentDigest][]byte)
	s.backendV31.parseSigners(block, v31ContentDigest, result)

	s.backendV3.rotationMinSdkVersion = math.MaxInt32
	for _, signer := range s.backendV31.signers {
		if signer.minSdkVersion < s.backendV3.rotationMinSdkVersion {
			s.backendV3.rotationMinSdkVersion = signer.minSdkVersion
		}
	}
	if s.backendV3.rotationMinSdkVersion == math.MaxInt32 {
		s.backendV3.rotationMinSdkVersion = 0
	}

	v3Block := result.ExtraBlocks[blockIdSchemeV3]
	if v3Block == nil {
		result.addError("The APK contains a v3.1 signing block without a v3.0 base block")
		for k, v := range v31ContentDigest {
			contentDigests[k] = v
		}
		return
	}

	if result.ExtraResults == nil {
		result.ExtraResults = make(map[int]*VerificationResult)
	}
	v3result := &VerificationResult{
		SchemeId: schemeIdV3,
	}
	result.ExtraResults[schemeIdV3] = v3result

	s.backendV3.parseSigners(bytes.NewBuffer(v3Block), contentDigests, v3result)
}

func (s *schemeV31) finalizeResult(requestedMinSdkVersion, requestedMaxSdkVersion int32, result *VerificationResult) {
	v31minSdk := requestedMinSdkVersion
	if v31minSdk < apilevel.V13_0_TIRAMISU {
		v31minSdk = apilevel.V13_0_TIRAMISU
	}
	s.backendV31.finalizeResult(v31minSdk, requestedMaxSdkVersion, result)

	if v3Result := result.ExtraResults[schemeIdV3]; v3Result != nil {
		s.backendV3.finalizeResult(requestedMinSdkVersion, requestedMaxSdkVersion, v3Result)
		result.Warnings = append(result.Warnings, v3Result.Warnings...)
		result.Errors = append(result.Errors, v3Result.Errors...)
	}
}
