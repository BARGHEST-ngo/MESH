#!/usr/bin/env bash

set -eu

echo "Building meshcli with selected features..."
echo ""

# Use the new --custom-tailscaled flag we added to build_dist.sh
exec ./build_dist.sh --custom-tailscaled "$@"
