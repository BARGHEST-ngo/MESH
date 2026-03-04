# Analyst client (Linux/macOS)

The MESH analyst client is the forensic workstation component that connects to target devices over the mesh network. This guide covers installation, configuration, and usage of the analyst client on Linux and macOS.

## Overview

The analyst client consists of two main components:

- **`tailscaled-amnezia`**: The daemon that manages the WireGuard/AmneziaWG tunnel
- **`meshcli`**: The command-line interface for controlling the daemon

## System Requirements

### Linux

- **OS**: Ubuntu 20.04+, Debian 11+, Fedora 35+, or similar
- **Kernel**: 5.6+ (for WireGuard kernel module support)
- **RAM**: 512MB minimum, 1GB recommended
- **Disk**: 100MB for binaries and state
- **Network**: Internet connectivity, UDP port access (or HTTPS for DERP fallback)

### macOS

- **OS**: macOS 11 (Big Sur) or later
- **RAM**: 512MB minimum, 1GB recommended
- **Disk**: 100MB for binaries and state
- **Network**: Internet connectivity, UDP port access (or HTTPS for DERP fallback)

### Required Tools

- **Go**: 1.21 or later (for building from source)
- **Git**: For cloning the repository
- **ADB**: Android Debug Bridge (for Android forensics)
- **AndroidQF**: Optional, for automated spyware detection
- **MVT**: Optional, Mobile Verification Toolkit

## Installation

### Option 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh/mesh-linux-macos-analyst

# Build the binaries
./build_mesh.sh

# Optionally install to system paths
sudo cp meshcli /usr/local/bin/
sudo cp tailscaled-amnezia /usr/local/bin/
```

### Option 2: Download Pre-built Binaries

```bash
# Download the latest release
wget https://github.com/BARGHEST-ngo/mesh/releases/latest/download/mesh-analyst-linux-amd64.tar.gz

# Extract
tar -xzf mesh-analyst-linux-amd64.tar.gz

# Move to system path
sudo mv meshcli tailscaled-amnezia /usr/local/bin/
```

### Install ADB (Android Forensics)

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install android-tools-adb
```

**macOS:**

```bash
brew install android-platform-tools
```

**Fedora:**

```bash
sudo dnf install android-tools
```

### Install AndroidQF (Optional)

```bash
pip3 install androidqf
```

### Install MVT (Optional)

```bash
pip3 install mvt
```

## Configuration

### Create State Directory

```bash
sudo mkdir -p /var/lib/mesh
sudo mkdir -p /var/run/mesh
```

### Configure AmneziaWG (Optional)

For censorship resistance, create an AmneziaWG configuration:

```bash
sudo mkdir -p /etc/mesh
sudo cat > /etc/mesh/amneziawg.conf << EOF
[Interface]
Jc = 5
Jmin = 50
Jmax = 1000
S1 = 30
S2 = 40
H1 = 100
H2 = 200
H3 = 300
H4 = 400
EOF
```

See [AmneziaWG Configuration](../advanced/amneziawg.md) for details.

## Usage

### Starting the Daemon

#### Manual Start

```bash
sudo tailscaled-amnezia \
  --socket=/var/run/mesh/tailscaled.sock \
  --state=/var/lib/mesh/tailscaled.state \
  --statedir=/var/lib/mesh
```

#### Systemd Service (Linux)

Create a systemd service file:

```bash
sudo cat > /etc/systemd/system/mesh.service << EOF
[Unit]
Description=MESH Analyst client
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/tailscaled-amnezia --socket=/var/run/mesh/tailscaled.sock --state=/var/lib/mesh/tailscaled.state --statedir=/var/lib/mesh
Restart=on-failure
Environment="TS_DEBUG_TRIM_WIREGUARD=false"

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable mesh
sudo systemctl start mesh
```

### Connecting to the Control plane

```bash
# Connect with pre-auth key
sudo meshcli up \
  --login-server=https://your-control-plane.com \
  --authkey=your-preauth-key \
  --accept-dns=false

# Or connect with interactive authentication
sudo meshcli up --login-server=https://your-control-plane.com
```

### Checking Status

```bash
# View connection status
sudo meshcli status

# List all peers in the mesh
sudo meshcli status --peers

# View detailed peer information
sudo meshcli status --json | jq
```

### Disconnecting

```bash
# Disconnect from the mesh
sudo meshcli down
```

## Forensic Workflows

### Connecting to Android Device

Once the Android endpoint is connected to the mesh:

```bash
# Find the device's mesh IP
sudo meshcli status --peers

# Connect via ADB
adb connect 100.64.X.X:5555

# Verify connection
adb devices
```

### Collecting Artifacts

#### Bug Report

```bash
adb bugreport bugreport-$(date +%Y%m%d-%H%M%S).zip
```

#### System Dump

```bash
adb shell dumpsys > dumpsys-$(date +%Y%m%d-%H%M%S).txt
```

#### Logcat

```bash
adb logcat -d > logcat-$(date +%Y%m%d-%H%M%S).txt
```

#### Package List

```bash
adb shell pm list packages -f > packages-$(date +%Y%m%d-%H%M%S).txt
```

### Running AndroidQF

```bash
# Run AndroidQF on the remote device
androidqf --adb 100.64.X.X:5555 --output ./artifacts/

# Check for specific IOCs
androidqf --adb 100.64.X.X:5555 --iocs iocs.json --output ./artifacts/
```

### Running MVT

```bash
# Run MVT spyware check
mvt-android check-adb --serial 100.64.X.X:5555 --output ./mvt-results/

# Check against specific IOCs
mvt-android check-adb --serial 100.64.X.X:5555 --iocs pegasus.stix2 --output ./mvt-results/
```

## Advanced Features

### Subnet Routing

Access devices on the endpoint's local network:

```bash
# The endpoint must advertise its subnet
# On analyst, accept the route
sudo meshcli up --accept-routes
```

### Exit Node

Route your internet traffic through an endpoint:

```bash
# Use a peer as exit node
sudo meshcli up --exit-node=100.64.X.X

# Stop using exit node
sudo meshcli up --exit-node=
```

### Kill Switch

Block all non-MESH traffic:

```bash
# Enable kill switch (Linux only)
sudo iptables -A OUTPUT -o tun0 -j ACCEPT
sudo iptables -A OUTPUT -j DROP
```

### SSH over MESH

```bash
# SSH to endpoint (if SSH server is running)
ssh user@100.64.X.X
```

## Troubleshooting

### Daemon Won't Start

**Check if another instance is running:**

```bash
ps aux | grep tailscaled
sudo pkill tailscaled-amnezia
```

**Check socket permissions:**

```bash
sudo rm -rf /var/run/mesh/tailscaled.sock
```

### Can't Connect to Control plane

**Verify control plane URL:**

```bash
curl https://your-control-plane.com/health
```

**Check logs:**

```bash
# If running as systemd service
sudo journalctl -u mesh -f

# If running manually, check terminal output
```

### No Peer-to-Peer Connection

**Check if UDP is blocked:**

```bash
# Force DERP relay
export TS_DEBUG_ALWAYS_USE_DERP=true
sudo -E meshcli up --login-server=https://your-control-plane.com
```

**Enable AmneziaWG obfuscation:**

See [AmneziaWG Configuration](../advanced/amneziawg.md).

### ADB Connection Fails

**Verify device is online:**

```bash
sudo meshcli status --peers
ping 100.64.X.X
```

**Reset ADB:**

```bash
adb kill-server
adb start-server
adb connect 100.64.X.X:5555
```

**Check ADB port:**

```bash
# On Android device (via USB)
adb shell getprop service.adb.tcp.port
# Should return 5555
```

## CLI reference

See [CLI reference](../reference/cli-reference.md) for complete command documentation.

## Security Best practices

1. **Use strong pre-auth keys**: Generate keys with short expiration times
2. **Limit ACLs**: Only allow necessary connections between nodes
3. **Rotate keys**: Regularly regenerate WireGuard keys
4. **Monitor connections**: Regularly check `meshcli status --peers`
5. **Use AmneziaWG**: Enable obfuscation in hostile environments
6. **Secure the control plane**: Use HTTPS, strong authentication
7. **Audit logs**: Review control plane and daemon logs regularly

## Performance Tuning

### Increase MTU

```bash
# Check current MTU
ip link show tun0

# Increase MTU (Linux)
sudo ip link set tun0 mtu 1420
```

### Disable DNS

If you don't need MagicDNS:

```bash
sudo meshcli up --accept-dns=false
```

### Prefer IPv6

```bash
sudo meshcli up --prefer-ipv6
```

## Uninstallation

```bash
# Stop the daemon
sudo systemctl stop mesh
sudo systemctl disable mesh

# Remove binaries
sudo rm /usr/local/bin/meshcli
sudo rm /usr/local/bin/tailscaled-amnezia

# Remove state
sudo rm -rf /var/lib/mesh
sudo rm -rf /var/run/mesh
sudo rm -rf /etc/mesh

# Remove systemd service
sudo rm /etc/systemd/system/mesh.service
sudo systemctl daemon-reload
```

## Next steps

- **[User guide](../user-guide/user-guide.md)** - Learn forensic workflows
- **[AmneziaWG Configuration](../advanced/amneziawg.md)** - Enable censorship resistance
- **[CLI reference](../reference/cli-reference.md)** - Complete command documentation
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
