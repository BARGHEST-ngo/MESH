# apkparser

[![GoDoc](https://godoc.org/github.com/avast/apkparser?status.svg)](https://godoc.org/github.com/avast/apkparser)
[![Build Status](https://travis-ci.org/avast/apkparser.svg?branch=master)](https://travis-ci.org/avast/apkparser)

APK AndroidManifest.xml and resources.arsc parsing.

**Works with Go 1.9 or higher.**

Documentation on [GoDoc](https://godoc.org/github.com/avast/apkparser)

    go get github.com/avast/apkparser

## ZipReader
Because Android can handle even broken ZIP archives, this packages has it's own zip reader,
based on archive/zip.

## axml2xml
A tool to extract AndroidManifest.xml and verify APK signature is also part of this repo.

    go get github.com/avast/apkparser
    go install github.com/avast/apkparser/axml2xml
    ./axml2xml -v application.apk

## Example

```go
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/avast/apkparser"
	"os"
)

func main() {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(os.Args[1], enc)
	if zipErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the APK: %s", zipErr.Error())
		os.Exit(1)
		return
	}

	if resErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse resources: %s", resErr.Error())
	}
	if manErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse AndroidManifest.xml: %s", manErr.Error())
		os.Exit(1)
		return
	}
	fmt.Println()
}
```
