#!/bin/bash

set -e

echo "Building MESH CLI..."

# Build the CLI with specific tags (exclude most features but keep clientupdate)
./tool/go build -o meshcli ./cmd/meshcli

echo "CLI built correctly"
echo "CLI ready"
