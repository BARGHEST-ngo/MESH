// Package apkverifier does APK signature verification.
// It should support all algorithms and schemes supported Android, including scheme v2 verification
// and checks for downgrade attack to v1.
package apkverifier

import (
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/avast/apkparser"
	"github.com/avast/apkverifier/apilevel"
	"github.com/avast/apkverifier/signingblock"
)

// Result Contains result of Apk verification
type Result struct {
	SigningSchemeId    int
	SignerCerts        [][]*x509.Certificate
	SigningBlockResult *signingblock.VerificationResult
}

// ErrMixedDexApkFile Returned from the Verify method if the file starts with the DEX magic value,
// but otherwise looks like a properly signed APK.
//
// This detects 'Janus' Android vulnerability where a DEX is prepended to a valid,
// signed APK file. The signature verification passes because with v1 scheme,
// only the APK portion of the file is checked, but Android then loads the prepended,
// unsigned DEX file instead of the one from APK.
// https://www.guardsquare.com/en/blog/new-android-vulnerability-allows-attackers-modify-apps-without-affecting-their-signatures
//
// If this error is returned, the signature is otherwise valid (the err would be nil
// had it not have the DEX file prepended).
var ErrMixedDexApkFile = errors.New("This file is both DEX and ZIP archive! Exploit?")

const (
	dexHeaderMagic uint32 = 0xa786564 // "dex\n", littleendinan

	maxApkSigners = 10
)

// Verify Calls VerifyWithSdkVersion with sdk versions <apilevel.V_AnyMin; apilevel.V_AnyMax>
func Verify(path string, optionalZip *apkparser.ZipReader) (res Result, err error) {
	return VerifyWithSdkVersion(path, optionalZip, apilevel.V_AnyMin, apilevel.V_AnyMax)
}

// VerifyReader Calls VerifyWithSdkVersionReader with sdk versions <apilevel.V_AnyMin; apilevel.V_AnyMax>
func VerifyReader(r io.ReadSeeker, optionalZip *apkparser.ZipReader) (res Result, err error) {
	return VerifyWithSdkVersionReader(r, optionalZip, apilevel.V_AnyMin, apilevel.V_AnyMax)
}

// VerifyWithSdkVersion see VerifyWithSdkVersionReader
func VerifyWithSdkVersion(path string, optionalZip *apkparser.ZipReader, minSdkVersion, maxSdkVersion int32) (res Result, err error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()
	return VerifyWithSdkVersionReader(f, optionalZip, minSdkVersion, maxSdkVersion)
}

// VerifyWithSdkVersionReader Verify the application signature. If err is nil, the signature is correct,
// otherwise it is not and res may or may not contain extracted certificates,
// depending on how the signature verification failed.
// Path is required, pass optionalZip if you have the ZipReader already opened and want to reuse it.
// This method will not close it.
// minSdkVersion and maxSdkVersion means the apk has to successfuly verify on real devices with sdk version
// inside the <minSdkVersion;maxSdkVersion> interval.
// minSdkVersion == apilevel.V_AnyMin means it will obtain the minSdkVersion from AndroidManifest.
func VerifyWithSdkVersionReader(r io.ReadSeeker, optionalZip *apkparser.ZipReader, minSdkVersion, maxSdkVersion int32) (res Result, err error) {
	if optionalZip == nil {
		optionalZip, err = apkparser.OpenZipReader(r)
		if err != nil {
			return Result{}, err
		}
		defer optionalZip.Close()
	}

	var sandboxVersion, targetSdkVersion int32
	var manifestError error
	if minSdkVersion == apilevel.V_AnyMin || apilevel.RequiresSandboxV2(maxSdkVersion) {
		var manifestMinSdkVersion int32
		manifestMinSdkVersion, targetSdkVersion, sandboxVersion, err = getManifestInfo(optionalZip)
		if err != nil {
			manifestError = err
		} else {
			if minSdkVersion == apilevel.V_AnyMin {
				minSdkVersion = manifestMinSdkVersion
			}
		}
	}

	if minSdkVersion > maxSdkVersion {
		err = fmt.Errorf("invalid sdk version range <%d;%d>", minSdkVersion, maxSdkVersion)
		return
	}

	var fileMagic uint32
	var signingBlockError error
	res.SigningBlockResult, fileMagic, signingBlockError = signingblock.VerifySigningBlockReaderWithZip(r, minSdkVersion, maxSdkVersion, optionalZip)

	if res.SigningBlockResult != nil && !signingblock.IsSigningBlockNotFoundError(signingBlockError) {
		res.SignerCerts = res.SigningBlockResult.Certs
		res.SigningSchemeId = res.SigningBlockResult.SchemeId
	} else {
		res.SigningSchemeId = 1
	}

	if signingBlockError != nil && !signingblock.IsSigningBlockNotFoundError(signingBlockError) {
		return res, signingBlockError
	} else if apilevel.SupportsSigV2(minSdkVersion) && signingBlockError == nil {
		return res, nil
	}

	// Android O and newer requires that APKs targeting security sandbox version 2 and higher
	// are signed using APK Signature Scheme v2 or newer.
	var sandboxError error
	if apilevel.RequiresSandboxV2(maxSdkVersion) && sandboxVersion > 1 && (signingBlockError != nil || res.SigningSchemeId < 2) {
		sandboxError = fmt.Errorf("no valid signature for sandbox version %d", sandboxVersion)
	} else if targetSdkVersion >= apilevel.V11_0_Eleven && maxSdkVersion >= targetSdkVersion && res.SigningSchemeId < 2 {
		sandboxError = fmt.Errorf("target SDK version %d requires a minimum of signature scheme v2; the APK is not signed with this or a later signature scheme", targetSdkVersion)
	}

	var certsv1 [][]*x509.Certificate
	certsv1, err = verifySchemeV1(optionalZip, signingBlockError == nil, minSdkVersion, maxSdkVersion)
	if len(res.SignerCerts) == 0 {
		res.SignerCerts = certsv1
	}

	if sandboxError != nil {
		err = sandboxError
	} else if err == nil {
		if res.SigningSchemeId != 1 && manifestError != nil {
			err = manifestError
		} else if fileMagic == dexHeaderMagic {
			err = ErrMixedDexApkFile
		}
	}

	return
}

// ExtractCerts Extract certs without verifying the signature.
func ExtractCerts(path string, optionalZip *apkparser.ZipReader) ([][]*x509.Certificate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ExtractCertsReader(f, optionalZip)
}

func ExtractCertsReader(r io.ReadSeeker, optionalZip *apkparser.ZipReader) ([][]*x509.Certificate, error) {
	var err error
	if optionalZip == nil {
		optionalZip, err = apkparser.OpenZipReader(r)
		if err != nil {
			return nil, err
		}
		defer optionalZip.Close()
	}

	certs, signingBlockError := signingblock.ExtractCertsReader(r, apilevel.V_AnyMin, apilevel.V_AnyMax)
	if !signingblock.IsSigningBlockNotFoundError(signingBlockError) {
		return certs, signingBlockError
	}

	var certsv1 [][]*x509.Certificate
	certsv1, err = extractCertsSchemeV1(optionalZip, apilevel.V_AnyMin, apilevel.V_AnyMax)
	certs = append(certs, certsv1...)
	return certs, err
}

type sandboxVersionEncoder struct {
	minSdkVersion    int32
	targetSdkVersion int32
	sandboxVersion   int32
}

func (e *sandboxVersionEncoder) EncodeToken(t xml.Token) error {
	st, ok := t.(xml.StartElement)
	if !ok {
		return nil
	}

	switch st.Name.Local {
	case "manifest":
		val, err := e.getAttrIntValue(&st, "targetSandboxVersion")
		if err == nil {
			e.sandboxVersion = val
		} else if err != io.EOF {
			return err
		}
	case "uses-sdk":
		val, err := e.getAttrIntValue(&st, "minSdkVersion")
		if err == nil {
			e.minSdkVersion = val
		} else if err != io.EOF {
			return err
		}

		val, err = e.getAttrIntValue(&st, "targetSdkVersion")
		if err == nil {
			e.targetSdkVersion = val
		} else if err != io.EOF {
			return err
		}

		return apkparser.ErrEndParsing
	}
	return nil
}

func (e *sandboxVersionEncoder) Flush() error {
	return nil
}

func (e *sandboxVersionEncoder) getAttrIntValue(st *xml.StartElement, name string) (int32, error) {
	for _, attr := range st.Attr {
		if attr.Name.Local == name {
			val, err := strconv.ParseInt(attr.Value, 10, 32)
			if err != nil {
				return 0, fmt.Errorf("failed to decode %s '%s': %s", name, attr.Value, err.Error())
			}
			return int32(val), nil
		}
	}
	return 0, io.EOF
}

func getManifestInfo(zip *apkparser.ZipReader) (minSdkVersion, targetSdkVersion, sandboxVersion int32, err error) {
	manifest := zip.File["AndroidManifest.xml"]
	if manifest == nil {
		return 1, 0, 1, nil
	}

	if err = manifest.Open(); err != nil {
		err = fmt.Errorf("failed to open AndroidManifest.xml: %s", err.Error())
		return
	}
	defer manifest.Close()

	for manifest.Next() {
		enc := sandboxVersionEncoder{1, 0, 1}
		if err = apkparser.ParseXml(manifest, &enc, nil); err != nil {
			err = fmt.Errorf("failed to parse AndroidManifest.xml: %s", err.Error())
			return
		}
		return enc.minSdkVersion, enc.targetSdkVersion, enc.sandboxVersion, nil
	}
	return
}
