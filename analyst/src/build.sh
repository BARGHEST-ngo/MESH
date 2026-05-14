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

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

BUILD_DIR="$ROOT_DIR/tailscale"

echo "ROOT_DIR=$ROOT_DIR"
echo "BUILD_DIR=$BUILD_DIR"

if [ ! -d "$BUILD_DIR" ]; then
	echo "${RED}Error: tailscale submodule missing${NC}" >&2
	echo "Run: git submodule update --init --recursive" >&2
	exit 1
fi

tags="${TAGS:+$TAGS,}ts_omit_logtail"

cd "$BUILD_DIR"

echo "${GREEN}Building tailscaled daemon...${NC}"
go build \
	-tags "$tags" \
	-trimpath \
	-o "$ROOT_DIR/analyst/tailscaled" \
	./cmd/tailscaled

echo "${GREEN}Building tailscale CLI...${NC}"
go build \
	-tags "$tags" \
	-trimpath \
	-o "$ROOT_DIR/analyst/tailscale" \
	./cmd/tailscale

echo "${GREEN}Building mesh CLI...${NC}"
cd "$ROOT_DIR/analyst"
go build \
	-trimpath \
	-o "$ROOT_DIR/analyst/mesh-analyst" \
	.

echo "${GREEN}Build complete!${NC}"
