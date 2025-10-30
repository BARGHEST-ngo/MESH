#!/bin/bash
# SPDX-License-Identifier: CC0-1.0
# Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.
#
# MESH Build Script
# Builds tailscaled-amnezia daemon and meshcli CLI tool
# Sets up configuration directories and example files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print banner
echo -e "${BLUE}"
echo "        ███╗   ███╗███████╗███████╗██╗  ██╗"
echo "        ████╗ ████║██╔════╝██╔════╝██║  ██║"
echo "        ██╔████╔██║█████╗  ███████╗███████║"
echo "        ██║╚██╔╝██║██╔══╝  ╚════██║██╔══██║"
echo "        ██║ ╚═╝ ██║███████╗███████║██║  ██║"
echo "        ╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝"
echo "        by Barghest.asia. No rights reserved."
echo -e "${NC}"
echo ""

# Configuration
BUILD_DIR="$(pwd)"
OUTPUT_DIR="${BUILD_DIR}"
DAEMON_NAME="tailscaled-amnezia"
CLI_NAME="meshcli"
CONFIG_DIR_SYSTEM="/etc/mesh"
CONFIG_DIR_TMP="/tmp/mesh"

# Parse command line arguments
INSTALL_CONFIG=false
SETUP_SYSTEMD=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --install-config)
            INSTALL_CONFIG=true
            shift
            ;;
        --setup-systemd)
            SETUP_SYSTEMD=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --install-config    Install example AmneziaWG config to /etc/mesh/"
            echo "  --setup-systemd     Create systemd service file (requires root)"
            echo "  --verbose, -v       Enable verbose output"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                           # Build daemon and CLI"
            echo "  $0 --install-config          # Build and install example config"
            echo "  sudo $0 --setup-systemd      # Build and setup systemd service"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/tailscaled" ] || [ ! -d "cmd/meshcli" ]; then
    echo -e "${RED}Error: This script must be run from the MESH-Linux-client directory${NC}"
    exit 1
fi

# Check Go version
echo -e "${BLUE}Checking Go version...${NC}"
GO_VERSION=$(./tool/go version | awk '{print $3}' | sed 's/go//')
REQUIRED_GO_VERSION="1.25.3"
if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
    echo -e "${YELLOW}Warning: Go version $GO_VERSION detected. Recommended: $REQUIRED_GO_VERSION or higher${NC}"
fi

# Build the daemon
echo -e "${GREEN}Building tailscaled-amnezia daemon...${NC}"
if [ "$VERBOSE" = true ]; then
    ./tool/go build -v -o "${OUTPUT_DIR}/${DAEMON_NAME}" ./cmd/tailscaled
else
    ./tool/go build -o "${OUTPUT_DIR}/${DAEMON_NAME}" ./cmd/tailscaled
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Daemon built successfully: ${OUTPUT_DIR}/${DAEMON_NAME}${NC}"
else
    echo -e "${RED}✗ Failed to build daemon${NC}"
    exit 1
fi

# Build the CLI
echo -e "${GREEN}Building meshcli CLI...${NC}"
if [ "$VERBOSE" = true ]; then
    ./tool/go build -v -o "${OUTPUT_DIR}/${CLI_NAME}" ./cmd/meshcli
else
    ./tool/go build -o "${OUTPUT_DIR}/${CLI_NAME}" ./cmd/meshcli
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI built successfully: ${OUTPUT_DIR}/${CLI_NAME}${NC}"
else
    echo -e "${RED}✗ Failed to build CLI${NC}"
    exit 1
fi

# Create temporary config directory
echo -e "${BLUE}Setting up temporary config directory...${NC}"
mkdir -p "${CONFIG_DIR_TMP}"
echo -e "${GREEN}✓ Created ${CONFIG_DIR_TMP}${NC}"

# Install example config if requested
if [ "$INSTALL_CONFIG" = true ]; then
    echo -e "${BLUE}Installing example AmneziaWG configuration...${NC}"
    
    if [ ! -f "amneziawg.conf.example" ]; then
        echo -e "${RED}Error: amneziawg.conf.example not found${NC}"
        exit 1
    fi
    
    # Create system config directory (requires root)
    if [ "$EUID" -ne 0 ]; then
        echo -e "${YELLOW}Warning: Root privileges required to install to ${CONFIG_DIR_SYSTEM}${NC}"
        echo -e "${YELLOW}Run with sudo to install system-wide configuration${NC}"
    else
        mkdir -p "${CONFIG_DIR_SYSTEM}"
        
        # Check if config already exists
        if [ -f "${CONFIG_DIR_SYSTEM}/amneziawg.conf" ]; then
            echo -e "${YELLOW}Config already exists at ${CONFIG_DIR_SYSTEM}/amneziawg.conf${NC}"
            echo -e "${YELLOW}Backing up to ${CONFIG_DIR_SYSTEM}/amneziawg.conf.backup${NC}"
            cp "${CONFIG_DIR_SYSTEM}/amneziawg.conf" "${CONFIG_DIR_SYSTEM}/amneziawg.conf.backup"
        fi
        
        # Install example config as disabled by default
        cp "amneziawg.conf.example" "${CONFIG_DIR_SYSTEM}/amneziawg.conf.example"
        echo -e "${GREEN}✓ Installed example config to ${CONFIG_DIR_SYSTEM}/amneziawg.conf.example${NC}"
        echo -e "${BLUE}To enable obfuscation, copy the example and edit:${NC}"
        echo -e "  sudo cp ${CONFIG_DIR_SYSTEM}/amneziawg.conf.example ${CONFIG_DIR_SYSTEM}/amneziawg.conf"
        echo -e "  sudo nano ${CONFIG_DIR_SYSTEM}/amneziawg.conf"
    fi
fi

# Setup systemd service if requested
if [ "$SETUP_SYSTEMD" = true ]; then
    echo -e "${BLUE}Setting up systemd service...${NC}"
    
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}Error: Root privileges required to setup systemd service${NC}"
        exit 1
    fi
    
    SYSTEMD_SERVICE="/etc/systemd/system/mesh-tailscaled.service"
    
    cat > "${SYSTEMD_SERVICE}" <<EOF
[Unit]
Description=MESH Tailscaled Daemon (using AmneziaWG)
Documentation=https://github.com/barghest/mesh
After=network-pre.target
Wants=network-pre.target

[Service]
Type=notify
ExecStartPre=/bin/mkdir -p /var/run/mesh /var/cache/mesh /var/lib/mesh
ExecStart=${OUTPUT_DIR}/${DAEMON_NAME} --state=/var/lib/mesh/tailscaled.state --socket=/var/run/mesh/tailscaled.sock
Restart=on-failure
RestartSec=5
Environment="TS_DEBUG_TRIM_WIREGUARD=false"

[Install]
WantedBy=multi-user.target
EOF
    
    echo -e "${GREEN}✓ Created systemd service: ${SYSTEMD_SERVICE}${NC}"
    echo -e "${BLUE}To enable and start the service:${NC}"
    echo -e "  sudo systemctl daemon-reload"
    echo -e "  sudo systemctl enable mesh-tailscaled"
    echo -e "  sudo systemctl start mesh-tailscaled"
fi

# Print summary
echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}Built binaries:${NC}"
echo -e "  Daemon: ${OUTPUT_DIR}/${DAEMON_NAME}"
echo -e "  CLI:    ${OUTPUT_DIR}/${CLI_NAME}"
echo ""
echo -e "${BLUE}Configuration directories:${NC}"
echo -e "  Temporary: ${CONFIG_DIR_TMP}"
if [ "$INSTALL_CONFIG" = true ] && [ "$EUID" -eq 0 ]; then
    echo -e "  System:    ${CONFIG_DIR_SYSTEM}"
fi
echo ""
echo -e "${BLUE}Quick start:${NC}"
echo -e "  1. Start daemon (temporary):"
echo -e "     ${YELLOW}export TS_DEBUG_TRIM_WIREGUARD=false${NC}"
echo -e "     ${YELLOW}sudo -E ./${DAEMON_NAME} --socket=/tmp/mesh/tailscaled.sock --state=/tmp/mesh/tailscaled.state --statedir=/tmp/mesh${NC}"
echo ""
echo -e "  2. Connect to your network (in another terminal):"
echo -e "     ${YELLOW}sudo ./${CLI_NAME} up --login-server=https://your-headscale-server.com --accept-dns=false${NC}"
echo ""
echo -e "  3. Check status:"
echo -e "     ${YELLOW}sudo ./${CLI_NAME} status${NC}"
echo ""
echo -e "${BLUE}AmneziaWG obfuscation:${NC}"
if [ "$INSTALL_CONFIG" = true ] && [ "$EUID" -eq 0 ]; then
    echo -e "  Config example: ${CONFIG_DIR_SYSTEM}/amneziawg.conf.example"
    echo -e "  To enable: sudo cp ${CONFIG_DIR_SYSTEM}/amneziawg.conf.example ${CONFIG_DIR_SYSTEM}/amneziawg.conf"
else
    echo -e "  Example config: ./amneziawg.conf.example"
    echo -e "  To install: sudo $0 --install-config"
fi
echo -e "  Edit config and restart daemon to enable obfuscation"
echo ""
echo -e "${BLUE}For more information:${NC}"
echo -e "  Run: $0 --help"
echo ""

