package apkverifier

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/avast/apkverifier/signingblock"
	"math/big"
	"sort"
	"time"
)

// CertInfo Nicer looking certificate info
type CertInfo struct {
	Md5                string
	Sha1               string
	Sha256             string
	ValidFrom, ValidTo time.Time
	Issuer, Subject    string
	SignatureAlgorithm string
	SerialNumber       *big.Int
}

type byPreference [][]*x509.Certificate

func (c byPreference) Len() int      { return len(c) }
func (c byPreference) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c byPreference) Less(i, j int) bool {
	ci, cj := c[i][0], c[j][0]

	if ci.SignatureAlgorithm != cj.SignatureAlgorithm {
		return ci.SignatureAlgorithm > cj.SignatureAlgorithm
	}

	now := time.Now()
	// expired cert is "More" than not-expired one
	if ci.NotAfter.After(now) || ci.NotBefore.Before(now) {
		return false
	} else if cj.NotAfter.After(now) || cj.NotBefore.Before(now) {
		return true
	}

	if !ci.NotBefore.Equal(cj.NotBefore) {
		return ci.NotBefore.After(cj.NotBefore)
	}

	if !ci.NotAfter.Equal(cj.NotAfter) {
		return ci.NotAfter.Sub(ci.NotBefore) > cj.NotAfter.Sub(cj.NotBefore)
	}
	return bytes.Compare(ci.Raw, cj.Raw) > 0
}

// PickBestApkCert Picks the "best-looking" (most likely the correct one) certificate from the chain
// extracted from APK. Is noop for most APKs, as they usually contain only one certificate.
func PickBestApkCert(chains [][]*x509.Certificate) (*CertInfo, *x509.Certificate) {
	if len(chains) == 0 {
		return nil, nil
	}

	sort.Sort(byPreference(chains))

	return NewCertInfo(chains[0][0]), chains[0][0]
}

// NewCertInfo Returns new CertInfo with information from the x509.Certificate.
func NewCertInfo(cert *x509.Certificate) *CertInfo {
	var res CertInfo
	res.Fill(cert)
	return &res
}

// Fill Replaces CertInfo's data with information from the x509.Certificate.
func (ci *CertInfo) Fill(cert *x509.Certificate) {
	md5sum := md5.Sum(cert.Raw)
	sha1sum := sha1.Sum(cert.Raw)
	sha256sum := sha256.Sum256(cert.Raw)

	ci.Md5 = hex.EncodeToString(md5sum[:])
	ci.Sha1 = hex.EncodeToString(sha1sum[:])
	ci.Sha256 = hex.EncodeToString(sha256sum[:])
	ci.ValidFrom = cert.NotBefore
	ci.ValidTo = cert.NotAfter
	ci.Issuer = signingblock.PkixNameToString(&cert.Issuer)
	ci.Subject = signingblock.PkixNameToString(&cert.Subject)
	ci.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	ci.SerialNumber = cert.SerialNumber
}

// Returns description of the cert, like this:
// Cert 90d0f1ac70d647edfdf905ff129379bfae469ad6, valid from 2015-08-05 08:01:53 +0000 UTC to 2045-07-28 08:01:53 +0000 UTC,
// Subject C=US, O=Android, CN=Android Debug, Issuer C=US, O=Android, CN=Android Debug
func (ci *CertInfo) String() string {
	return fmt.Sprintf("Cert %s, valid from %s to %s, Subject: %s, Issuer: %s",
		ci.Sha1, ci.ValidFrom.Format(time.RFC3339), ci.ValidTo.Format(time.RFC3339), ci.Subject, ci.Issuer)
}
