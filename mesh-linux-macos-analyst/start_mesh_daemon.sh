#!/bin/bash
# Script to start the MESH daemon in userspace mode

set -e
MESH_DIR="/tmp/tailscale"
SOCKET_PATH="/tmp/tailscale/tailscaled.sock"
STATE_PATH="/tmp/tailscale/tailscaled.state"
echo "Starting Tailscale-AmenziaWG daemon..."
mkdir -p "$MESH_DIR"
rm -f "$SOCKET_PATH"
echo "Socket: $SOCKET_PATH"
echo "State: $STATE_PATH"
echo ""
# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run with sudo for TUN interface access"
    echo "Usage: sudo ./start_mesh_daemon.sh"
    exit 1
fi

./tailscaled-amnezia \
    --socket="$SOCKET_PATH" \
    --state="$STATE_PATH" \
    --statedir="$MESH_DIR" \
    --verbose=1

echo "Tailscale-AmenziaWG daemon stopped."