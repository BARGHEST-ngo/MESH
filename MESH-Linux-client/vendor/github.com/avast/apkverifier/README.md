# apkverifier

[![GoDoc](https://godoc.org/github.com/avast/apkverifier?status.svg)](https://godoc.org/github.com/avast/apkverifier)
[![Build Status](https://travis-ci.org/avast/apkverifier.svg?branch=master)](https://travis-ci.org/avast/apkverifier)

APK signature verification, should support all algorithms and both scheme v1 and v2,
including downgrade attack protection.

**Works with Go 1.17 or higher.**

Documentation on [GoDoc](https://godoc.org/github.com/avast/apkverifier)

    go get github.com/avast/apkverifier

## Vendored stuff
Because Android can handle even broken x509 cerficates and ZIP files, apkverifier is using the ZipReader from apkparser
package and vendors `crypto/x509` in `internal/x509andr` and [github.com/fullsailor/pkcs7](https://github.com/fullsailor/pkcs7)
in the `fullsailor/pkcs7` folder.
The last two have some changes to handle some not-entirely-according-to-spec certificates.

## Example

```go
package main

import (
	"fmt"
	"github.com/avast/apkverifier"
	"os"
)

func main() {
	res, err := apkverifier.Verify(os.Args[1], nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Verification failed: %s\n", err.Error())
	}

	fmt.Printf("Verification scheme used: v%d\n", res.SigningSchemeId)
	cert, _ := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		fmt.Printf("No certificate found.\n")
	} else {
		fmt.Println(cert)
	}
}

```
