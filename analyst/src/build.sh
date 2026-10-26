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

echo "ROOT_DIR=$ROOT_DIR"

echo "${GREEN}Building mesh CLI...${NC}"
cd "$ROOT_DIR/analyst"
go build \
	-trimpath \
	-o "$ROOT_DIR/analyst/mesh" \
	.

echo "${GREEN}Build complete!${NC}"
