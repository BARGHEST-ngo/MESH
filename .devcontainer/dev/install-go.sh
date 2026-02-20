#!/bin/bash
# Setup script for Go in devcontainer
# This script installs Go

set -e

# Define variables
GO_VERSION="1.26.0"
GO_DOWNLOAD_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
GO_INSTALL_DIR="/usr/local"

# ============================================================================
# Install Go
# ============================================================================

echo ""
echo "Installing Go ${GO_VERSION}..."

# Check if Go is already installed with the correct version
if [ -x "${GO_INSTALL_DIR}/go/bin/go" ]; then
    INSTALLED_VERSION=$(${GO_INSTALL_DIR}/go/bin/go version | awk '{print $3}' | sed 's/go//')
    if [ "$INSTALLED_VERSION" = "$GO_VERSION" ]; then
        echo "Go ${GO_VERSION} is already installed."
    else
        echo "Found Go ${INSTALLED_VERSION}, but need ${GO_VERSION}. Reinstalling..."
        sudo rm -rf ${GO_INSTALL_DIR}/go
    fi
fi

# Install Go if not present or version mismatch
if [ ! -x "${GO_INSTALL_DIR}/go/bin/go" ]; then
    echo "Downloading Go ${GO_VERSION} from ${GO_DOWNLOAD_URL}..."

    # Create temporary directory for download
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Download Go
    curl -L -o go.tar.gz "$GO_DOWNLOAD_URL"

    # Extract to /usr/local
    echo "Extracting Go to ${GO_INSTALL_DIR}..."
    sudo tar -C ${GO_INSTALL_DIR} -xzf go.tar.gz

    # Cleanup
    cd -
    rm -rf "$TMP_DIR"

    echo "Go ${GO_VERSION} installed successfully."
fi

# Add Go to PATH for current session
export PATH="${GO_INSTALL_DIR}/go/bin:$PATH"

# Verify Go installation
echo "Verifying Go installation..."
which go && echo "✓ go found ($(go version))"

echo ""
echo "Go setup complete!"
