package signingblock

import (
	"crypto/x509"
	"fmt"
)

type VerificationResult struct {
	Certs          [][]*x509.Certificate
	SchemeId       int
	SigningLineage *V3SigningLineage

	// When APK is signed with v3.1, the v3 result is stored here. Any v3 errors are lifted to the main Warnings/errors though
	ExtraResults map[int]*VerificationResult

	Frosting *FrostingResult

	SourceStamp *SourceStampResult

	// Extra blocks found in the signing block that are not used by apkverifier,
	// either completely unknown, or those found in BlockId constants.
	// Parsed block types (schemeV2, V3, play frosting..) will NOT be in this map.
	// May be nil.
	ExtraBlocks map[BlockId][]byte

	Warnings []string
	Errors   []error
}

type FrostingResult struct {
	Error        error
	KeySha256    string
	ProtobufInfo []byte
}

type certAdder struct {
	Certs []*x509.Certificate

	res *VerificationResult
}

func (r *VerificationResult) addWarning(format string, args ...interface{}) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

func (r *VerificationResult) addError(format string, args ...interface{}) {
	r.Errors = append(r.Errors, fmt.Errorf(format, args...))
}

func (r *VerificationResult) ContainsErrors() bool {
	return len(r.Errors) != 0
}

func (r *VerificationResult) GetLastError() error {
	if l := len(r.Errors); l != 0 {
		return r.Errors[l-1]
	}
	return nil
}

func (r *VerificationResult) getCertAdder() certAdder {
	return certAdder{
		res: r,
	}
}

func (a *certAdder) append(cert *x509.Certificate) {
	a.Certs = append(a.Certs, cert)
	if len(a.Certs) == 1 {
		a.res.Certs = append(a.res.Certs, a.Certs)
	} else {
		idx := len(a.res.Certs) - 1
		a.res.Certs[idx] = a.Certs
	}
}
