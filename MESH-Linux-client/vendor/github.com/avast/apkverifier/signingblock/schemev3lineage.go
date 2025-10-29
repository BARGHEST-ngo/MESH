package signingblock

import (
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/avast/apkverifier/apilevel"
	"io"
	"strings"
)

type V3LineageSigningCertificateNode struct {
	SigningCert        *x509.Certificate
	ParentSigAlgorithm SignatureAlgorithm
	SigAlgorithm       SignatureAlgorithm
	Signature          []byte
	Flags              LineageCertCaps
}

// frameworks/base/core/java/android/content/pm/PackageParser.java
// public @interface CertCapabilities
type LineageCertCaps int32

const (
	CapInstalledData LineageCertCaps = 1  // accept data from already installed pkg with this cert
	CapSharedUserId  LineageCertCaps = 2  // accept sharedUserId with pkg with this cert
	CapPermission    LineageCertCaps = 4  // grant SIGNATURE permissions to pkgs with this cert
	CapRollback      LineageCertCaps = 8  // allow pkg to update to one signed by this certificate
	CapAuth          LineageCertCaps = 16 // allow pkg to continue to have auth access gated by this cert
)

func (c LineageCertCaps) String() string {
	var values []string
	for i := uint(0); i < 31; i++ {
		mask := LineageCertCaps(1 << i)
		if (c & mask) == 0 {
			continue
		}
		switch mask {
		case CapInstalledData:
			values = append(values, "InstalledData")
		case CapSharedUserId:
			values = append(values, "SharedUserId")
		case CapPermission:
			values = append(values, "Permission")
		case CapRollback:
			values = append(values, "Rollback")
		case CapAuth:
			values = append(values, "Auth")
		default:
			values = append(values, fmt.Sprintf("0x%x", mask))
		}
	}

	if len(values) != 0 {
		return strings.Join(values, "|")
	}
	return "None"
}

func (n *V3LineageSigningCertificateNode) Equal(o *V3LineageSigningCertificateNode) bool {
	if n == o {
		return true
	}

	if !n.SigningCert.Equal(o.SigningCert) {
		return false
	}

	if n.ParentSigAlgorithm != o.ParentSigAlgorithm {
		return false
	}

	if n.SigAlgorithm != o.SigAlgorithm {
		return false
	}

	if !bytes.Equal(n.Signature, o.Signature) {
		return false
	}

	if n.Flags != o.Flags {
		return false
	}

	return true
}

func (n *V3LineageSigningCertificateNode) Dump(w io.Writer) error {
	sha1sum := sha1.Sum(n.SigningCert.Raw)
	if _, err := fmt.Fprintf(w, "Cert %s:\n", hex.EncodeToString(sha1sum[:])); err != nil {
		return err
	}

	_, err := fmt.Fprintf(w, "  ValidFrom: %v\n  ValidTo: %v\n", n.SigningCert.NotBefore, n.SigningCert.NotAfter)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "  Issuer: %s\n  Subject: %s\n", PkixNameToString(&n.SigningCert.Issuer), PkixNameToString(&n.SigningCert.Subject))
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "ParentSigAlgorithm: %s\n", n.ParentSigAlgorithm.String()); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "SigAlgorithm: %s\n", n.SigAlgorithm.String()); err != nil {
		return err
	}

	sha1sum = sha1.Sum(n.Signature)
	if _, err := fmt.Fprintf(w, "Signature: %s\n", hex.EncodeToString(sha1sum[:])); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "Flags: 0x%04x (%s)\n", uint32(n.Flags), n.Flags.String()); err != nil {
		return err
	}
	return nil
}

type V3LineageSigningCertificateNodeList []*V3LineageSigningCertificateNode

func (l V3LineageSigningCertificateNodeList) Equal(o V3LineageSigningCertificateNodeList) bool {
	if len(l) != len(o) {
		return false
	}

	for i := range l {
		if !l[i].Equal(o[i]) {
			return false
		}
	}
	return true
}

func lineageCalculateMinSdkVersion(nodes V3LineageSigningCertificateNodeList) int32 {
	minSdkVersion := apilevel.V9_0_Pie // lineage introduced in P
	for _, n := range nodes {
		if n.SigAlgorithm.isSupported() {
			if nmin := n.SigAlgorithm.getMinSdkVersion(); nmin > minSdkVersion {
				minSdkVersion = nmin
			}
		}
	}
	return minSdkVersion
}

type V3SigningLineage struct {
	MinSdkVersion int32
	Nodes         V3LineageSigningCertificateNodeList
}

func (l *V3SigningLineage) getSubLineage(cert *x509.Certificate) (*V3SigningLineage, error) {
	for i, n := range l.Nodes {
		if n.SigningCert.Equal(cert) {
			return &V3SigningLineage{
				MinSdkVersion: l.MinSdkVersion,
				Nodes:         l.Nodes[0 : i+1],
			}, nil
		}
	}
	return nil, fmt.Errorf("certificate not found in signing lineage")
}
