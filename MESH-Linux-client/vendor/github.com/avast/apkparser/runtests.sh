#!/bin/sh
set -eux

if [ -n "${GIMME_OS:-}" ] && [ "$GIMME_OS" != "linux" ]; then
    echo "Non-linux travis platform, skipping tests..."
    exit 0
fi

go test -v ./...
