package signingblock

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func getLenghtPrefixedSlice(r *bytes.Buffer) (*bytes.Buffer, error) {
	if r.Len() < 4 {
		return nil, fmt.Errorf("Remaining buffer too short to contain length of length-prefixed field. Remaining: %d", r.Len())
	}

	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.New("negative length")
	} else if int(length) > r.Len() {
		return nil, fmt.Errorf("Length-prefixed field longer than remaining buffer. "+
			"Field length: %d, remaining: %d", length, r.Len())
	}
	return bytes.NewBuffer(r.Next(int(length))), nil
}

func verifySignature(publicKey interface{}, algo SignatureAlgorithm, signedDataBytes, signature []byte) error {
	switch algo {
	case SigRsaPssWithSha256:
		hashed := sha256.Sum256(signedDataBytes)
		return rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
			SaltLength: 256 / 8,
		})
	case SigRsaPssWithSha512:
		hashed := sha512.Sum512(signedDataBytes)
		return rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA512, hashed[:], signature, &rsa.PSSOptions{
			SaltLength: 512 / 8,
		})
	case SigRsaPkcs1V15WithSha256, SigVerityRsaPkcs1V15WithSha256:
		hashed := sha256.Sum256(signedDataBytes)
		return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature)
	case SigRsaPkcs1V15WithSha512:
		hashed := sha512.Sum512(signedDataBytes)
		return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA512, hashed[:], signature)
	case SigEcdsaWithSha256, SigEcdsaWithSha512, SigDsaWithSha256, SigVerityEcdsaWithSha256, SigVerityDsaWithSha256:
		var params []*big.Int
		if _, err := asn1.Unmarshal(signature, &params); err != nil {
			return fmt.Errorf("failed to unmarshal ECDSA signature: %s", err.Error())
		}

		var hashed []byte
		switch algo {
		case SigEcdsaWithSha256, SigDsaWithSha256, SigVerityEcdsaWithSha256, SigVerityDsaWithSha256:
			h := sha256.Sum256(signedDataBytes)
			hashed = h[:]
		case SigEcdsaWithSha512:
			h := sha512.Sum512(signedDataBytes)
			hashed = h[:]
		default:
			panic("Unandled")
		}

		if algo == SigDsaWithSha256 || algo == SigVerityDsaWithSha256 {
			k := publicKey.(*dsa.PublicKey)
			hashed = hashed[:k.Q.BitLen()/8]
			if !dsa.Verify(k, hashed, params[0], params[1]) {
				return errors.New("DSA verification failed.")
			}
		} else {
			if !ecdsa.Verify(publicKey.(*ecdsa.PublicKey), hashed, params[0], params[1]) {
				return errors.New("ECDSA verification failed.")
			}
		}
	default:
		return fmt.Errorf("unhandled signature type: 0x%04x", int32(algo))
	}

	return nil
}

func PkixNameToString(n *pkix.Name) string {
	var buf bytes.Buffer

	if len(n.Country) != 0 {
		fmt.Fprintf(&buf, "C=%s, ", strings.Join(n.Country, ";"))
	}
	if len(n.Province) != 0 {
		fmt.Fprintf(&buf, "ST=%s, ", strings.Join(n.Province, ";"))
	}
	if len(n.Locality) != 0 {
		fmt.Fprintf(&buf, "L=%s, ", strings.Join(n.Locality, ";"))
	}
	if len(n.Organization) != 0 {
		fmt.Fprintf(&buf, "O=%s, ", strings.Join(n.Organization, ";"))
	}
	if len(n.OrganizationalUnit) != 0 {
		fmt.Fprintf(&buf, "OU=%s, ", strings.Join(n.OrganizationalUnit, ";"))
	}
	if len(n.CommonName) != 0 {
		fmt.Fprintf(&buf, "CN=%s, ", n.CommonName)
	}

	// Remove last ', '
	if buf.Len() != 0 {
		buf.Truncate(buf.Len() - 2)
	}

	return buf.String()
}
