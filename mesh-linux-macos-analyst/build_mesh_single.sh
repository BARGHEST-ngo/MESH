#!/bin/bash
# SPDX-License-Identifier: CC0-1.0
# Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.
#
# MESH Single Binary Build Script
# Builds a single "mesh" binary that acts as both daemon and CLI.
#
# Usage:
#   mesh up ...         → runs CLI (connect to network)
#   mesh status         → runs CLI (show status)
#   mesh daemon [flags] → starts the daemon
#   mesh down           → runs CLI (disconnect)

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "        ███╗   ███╗███████╗███████╗██╗  ██╗"
echo "        ████╗ ████║██╔════╝██╔════╝██║  ██║"
echo "        ██╔████╔██║█████╗  ███████╗███████║"
echo "        ██║╚██╔╝██║██╔══╝  ╚════██║██╔══██║"
echo "        ██║ ╚═╝ ██║███████╗███████║██║  ██║"
echo "        ╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝"
echo "        by Barghest.asia. No rights reserved."
echo -e "${NC}"
echo -e "${GREEN}        ── Single Binary Mode ──${NC}"
echo ""

BUILD_DIR="$(pwd)"
OUTPUT_DIR="${BUILD_DIR}"
BINARY_NAME="mesh"
CONFIG_DIR_TMP="/tmp/mesh"

VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --output|-o)
            BINARY_NAME="$2"
            shift 2
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Builds a single 'mesh' binary that includes both the daemon and CLI."
            echo ""
            echo "Options:"
            echo "  --output, -o NAME   Output binary name (default: mesh)"
            echo "  --verbose, -v       Enable verbose output"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "The single binary routes commands based on context:"
            echo "  mesh up ...           → CLI mode (connect to network)"
            echo "  mesh status           → CLI mode (show status)"
            echo "  mesh daemon [flags]   → Daemon mode (start the daemon)"
            echo "  mesh down             → CLI mode (disconnect)"
            echo ""
            echo "Examples:"
            echo "  $0                           # Build single binary"
            echo "  $0 --verbose                 # Build with verbose output"
            echo "  $0 -o mesh-analyst           # Build with custom name"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

if [ ! -f "go.mod" ] || [ ! -d "cmd/tailscaled" ] || [ ! -d "cmd/meshcli" ]; then
    echo -e "${RED}Error: This script must be run from the MESH-Linux-client directory${NC}"
    exit 1
fi

echo -e "${BLUE}Checking Go version...${NC}"
GO_VERSION=$(./tool/go version | awk '{print $3}' | sed 's/go//')
REQUIRED_GO_VERSION="1.25.3"
if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
    echo -e "${YELLOW}Warning: Go version $GO_VERSION detected. Recommended: $REQUIRED_GO_VERSION or higher${NC}"
fi

echo -e "${GREEN}Building single mesh binary (daemon + CLI)...${NC}"
if [ "$VERBOSE" = true ]; then
    ./tool/go build -v -tags ts_include_cli -o "${OUTPUT_DIR}/${BINARY_NAME}" ./cmd/tailscaled
else
    ./tool/go build -tags ts_include_cli -o "${OUTPUT_DIR}/${BINARY_NAME}" ./cmd/tailscaled
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Single binary built successfully: ${OUTPUT_DIR}/${BINARY_NAME}${NC}"
else
    echo -e "${RED}Failed to build single binary${NC}"
    exit 1
fi

echo -e "${BLUE}Setting up temporary config directory...${NC}"
mkdir -p "${CONFIG_DIR_TMP}"
echo -e "${GREEN}Created ${CONFIG_DIR_TMP}${NC}"

# Print summary
echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}Built binary:${NC}"
echo -e "  ${OUTPUT_DIR}/${BINARY_NAME}"
BINARY_SIZE=$(du -h "${OUTPUT_DIR}/${BINARY_NAME}" | cut -f1)
echo -e "  Size: ${BINARY_SIZE}"
echo ""
echo -e "${BLUE}Quick start (single terminal!):${NC}"
echo -e "  1. Start daemon:"
echo -e "     ${YELLOW}sudo ./${BINARY_NAME} daemon &${NC}"
echo ""
echo -e "  2. Connect to your network (in another terminal):"
echo -e "     ${YELLOW}sudo ./${BINARY_NAME} up --login-server=https://controlplane.com --auth-key=YOUR_KEY${NC}"
echo ""
echo -e "  3. Check status:"
echo -e "     ${YELLOW}sudo ./${BINARY_NAME} status${NC}"
echo ""
echo -e "  4. Disconnect:"
echo -e "     ${YELLOW}sudo ./${BINARY_NAME} down${NC}"
echo ""
echo -e "${BLUE}Daemon with custom flags:${NC}"
echo -e "  ${YELLOW}sudo ./${BINARY_NAME} daemon --socket=/tmp/tailscale/tailscaled.sock --state=/tmp/tailscale/tailscaled.state${NC}"
echo ""
echo ""

