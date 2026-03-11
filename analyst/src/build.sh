#!/bin/sh
# Copyright (c) BARGHEST
# SPDX-License-Identifier: AGPL-3.0-or-later

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "${BLUE}"
echo "        ███╗   ███╗███████╗███████╗██╗  ██╗"
echo "        ████╗ ████║██╔════╝██╔════╝██║  ██║"
echo "        ██╔████╔██║█████╗  ███████╗███████║"
echo "        ██║╚██╔╝██║██╔══╝  ╚════██║██╔══██║"
echo "        ██║ ╚═╝ ██║███████╗███████║██║  ██║"
echo "        ╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝"
echo "        by Barghest.asia. No rights reserved."
echo "${NC}"
echo ""

GO_MOD_DIR="$(dirname "$(go env GOMOD)")"
if [ -z "$GO_MOD_DIR" ] || [ ! -d "$GO_MOD_DIR" ]; then
	echo "${RED}Error: Could not determine Go module directory${NC}" >&2
	exit 1
fi

tags="${TAGS:+$TAGS,}ts_omit_logtail" # ts_omit_aws,ts_omit_cloud,ts_omit_kube,ts_omit_synology,ts_omit_appconnectors,ts_omit_cli,ts_omit_completion,ts_omit_cliconndiag,ts_omit_clientupdate,ts_omit_c2n,ts_omit_oauthkey,ts_omit_outboundproxy,ts_omit_peerapiclient,ts_omit_peerapiserver,ts_omit_portlist,ts_omit_relayserver,ts_omit_wakeonlan,ts_omit_tap,ts_omit_bird"

BUILD_DIR="$GO_MOD_DIR/tailscale"
if [ ! -d "$BUILD_DIR" ]; then
	echo "${RED}Error: Build directory does not exist: $BUILD_DIR${NC}" >&2
	exit 1
fi
cd "$BUILD_DIR"

echo "${GREEN}Building mesh binary...${NC}"
go build -tags "$tags" -trimpath -o "$GO_MOD_DIR/analyst/mesh" ./cmd/tailscaled

echo "${GREEN}Build complete!${NC}"
