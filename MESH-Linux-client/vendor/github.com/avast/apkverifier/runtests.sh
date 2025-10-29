#!/bin/sh
set -eux

if ! [ -d "apksig_for_tests" ]; then
    git clone --depth=1 -b android14-s1-release https://android.googlesource.com/platform/tools/apksig apksig_for_tests
else
    echo "Using cached apksig_for_test directory"
fi

export APKSIG_PATH=apksig_for_tests

if [ "$(go env GOARCH)" = "amd64" ]; then
    go test -race -parallel 8 -v ./...
else
    go test -v ./...
fi
