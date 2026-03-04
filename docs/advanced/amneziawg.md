# AmneziaWG configuration

!!! danger
    AmneziaWG configuration only works on the analyst client right now. Android intergration is planned for February 2026.

    This guide explains how to configure AmneziaWG obfuscation in MESH for censorship resistance and Deep packet inspection (DPI) evasion.

## Overview

AmneziaWG is a backward-compatible fork of WireGuard that adds obfuscation parameters to evade DPI systems that block or throttle WireGuard traffic. This integration in MESH replaces the standard `wireguard-go` dependency with `amneziawg-go` while maintaining full compatibility with standard WireGuard peers.

### Why AmneziaWG?

Many countries and networks use Deep packet inspection to identify and block VPN traffic, including WireGuard. AmneziaWG makes WireGuard traffic look like random UDP packets, bypassing these restrictions.

**Use cases:**

- Operating in countries with aggressive censorship (China, Russia, Iran, etc.)
- Bypassing corporate firewalls that block VPN traffic
- Evading ISP throttling of VPN connections
- Conducting forensic investigations in hostile network environments

## Features

- **Backward Compatible**: When obfuscation is disabled (default), behaves identically to standard WireGuard
- **DPI Evasion**: Obfuscates WireGuard traffic to bypass DPI systems
- **Configurable**: Fine-grained control over obfuscation parameters
- **Zero Performance Impact**: When disabled, no overhead compared to standard WireGuard

## Obfuscation Parameters

AmneziaWG adds the following obfuscation parameters:

### Junk Packets

- **Jc** (uint8, 0-128): Number of junk packets sent before handshake
- **Jmin** (uint16): Minimum size of junk packets in bytes
- **Jmax** (uint16, ≤1280): Maximum size of junk packets in bytes

Junk packets are random data sent before the WireGuard handshake to confuse DPI systems.

### Handshake Padding

- **S1** (uint16, 15-150): Junk bytes added to handshake initiation
- **S2** (uint16, 15-150, S1+56≠S2): Junk bytes added to handshake response

Padding makes handshake packets variable-sized, preventing fingerprinting.

### Message Type Obfuscation

- **H1, H2, H3, H4** (uint32, ≥1, all different): Custom message type identifiers

Standard WireGuard uses fixed message types (1, 2, 3, 4). AmneziaWG allows custom values to evade signature-based detection.

## Configuration

### Configuration File Location

- **Linux**: `/etc/mesh/amneziawg.conf`
- **macOS**: `/usr/local/etc/mesh/amneziawg.conf`
- **Windows**: `C:\ProgramData\MESH\amneziawg.conf`
- **Android**: Configured via the MESH app settings

### Preset Configurations

#### Standard WireGuard Mode (Default)

No obfuscation, identical to standard WireGuard:

```ini
[Interface]
Jc = 0
Jmin = 0
Jmax = 0
S1 = 0
S2 = 0
H1 = 1
H2 = 2
H3 = 3
H4 = 4
```

#### Light Obfuscation (Minimal Overhead)

Minimal obfuscation for light censorship:

    [Interface]
    Jc = 3
    Jmin = 10
    Jmax = 50
    S1 = 15
    S2 = 20
    H1 = 5
    H2 = 6
    H3 = 7
    H4 = 8

```ini
[Interface]
Jc = 3
Jmin = 10
Jmax = 50
S1 = 15
S2 = 20
H1 = 5
H2 = 6
H3 = 7
H4 = 8
```

**Performance impact**: <5% overhead
**Use case**: Corporate firewalls, light ISP throttling

#### Balanced Configuration (Recommended)

Good balance between obfuscation and performance:

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

```ini
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
```

**Performance impact**: 5-10% overhead
**Use case**: Most censored networks, recommended default

#### Heavy Obfuscation (Maximum DPI Evasion)

Maximum obfuscation for aggressive censorship:

    [Interface]
    Jc = 10
    Jmin = 50
    Jmax = 1000
    S1 = 100
    S2 = 150
    H1 = 1234567
    H2 = 2345678
    H3 = 3456789
    H4 = 4567890

```ini
[Interface]
Jc = 10
Jmin = 50
Jmax = 1000
S1 = 100
S2 = 150
H1 = 1234567
H2 = 2345678
H3 = 3456789
H4 = 4567890
```

**Performance impact**: 10-15% overhead
**Use case**: Great Firewall of China, aggressive DPI systems

## Setup Instructions

### Linux/macOS Analyst client

#### 1. Create Configuration File

```bash
# Create config directory
sudo mkdir -p /etc/mesh

# Create config file with balanced preset
sudo cat > /etc/mesh/amneziawg.conf << 'EOF'
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

#### 2. Restart MESH Daemon

```bash
# If running manually
sudo pkill tailscaled-amnezia
sudo tailscaled-amnezia \
  --socket=/var/run/mesh/tailscaled.sock \
  --state=/var/lib/mesh/tailscaled.state \
  --statedir=/var/lib/mesh

# If running as systemd service
sudo systemctl restart mesh
```

#### 3. Verify Obfuscation

Check daemon logs for:

```
amneziawg: loading config from /etc/mesh/amneziawg.conf
amneziawg: obfuscation ENABLED - Jc=5 Jmin=50 Jmax=1000 S1=30 S2=40 H1=100 H2=200 H3=300 H4=400
```

### Android Client

#### 1. Open MESH App Settings

1. Open the MESH app
2. Tap the menu icon (☰)
3. Tap "Settings"
4. Tap "AmneziaWG Configuration"

#### 2. Select Preset

Choose from:

- **Disabled**: Standard WireGuard (default)
- **Light**: Minimal obfuscation
- **Balanced**: Recommended for most cases
- **Heavy**: Maximum DPI evasion
- **Custom**: Enter custom parameters

#### 3. Apply and Reconnect

1. Tap "Apply"
2. Disconnect and reconnect to the mesh

The new obfuscation settings will take effect.

## Verification

### Packet Capture

You can verify obfuscation by capturing packets:

```bash
# Find WireGuard port
sudo ss -unlp | grep tailscaled

# Capture packets (replace PORT with actual port)
sudo tcpdump -i any -n 'udp port PORT' -X -c 10
```

**Standard WireGuard packets** start with:

- `01 00 00 00` (handshake init)
- `02 00 00 00` (handshake response)
- `04 00 00 00` (transport data)

**AmneziaWG obfuscated packets** will have:

- Custom message types (H1, H2, H3, H4 values)
- Variable packet sizes (due to padding)
- Junk packets before handshake

### Connection test

Test if obfuscation helps bypass restrictions:

```bash
# Disable obfuscation
sudo rm /etc/mesh/amneziawg.conf
sudo systemctl restart mesh
sudo meshcli up --login-server=https://mesh.yourdomain.com

# Test connection
ping 100.64.X.X

# Enable obfuscation
sudo cat > /etc/mesh/amneziawg.conf << 'EOF'
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

sudo systemctl restart mesh
sudo meshcli up --login-server=https://mesh.yourdomain.com

# Test connection again
ping 100.64.X.X
```

If the connection works with obfuscation but not without, DPI is likely blocking standard WireGuard.

## References

- [AmneziaWG Specification](https://github.com/amnezia-vpn/amneziawg-go)
- [WireGuard Protocol](https://www.wireguard.com/protocol/)
- [MESH Project](https://github.com/barghest-ngo/mesh)

## Next steps

- **[Getting started](../getting-started/index.md)** - Set up your first mesh network
- **[Analyst client](../installation/analyst-client.md)** - Install the analyst client
- **[Endpoint client](../installation/endpoint-clients.md)** - Deploy to Android devices
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
