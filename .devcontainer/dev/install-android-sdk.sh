#!/bin/bash
# Setup script for Android SDK in devcontainer
# This script installs Android SDK, NDK, and required build tools

set -e

# Define variables
ANDROID_SDK_ROOT="${ANDROID_SDK_ROOT:-/usr/local/android-sdk}"
ANDROID_HOME="${ANDROID_HOME:-$ANDROID_SDK_ROOT}"
ANDROID_TOOLS_URL="https://dl.google.com/android/repository/commandlinetools-linux-9477386_latest.zip"
ANDROID_TOOLS_SUM="bd1aa17c7ef10066949c88dc6c9c8d536be27f992a1f3b5a584f9bd2ba5646a0"

# SDK packages to install (matching Makefile requirements)
ANDROID_SDK_PACKAGES=(
    "platforms;android-34"
    "extras;android;m2repository"
    "ndk;23.1.7779620"
    "platform-tools"
    "build-tools;34.0.0"
)

# ============================================================================
# Install Android SDK
# ============================================================================

echo ""
echo "Setting up Android SDK environment..."

# Create SDK directory
echo "Creating Android SDK directory at $ANDROID_SDK_ROOT..."
sudo mkdir -p "$ANDROID_SDK_ROOT"
sudo chown -R $(whoami):$(whoami) "$ANDROID_SDK_ROOT"

# Check if SDK is already installed
if [ -f "$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" ]; then
    echo "Android SDK command-line tools already installed."
else
    echo "Downloading Android command-line tools..."
    mkdir -p "$ANDROID_SDK_ROOT/tmp"
    cd "$ANDROID_SDK_ROOT/tmp"
    
    # Download command-line tools
    curl --silent -O -L "$ANDROID_TOOLS_URL"
    
    # Verify checksum
    echo "$ANDROID_TOOLS_SUM  commandlinetools-linux-9477386_latest.zip" | sha256sum -c -
    
    # Extract and move to correct location
    unzip -q commandlinetools-linux-9477386_latest.zip
    mkdir -p "$ANDROID_SDK_ROOT/cmdline-tools"
    mv cmdline-tools "$ANDROID_SDK_ROOT/cmdline-tools/latest"
    
    # Cleanup
    cd -
    rm -rf "$ANDROID_SDK_ROOT/tmp"
    
    echo "Android command-line tools installed successfully."
fi

# Accept licenses
echo "Accepting Android SDK licenses..."
yes | "$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --licenses > /dev/null 2>&1 || true

# Update SDK manager
echo "Updating SDK manager..."
"$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --update

# Install required SDK packages
echo "Installing Android SDK packages..."
for package in "${ANDROID_SDK_PACKAGES[@]}"; do
    echo "  Installing $package..."
    "$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" "$package"
done

# Verify installations
echo ""
echo "Verifying installations..."
echo "Installed SDK packages:"
"$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --list_installed

# Set up environment for current session
export ANDROID_SDK_ROOT="$ANDROID_SDK_ROOT"
export ANDROID_HOME="$ANDROID_HOME"
export PATH="$ANDROID_SDK_ROOT/cmdline-tools/latest/bin:$ANDROID_SDK_ROOT/platform-tools:$ANDROID_SDK_ROOT/build-tools/34.0.0:$PATH"

# Verify key tools are accessible
echo ""
echo "Verifying tools are accessible..."
which sdkmanager && echo "✓ sdkmanager found"
which adb && echo "✓ adb found"
which avdmanager && echo "✓ avdmanager found"

# Display NDK location
NDK_ROOT="$ANDROID_SDK_ROOT/ndk/23.1.7779620"
if [ -d "$NDK_ROOT" ]; then
    echo "✓ NDK 23.1.7779620 installed at: $NDK_ROOT"
else
    echo "⚠ NDK not found at expected location: $NDK_ROOT"
fi

echo ""
echo "Android SDK setup complete!"
echo ""
echo "Installed versions:"
echo "  Java: $(java -version 2>&1 | head -n 1)"
echo "  Gradle: $(gradle --version | grep Gradle | awk '{print $2}')"
echo ""
echo "Environment variables:"
echo "  ANDROID_SDK_ROOT: $ANDROID_SDK_ROOT"
echo "  ANDROID_HOME: $ANDROID_HOME"
