#!/usr/bin/env bash
#
# Build script for custom tailscaled with selected features
# This builds tailscaled with most features enabled except:
# - Cloud/platform integrations (aws, cloud, kube, synology, appconnectors)  
# - CLI features (cli, completion, cliconndiag)
# - Client update features (clientupdate)
# - Advanced/specialized features (c2n, oauthkey, outboundproxy, etc.)

set -eu

echo "Building meshcli with selected features..."
echo ""

# Use the new --custom-tailscaled flag we added to build_dist.sh
exec ./build_dist.sh --custom-tailscaled "$@"
