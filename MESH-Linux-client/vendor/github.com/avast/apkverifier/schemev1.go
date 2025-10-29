package apkverifier

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/avast/apkverifier/internal/asn1andr"

	"github.com/avast/apkverifier/apilevel"

	"crypto/ecdsa"

	"github.com/avast/apkparser"
	"github.com/avast/apkverifier/fullsailor/pkcs7"
)

// These two arrays are synchronized
var (
	digestAlgorithms = [...]string{
		"sha-512",
		"sha-384",
		"sha-256",
		"sha1",
	}
	digestHashers = map[string]func() hash.Hash{
		"sha-512": sha512.New,
		"sha-384": sha512.New384,
		"sha-256": sha256.New,
		"sha1":    sha1.New,
	}
)

const (
	SHA224WithRSA x509.SignatureAlgorithm = iota + 65535
	DSAWithSHA224
	ECDSAWithSHA224
)

var (
	oidAttributeContentType   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
	oidAttributeMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	oidDigestAlgorithmSHA1    = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
	oidDigestAlgorithmSHA256  = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
)

var oidToAlgo = map[string]x509.SignatureAlgorithm{
	" 1.2.840.113549.1.1.2":                   x509.MD2WithRSA,
	"1.3.14.3.2.24 1.2.840.113549.1.1.2":      x509.MD2WithRSA,
	"1.2.840.113549.2.2 1.2.840.113549.1.1.1": x509.MD2WithRSA,

	" 1.2.840.113549.1.1.4":                   x509.MD5WithRSA,
	"1.3.14.3.2.25 1.2.840.113549.1.1.4":      x509.MD5WithRSA,
	"1.3.14.3.2.3 1.2.840.113549.1.1.4":       x509.MD5WithRSA,
	"1.2.840.113549.2.5 1.2.840.113549.1.1.1": x509.MD5WithRSA,
	"1.2.840.113549.2.5 1.2.840.113549.1.1.4": x509.MD5WithRSA,

	" 1.2.840.113549.1.1.5":              x509.SHA1WithRSA,
	"1.3.14.3.2.26 1.2.840.113549.1.1.1": x509.SHA1WithRSA,
	"1.3.14.3.2.29 1.2.840.113549.1.1.1": x509.SHA1WithRSA,
	"1.3.14.3.2.26 1.2.840.113549.1.1.5": x509.SHA1WithRSA,

	"2.16.840.1.101.3.4.2.4 1.2.840.113549.1.1.1":  SHA224WithRSA,
	"2.16.840.1.101.3.4.2.4 1.2.840.113549.1.1.14": SHA224WithRSA,

	"2.16.840.1.101.3.4.2.1 1.2.840.113549.1.1.1":  x509.SHA256WithRSA,
	"2.16.840.1.101.3.4.2.1 1.2.840.113549.1.1.11": x509.SHA256WithRSA,

	"2.16.840.1.101.3.4.2.2 1.2.840.113549.1.1.1":  x509.SHA384WithRSA,
	"2.16.840.1.101.3.4.2.2 1.2.840.113549.1.1.12": x509.SHA384WithRSA,

	"2.16.840.1.101.3.4.2.3 1.2.840.113549.1.1.1":  x509.SHA512WithRSA,
	"2.16.840.1.101.3.4.2.3 1.2.840.113549.1.1.13": x509.SHA512WithRSA,

	"1.3.14.3.2.26 1.2.840.10040.4.1": x509.DSAWithSHA1,
	" 1.2.840.10040.4.3":              x509.DSAWithSHA1,
	"1.3.14.3.2.26 1.2.840.10040.4.3": x509.DSAWithSHA1,

	"2.16.840.1.101.3.4.2.4 1.2.840.10040.4.1":      DSAWithSHA224,
	"2.16.840.1.101.3.4.2.4 2.16.840.1.101.3.4.3.1": DSAWithSHA224,

	"2.16.840.1.101.3.4.2.1 1.2.840.10040.4.1":      x509.DSAWithSHA256,
	"2.16.840.1.101.3.4.2.1 2.16.840.1.101.3.4.3.2": x509.DSAWithSHA256,

	"1.3.14.3.2.26 1.2.840.10045.2.1": x509.ECDSAWithSHA1,
	"1.3.14.3.2.26 1.2.840.10045.4.1": x509.ECDSAWithSHA1,

	"2.16.840.1.101.3.4.2.4 1.2.840.10045.2.1":   ECDSAWithSHA224,
	"2.16.840.1.101.3.4.2.4 1.2.840.10045.4.3.1": ECDSAWithSHA224,

	"2.16.840.1.101.3.4.2.1 1.2.840.10045.2.1":   x509.ECDSAWithSHA256,
	"2.16.840.1.101.3.4.2.1 1.2.840.10045.4.3.2": x509.ECDSAWithSHA256,

	"2.16.840.1.101.3.4.2.2 1.2.840.10045.2.1":   x509.ECDSAWithSHA384,
	"2.16.840.1.101.3.4.2.2 1.2.840.10045.4.3.3": x509.ECDSAWithSHA384,

	"2.16.840.1.101.3.4.2.3 1.2.840.10045.2.1":   x509.ECDSAWithSHA512,
	"2.16.840.1.101.3.4.2.3 1.2.840.10045.4.3.4": x509.ECDSAWithSHA512,
}

var errNoKnownHashes = errors.New("No known hashes")

type schemeV1Signature struct {
	sigBlockFilename  string
	manifestFilename  string
	cert              *pkcs7.PKCS7
	signatureManifest *manifest
	chain             []*x509.Certificate
}

type schemeV1 struct {
	sigs     map[string]*schemeV1Signature
	manifest *manifest
	hashers  map[string]hash.Hash
	chain    [][]*x509.Certificate
}

// The order of signature block files is not deterministic on Android because the file list goes through two hash maps.
// First one is in the libziparchive and the second is java HashMap in StrictJarVerifier. This means that the verification
// can sometimes fail and sometimes succeed on files that have two types of signature block file (e.g. TEST.RSA and TEST.DSA),
// depending on which file Android parses first.
//
// We keep the behavior deterministic in our implementation, based on the file order in the ZIP.
// This unfortunately produces different result than Android in some corner cases.
func verifySchemeV1(apk *apkparser.ZipReader, hasValidSigningBlock bool, minSdkVersion, maxSdkVersion int32) ([][]*x509.Certificate, error) {
	scheme, err := newSchemeV1(apk)
	if err != nil {
		return nil, err
	}

	err = scheme.verify(apk, hasValidSigningBlock, minSdkVersion, maxSdkVersion)
	return scheme.chain, err
}

func extractCertsSchemeV1(apk *apkparser.ZipReader, minSdkVersion, maxSdkVersion int32) ([][]*x509.Certificate, error) {
	scheme, err := newSchemeV1(apk)
	if err != nil {
		return nil, err
	}

	var signatureErrors []error
	for _, sig := range scheme.sigs {
		sig.chain, err = scheme.verifySignature(sig, minSdkVersion, maxSdkVersion)
		if sig.chain != nil {
			scheme.chain = append(scheme.chain, sig.chain)
		}
		if err != nil {
			signatureErrors = append(signatureErrors, fmt.Errorf("%s: %s", sig.sigBlockFilename, err))
		}
	}

	if len(signatureErrors) != 0 {
		return scheme.chain, fmt.Errorf("One or more of the signatures are invalid: %v", signatureErrors)
	}
	return scheme.chain, nil
}

func newSchemeV1(apk *apkparser.ZipReader) (*schemeV1, error) {
	scheme := schemeV1{
		sigs:    make(map[string]*schemeV1Signature),
		hashers: make(map[string]hash.Hash),
	}

	const prefix = "META-INF/"
	var signatureBlocks []*apkparser.ZipReaderFile
	signatureFiles := map[string]*apkparser.ZipReaderFile{}
	for _, f := range apk.FilesOrdered {
		if !strings.HasPrefix(f.Name, prefix) {
			continue
		}

		switch {
		case f.Name == "META-INF/MANIFEST.MF":
			if err := scheme.addManifest(f); err != nil {
				return nil, fmt.Errorf("failed to parse main manifest: %s", err.Error())
			}
		case strings.HasSuffix(f.Name, ".RSA") || strings.HasSuffix(f.Name, ".DSA") || strings.HasSuffix(f.Name, ".EC"):
			signatureBlocks = append(signatureBlocks, f)
		case strings.HasSuffix(f.Name, ".SF"):
			if _, prs := signatureFiles[f.Name]; !prs {
				signatureFiles[f.Name] = f
			}
		}
	}

	var errors []error
	for _, blockFile := range signatureBlocks {
		name := blockFile.Name
		dot := strings.LastIndexByte(name, '.')
		sfname := name[:dot] + ".SF"

		sf, prs := signatureFiles[sfname]
		if !prs {
			continue
		}

		if err := scheme.addSignatureBlock(name, blockFile); err != nil {
			// Behavior changed in Android 7.0 - badly formed signature blocks are no longer ignored
			// errors = append(errors, fmt.Errorf("%s: %s", name, err.Error()))
			// continue
			return nil, fmt.Errorf("%s: %s", name, err.Error())
		}

		if err := scheme.addSignatureFile(sfname, sf); err != nil {
			errors = append(errors, fmt.Errorf("%s: %s", name, err.Error()))
			continue
		}

		// The same signatureFile can't be used by another signature block
		delete(signatureFiles, sfname)
	}

	if err := scheme.prepForVerification(); err != nil {
		if len(errors) == 0 {
			return nil, fmt.Errorf("Can't verify: %s", err.Error())
		} else {
			return nil, fmt.Errorf("Can't verify: %s %v", err.Error(), errors)
		}
	}
	return &scheme, nil
}

func (p *schemeV1) addManifest(f *apkparser.ZipReaderFile) (err error) {
	if p.manifest != nil {
		return fmt.Errorf("Manifest already parsed!")
	}

	p.manifest, err = parseManifest(f, true)
	return
}

func (p *schemeV1) addSignatureFile(pathUpper string, f *apkparser.ZipReaderFile) (err error) {
	prefix := p.signaturePrefix(pathUpper)
	s := p.sigs[prefix]
	if s == nil {
		s = &schemeV1Signature{}
		p.sigs[prefix] = s
	}

	s.manifestFilename = pathUpper
	s.signatureManifest, err = parseManifest(f, false)
	return
}

func (p *schemeV1) addSignatureBlock(pathUpper string, f *apkparser.ZipReaderFile) error {
	if err := f.Open(); err != nil {
		return err
	}
	defer f.Close()

	var err error
	var raw []byte
	var sig *pkcs7.PKCS7
	for f.Next() {
		raw, err = ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		sig, err = pkcs7.Parse(raw)
		if err != nil {
			continue
		}

		prefix := p.signaturePrefix(pathUpper)
		s := p.sigs[prefix]
		if s == nil {
			s = &schemeV1Signature{}
			p.sigs[prefix] = s
		}
		s.sigBlockFilename = f.Name
		s.cert = sig

		return nil
	}

	return fmt.Errorf("failed to open: %v", err)
}

func (p *schemeV1) signaturePrefix(pathUpper string) string {
	fn := filepath.Base(pathUpper)
	idx := strings.LastIndexByte(fn, '.')
	return fn[:idx]
}

func (p *schemeV1) prepForVerification() error {
	if p.manifest == nil {
		return errors.New("No valid MANIFEST.SF")
	}

	for prefix, sig := range p.sigs {
		if sig.cert == nil || sig.signatureManifest == nil {
			delete(p.sigs, prefix)
		}
	}

	if len(p.sigs) == 0 {
		return errors.New("No signatures.")
	}

	return nil
}

func (p *schemeV1) verify(apk *apkparser.ZipReader, hasValidSigningBlock bool, minSdkVersion, maxSdkVersion int32) error {
	var err error
	validSignatures := map[string]*schemeV1Signature{}
	var signatureErrors []error
	for sigName, sig := range p.sigs {
		sig.chain, err = p.verifySignature(sig, minSdkVersion, maxSdkVersion)
		if sig.chain != nil {
			p.chain = append(p.chain, sig.chain)
		}
		if err != nil {
			signatureErrors = append(signatureErrors, fmt.Errorf("%s: %s", sig.sigBlockFilename, err))
			continue
		}

		sm := sig.signatureManifest
		if idList, prs := sm.main[attrAndroidApkSigned]; !hasValidSigningBlock && prs && apilevel.SupportsSigV2(maxSdkVersion) {
			tokens := strings.Split(idList, ",")
			for _, tok := range tokens {
				tok = strings.TrimSpace(tok)
				if tok == "" {
					continue
				}

				id, err := strconv.ParseInt(tok, 10, 32)
				if err != nil {
					continue
				}

				if id == 2 {
					return fmt.Errorf("This apk has '%s: %s', cannot be verified using v1 scheme, downgrade attack?",
						attrAndroidApkSigned, idList)
				}
			}
		}

		if _, prs := sm.main[attrSignatureVersion]; !prs {
			// Android just ignores it
			//return fmt.Errorf("the manifest of %s does not have %s attribute", sig.manifestFilename, attrSignatureVersion)
			continue
		}

		createdBySigntool := strings.Contains(sm.main[attrCreatedBy], "signtool")

		if sm.mainAttributtesEnd > 0 && !createdBySigntool {
			err = p.verifyManifestEntry(sm.main, attrDigestMainAttrSuffix, minSdkVersion, maxSdkVersion, func(hash []byte, hasher hash.Hash) error {
				hasher.Write(p.manifest.rawData[:p.manifest.mainAttributtesEnd])
				if !bytes.Equal(hash, hasher.Sum(nil)) {
					return fmt.Errorf("Invalid manifest %s main attributes hash!", sig.manifestFilename)
				}
				return nil
			})

			if err != nil && err != errNoKnownHashes {
				return fmt.Errorf("failed to verify manifest %s main attributes: %s", sig.manifestFilename, err.Error())
			}
		}

		suffix := attrDigestSigntoolSuffix
		if createdBySigntool {
			suffix = attrDigestSuffix
		}

		err = p.verifyManifestEntry(sm.main, suffix, minSdkVersion, maxSdkVersion, func(hash []byte, hasher hash.Hash) error {
			if hasher.Write(p.manifest.rawData); !bytes.Equal(hash, hasher.Sum(nil)) {
				return errors.New("Invalid whole manifest hash!")
			}
			return nil
		})

		// file entries only checked if the whole-manifest fails/is not present
		if err != nil {
			for name, attrs := range sm.entries {
				err = p.verifyManifestEntry(attrs, attrDigestSuffix, minSdkVersion, maxSdkVersion, func(hash []byte, hasher hash.Hash) error {
					data, prs := p.manifest.chunks[name]
					if !prs {
						return fmt.Errorf("Signature entry %s not in manifest.mf file.", name)
					}

					if createdBySigntool && bytes.HasSuffix(data, []byte{'\n', '\n'}) {
						hasher.Write(data[:len(data)-1])
					} else {
						hasher.Write(data)
					}

					if !bytes.Equal(hash, hasher.Sum(nil)) {
						return fmt.Errorf("Invalid hash of manifest entry for %s", name)
					}
					return nil
				})

				if err != nil {
					break
				}
			}
		}

		if err == nil {
			validSignatures[sigName] = sig
		}
	}

	p.sigs = validSignatures

	if len(validSignatures) == 0 {
		return fmt.Errorf("No valid cert chains found, last error: %v", err)
	}

	if len(signatureErrors) != 0 {
		return fmt.Errorf("One or more of the signatures are invalid: %v", signatureErrors)
	}

	if len(p.chain) > maxApkSigners {
		return fmt.Errorf("APK Signature Scheme v1 only supports a maximum of %d signers, found %d", maxApkSigners, len(p.chain))
	}

	return p.verifyMainManifest(apk, minSdkVersion, maxSdkVersion)
}

func (p *schemeV1) verifyMainManifest(apk *apkparser.ZipReader, minSdkVersion, maxSdkVersion int32) error {
	for path := range p.manifest.entries {
		if _, prs := apk.File[path]; !prs {
			return fmt.Errorf("Manifest entry '%s' does not exists.", path)
		}
	}

	required := make([]string, 0, len(apk.File))
	hasAndroidManifest := false
	hasBundleConfig := false
	for path, zf := range apk.File {
		if zf.IsDir {
			continue
		}

		if !zf.IsDir && path != "AndroidManifest.xml" && !strings.HasPrefix(path, "META-INF/") {
			required = append(required, path)
		}

		if path == "AndroidManifest.xml" {
			hasAndroidManifest = true
		} else if path == "BundleConfig.pb" {
			hasBundleConfig = true
		}
	}

	if !hasAndroidManifest && !hasBundleConfig {
		required = append(required, "AndroidManifest.xml")
	}

	chainsSet := false
	for _, path := range required {
		attrs, prs := p.manifest.entries[path]
		if !prs {
			return fmt.Errorf("No manifest entry for required file '%s'", path)
		}

		err := p.verifyManifestEntry(attrs, attrDigestSuffix, minSdkVersion, maxSdkVersion, func(hash []byte, hasher hash.Hash) error {
			return p.verifyFileHash(apk.File[path], hash, hasher)
		})
		if err != nil {
			return err
		}

		var certChains [][]*x509.Certificate
		for _, sig := range p.sigs {
			if _, prs := sig.signatureManifest.entries[path]; prs {
				certChains = append(certChains, sig.chain)
			}
		}

		if len(certChains) == 0 {
			return fmt.Errorf("File '%s' is not in any signature manifests", path)
		}

		if !chainsSet {
			p.chain = certChains
			chainsSet = true
		} else if !p.certChainsMatch(p.chain, certChains) {
			return fmt.Errorf("Mismatched certificates at entry '%s'", path)
		}
	}
	return nil
}

func (p *schemeV1) certChainsMatch(a, b [][]*x509.Certificate) bool {
	if len(a) != len(b) {
		return false
	}

	for _, ca := range a {
		found := false
		for _, cb := range b {
			if p.chainEqual(ca, cb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	for _, cb := range b {
		found := false
		for _, ca := range a {
			if p.chainEqual(ca, cb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (p *schemeV1) chainEqual(a, b []*x509.Certificate) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func (p *schemeV1) getDigestsToVerify(entry map[string]string, suffix string, minSdkVersion, maxSdkVersion int32) []string {
	var res []string
	if minSdkVersion < 18 {
		algs := strings.ToLower(entry["digest-algorithms"])
		if algs == "" {
			algs = "sha sha1"
		}

		tokens := strings.Split(algs, " ")
		for _, algo := range tokens {
			if minSdkVersion >= 9 || (algo != "sha-384" && algo != "sha-512") {
				if _, prs := entry[algo+suffix]; prs {
					res = append(res, algo)
				}
			}
		}

		// apksig fails the verification in this case, because pre-18 Android will, too.
		// We don't want to, newer devices are more relevant to us.
		/*if len(res) == 0 {
			return res
		}*/
	}

	if maxSdkVersion >= 18 {
		for _, algo := range digestAlgorithms {
			if _, prs := entry[algo+suffix]; prs {
				res = append(res, algo)
				break
			}
		}
	}

	return res
}

func (p *schemeV1) verifyManifestEntry(entry map[string]string, digestSuffix string, minSdkVersion, maxSdkVersion int32, verify func(hash []byte, hasher hash.Hash) error) error {
	toVerify := p.getDigestsToVerify(entry, digestSuffix, minSdkVersion, maxSdkVersion)
	if len(toVerify) == 0 {
		return errNoKnownHashes
	}

	for _, algo := range toVerify {
		hash64 := entry[algo+digestSuffix]

		hash, err := base64.StdEncoding.DecodeString(hash64)
		if err != nil {
			return fmt.Errorf("Can't decode hash: %s", err.Error())
		}

		if p.hashers[algo] == nil {
			factory, prs := digestHashers[algo]
			if !prs {
				return errNoKnownHashes
			}
			p.hashers[algo] = factory()
		}
		p.hashers[algo].Reset()

		if err := verify(hash, p.hashers[algo]); err != nil {
			return err
		}
	}
	return nil
}

func (p *schemeV1) verifyFileHash(f *apkparser.ZipReaderFile, hash []byte, hasher hash.Hash) error {
	if err := f.Open(); err != nil {
		return fmt.Errorf("Can't generate hashes for '%s': %s", f.Name, err.Error())
	}
	defer f.Close()

	for f.Next() {
		hasher.Reset()
		if _, err := io.Copy(hasher, f); err == nil {
			if bytes.Equal(hasher.Sum(nil), hash) {
				return nil
			}
		}
	}

	return fmt.Errorf("No matching hash for '%s'!", f.Name)
}

func (p *schemeV1) getHashForOID(oid asn1.ObjectIdentifier) (crypto.Hash, error) {
	switch {
	case oid.Equal(oidDigestAlgorithmSHA1):
		return crypto.SHA1, nil
	case oid.Equal(oidDigestAlgorithmSHA256):
		return crypto.SHA256, nil
	}
	return crypto.Hash(0), fmt.Errorf("unsupported hash algorithm oid %s", oid.String())
}

func (p *schemeV1) verifySignature(sig *schemeV1Signature, minSdkVersion, maxSdkVersion int32) ([]*x509.Certificate, error) {
	if len(sig.cert.Signers) == 0 {
		return nil, errors.New("Empty signers slice!")
	}

	signers := sig.cert.GetSignerInfos()
	// Prior to Android N, Android attempts to verify only the first SignerInfo. From N
	// onwards, Android attempts to verify all SignerInfos and then picks the first verified
	// SignerInfo.
	if minSdkVersion < apilevel.V7_0_Nougat {
		signers = signers[:1]
	}

	var firstVerifiedSignerCert *x509.Certificate
	var chain []*x509.Certificate
	var lastError error
	for i := range signers {
		info := &signers[i]

		var issuerSeq pkix.RDNSequence
		if _, err := asn1andr.Unmarshal(info.IssuerAndSerialNumber.IssuerName.FullBytes, &issuerSeq); err != nil {
			return nil, err
		}
		var issuer pkix.Name
		issuer.FillFromRDNSequence(&issuerSeq)
		issuerCanonical := p.pkixCanonical(&issuer)

		snum := info.IssuerAndSerialNumber.SerialNumber
		signerCertIndex := -1
		for i, crt := range sig.cert.Certificates {
			if snum.Cmp(crt.SerialNumber) == 0 && issuerCanonical == p.pkixCanonical(&crt.Issuer) {
				signerCertIndex = i
				break
			}
		}

		if signerCertIndex == -1 {
			return nil, errors.New("No issuer certificate found")
		}

		signerCert := sig.cert.Certificates[signerCertIndex]
		chain = []*x509.Certificate{signerCert}

		if len(signerCert.UnhandledCriticalExtensions) != 0 {
			return chain, errors.New("Certificate has unhandled critical extensions.")
		}

		da := info.DigestAlgorithm.Algorithm.String()
		dea := info.DigestEncryptionAlgorithm.Algorithm.String()
		algo, prs := oidToAlgo[fmt.Sprintf("%s %s", da, dea)]
		if !prs {
			panic(fmt.Sprintf("Unknown digest algorithm: '%s %s'", da, dea))
		}

		var signedData []byte
		// Signed attributes present -- verify signature against the ASN.1 DER encoded form
		// of signed attributes. This verifies integrity of the signature file because
		// signed attributes must contain the digest of the signature file.
		if len(info.AuthenticatedAttributes) != 0 {
			// Prior to Android KitKat, APKs with signed attributes are unsafe:
			// * The APK's contents are not protected by the JAR signature because the
			//   digest in signed attributes is not verified. This means an attacker can
			//   arbitrarily modify the APK without invalidating its signature.
			// * Luckily, the signature over signed attributes was verified incorrectly
			//   (over the verbatim IMPLICIT [0] form rather than over re-encoded
			//   UNIVERSAL SET form) which means that JAR signatures which would verify on
			//   pre-KitKat Android and yet do not protect the APK from modification could
			//   be generated only by broken tools or on purpose by the entity signing the
			//   APK.
			//
			// We thus reject such unsafe APKs, even if they verify on platforms before
			// KitKat.

			// We do not care, it installs on real devices
			// 0f2b96555e09ef5dd0e18c360c0b1b35666f68d0ed312586a259a3d1ee39b68a
			/*if minSdkVersion < apilevel.V4_4_KitKat {
				return chain, errors.New("APKs with Signed Attributes broken on platforms API LEVEL < 19")
			}*/

			if maxSdkVersion >= apilevel.V7_0_Nougat {
				var typeVal asn1.ObjectIdentifier
				if err := info.UnmarshalSignedAttribute(oidAttributeContentType, &typeVal); err != nil {
					return chain, fmt.Errorf("failed to parse signed content type: %s", err.Error())
				}

				// Did not verify: Content type signed attribute does not match
				// SignedData.encapContentInfo.eContentType. This fails verification of
				// this SignerInfo but should not prevent verification of other
				// SignerInfos. Hence, no exception is thrown.
				if !typeVal.Equal(sig.cert.ContentType) {
					lastError = fmt.Errorf("PKCS7 content type does not match, %s != %s", typeVal.String(), sig.cert.ContentType.String())
					continue
				}
			}

			var digest []byte
			err := info.UnmarshalSignedAttribute(oidAttributeMessageDigest, &digest)
			if err != nil {
				return chain, fmt.Errorf("failed to parse signed message digest: %s", err.Error())
			}

			hash, err := p.getHashForOID(info.DigestAlgorithm.Algorithm)
			if err != nil {
				return chain, err
			}

			h := hash.New()
			h.Write(sig.signatureManifest.rawData)
			computed := h.Sum(nil)
			if !bytes.Equal(digest, computed) {
				// Skip verification: signature file digest in signed attributes does not
				// match the signature file. This fails verification of
				// this SignerInfo but should not prevent verification of other
				// SignerInfos. Hence, no exception is thrown.
				lastError = errors.New("signedAttributes hash mismatch")
				continue
			}

			signedData, err = info.MarshalAuthenticatedAttributes()
			if err != nil {
				return chain, err
			}
		} else {
			signedData = sig.signatureManifest.rawData
		}

		err := p.checkSignature(signerCert, algo, signedData, info.EncryptedDigest)
		if err != nil {
			lastError = err
			continue
		}

		if firstVerifiedSignerCert == nil {
			firstVerifiedSignerCert = signerCert
		}
	}

	if firstVerifiedSignerCert == nil {
		return nil, fmt.Errorf("no valid signers: %v", lastError)
	}

	// load cert chain if not self-signed
	chain = []*x509.Certificate{firstVerifiedSignerCert}
	if p.pkixCanonical(&firstVerifiedSignerCert.Issuer) == p.pkixCanonical(&firstVerifiedSignerCert.Subject) {
		return chain, nil
	}

	issuerCanonical := p.pkixCanonical(&firstVerifiedSignerCert.Issuer)
	for {
		var issuerCert *x509.Certificate
		for _, crt := range sig.cert.Certificates {
			if issuerCanonical == p.pkixCanonical(&crt.Subject) {
				issuerCert = crt
				break
			}
		}

		if issuerCert == nil {
			break
		}

		chain = append(chain, issuerCert)
		if len(chain) > len(sig.cert.Certificates) {
			break
		}

		issuerCanonical = p.pkixCanonical(&issuerCert.Issuer)
		if issuerCanonical == p.pkixCanonical(&issuerCert.Subject) {
			break
		}
	}

	return chain, nil
}

type dsaSignature struct {
	R, S *big.Int
}
type ecdsaSignature dsaSignature

func (p *schemeV1) checkSignature(cert *x509.Certificate, algo x509.SignatureAlgorithm, signed, signature []byte) error {
	// Go1.15 rejects signatures without padding, add one.
	if rsaPub, ok := cert.PublicKey.(*rsa.PublicKey); ok {
		pubSize := (rsaPub.N.BitLen() + 7) / 8 // rsaPub.Size(), but .Size() is only since go1.11
		if len(signature) < pubSize {
			signature = append(make([]byte, pubSize-len(signature)), signature...)
		}
	}

	switch algo {
	case x509.MD5WithRSA:
		digest := md5.Sum(signed)
		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("Unexpected public key type (%T)!", cert.PublicKey)
		}
		return rsa.VerifyPKCS1v15(pub, crypto.MD5, digest[:], signature)
	case SHA224WithRSA:
		digest := sha256.Sum224(signed)
		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("Unexpected public key type (%T)!", cert.PublicKey)
		}
		return rsa.VerifyPKCS1v15(pub, crypto.SHA224, digest[:], signature)
	case x509.DSAWithSHA1, DSAWithSHA224, x509.DSAWithSHA256:
		var hasher hash.Hash
		switch algo {
		case x509.DSAWithSHA1:
			hasher = sha1.New()
		case x509.DSAWithSHA256:
			hasher = sha256.New()
		case DSAWithSHA224:
			hasher = sha256.New224()
		}

		hasher.Write(signed)
		hash := hasher.Sum(nil)

		pub := cert.PublicKey.(*dsa.PublicKey)
		reqLen := pub.Q.BitLen() / 8
		if reqLen > len(hash) {
			reqLen = len(hash)
			// Android doesn't care?
			//return fmt.Errorf("Digest algorithm is too short for given DSA parameters.")
		}
		digest := hash[:reqLen]

		dsaSig := new(dsaSignature)
		if rest, err := asn1andr.Unmarshal(signature, dsaSig); err != nil {
			return err
		} else if len(rest) != 0 {
			return errors.New("x509: trailing data after DSA signature")
		}
		if dsaSig.R.Sign() <= 0 || dsaSig.S.Sign() <= 0 {
			return errors.New("x509: DSA signature contained zero or negative values")
		}
		if !dsa.Verify(pub, digest, dsaSig.R, dsaSig.S) {
			return errors.New("x509: DSA verification failure")
		}
		return nil
	case ECDSAWithSHA224:
		digest := sha256.Sum224(signed)
		pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return fmt.Errorf("Unexpected public key type (%T)!", cert.PublicKey)
		}

		ecdsaSig := new(ecdsaSignature)
		if rest, err := asn1andr.Unmarshal(signature, ecdsaSig); err != nil {
			return err
		} else if len(rest) != 0 {
			return errors.New("x509: trailing data after ECDSA signature")
		}
		if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
			return errors.New("x509: ECDSA signature contained zero or negative values")
		}
		if !ecdsa.Verify(pub, digest[:], ecdsaSig.R, ecdsaSig.S) {
			return errors.New("x509: ECDSA verification failure")
		}
		return nil
	default:
		return cert.CheckSignature(algo, signed, signature)
	}
}

type byX501Canonical []pkix.AttributeTypeAndValue

func (a byX501Canonical) Len() int      { return len(a) }
func (a byX501Canonical) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byX501Canonical) Less(i, j int) bool {
	ioid1 := a[i].Type
	ioid2 := a[j].Type
	min := len(ioid1)
	if len(ioid2) < min {
		min = len(ioid2)
	}
	for x := 0; x < min; x++ {
		if ioid1[x] < ioid2[x] {
			return true
		} else if ioid1[x] > ioid2[x] {
			return false
		}

		if (x+1) == len(ioid1) && (x+1) < len(ioid2) {
			return true
		} else if (x+1) < len(ioid1) && (x+1) == len(ioid2) {
			return false
		}
	}
	return false
}

func (p *schemeV1) pkixCanonical(n *pkix.Name) string {
	return p.pkixCanonicalSeq(n.ToRDNSequence())
}

func (p *schemeV1) pkixCanonicalSeq(n pkix.RDNSequence) string {
	var res bytes.Buffer
	for i := len(n) - 1; i >= 0; i-- {
		atavList := n[i]
		sort.Sort(byX501Canonical(atavList))

		for _, atav := range atavList {
			fmt.Fprintf(&res, "%s=", atav.Type.String())
			switch val := atav.Value.(type) {
			case string:
				length := len(val)
				if length == 0 {
					break
				}

				index := 0
				bufStart := res.Len()
				if val[0] == '#' {
					res.WriteString("\\#")
					index++
				}

				for ; index < length; index++ {
					switch val[index] {
					case ' ':
						bufLen := res.Len() - bufStart
						if bufLen == 0 || res.Bytes()[res.Len()-1] == ' ' {
							break
						}
						res.WriteByte(' ')
					case '"', '\\', ',', '+', '<', '>', ';':
						res.WriteByte('\\')
					default:
						res.WriteByte(val[index])
					}
				}

				x := res.Len() - 1
				for x >= bufStart && res.Bytes()[x] == ' ' {
					x--
				}
				res.Truncate(x + 1)
			default:
				fmt.Fprintf(&res, "%v", atav.Value)
			}
			res.WriteByte('+')
		}

		// remove last +
		if len(atavList) != 0 {
			res.Truncate(res.Len() - 1)
		}

		if i != 0 {
			res.WriteByte(',')
		}
	}
	return strings.ToLower(res.String())
}
