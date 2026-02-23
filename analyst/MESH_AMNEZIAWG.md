# AmneziaWG Integration

This document describes the AmneziaWG integration in MESH, which provides DPI (Deep Packet Inspection) evasion capabilities.

## Overview

AmneziaWG is a backward-compatible fork of WireGuard that adds obfuscation parameters to evade DPI systems that block or throttle WireGuard traffic. This integration in MESH replaces the standard `wireguard-go` dependency with `amneziawg-go` while maintaining full compatibility with standard WireGuard peers.

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

### Handshake Padding
- **S1** (uint16, 15-150): Junk bytes added to handshake initiation
- **S2** (uint16, 15-150, S1+56≠S2): Junk bytes added to handshake response

### Message Type Obfuscation
- **H1, H2, H3, H4** (uint32, ≥1, all different): Custom message type identifiers

## Configuration

### Configuration File Location

- **Linux**: `/etc/mesh/amneziawg.conf`
- **macOS**: `/usr/local/etc/mesh/amneziawg.conf`
- **Windows**: `C:\ProgramData\MESH\amneziawg.conf`
- **Android**: `/data/data/com.barghest.mesh/files/amneziawg.conf`

### Example Configurations

#### Standard WireGuard Mode (Default)
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

#### Recommended Balanced Configuration
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

#### Heavy Obfuscation (Maximum DPI Evasion)
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

## Building

### Quick Build

```bash
cd MESH-Linux-client
./build_mesh.sh
```

### Build with Config Installation

```bash
cd MESH-Linux-client
sudo ./build_mesh.sh --install-config
```

### Build with Systemd Service

```bash
cd MESH-Linux-client
sudo ./build_mesh.sh --setup-systemd
```

## Usage

### 1. Start the Daemon

```bash
# Set environment variable to prevent peer trimming
export TS_DEBUG_TRIM_WIREGUARD=false

# Start daemon
sudo -E ./tailscaled-amnezia \
  --socket=/tmp/mesh/tailscaled.sock \
  --state=/tmp/mesh/tailscaled.state \
  --statedir=/tmp/mesh
```

### 2. Connect to Your Network

```bash
# Connect to headscale server
sudo ./meshcli up \
  --login-server=https://controlplane.com \
  -auth-key 12387129837918 \
  --accept-dns=false ## VERY IMPORTANT, using the control server DNS will often cause issues on standard networks
```

### 3. Enable Obfuscation (Optional)

```bash
# Copy example config
sudo cp /etc/mesh/amneziawg.conf.example /etc/mesh/amneziawg.conf

# Edit config with your preferred parameters
sudo nano /etc/mesh/amneziawg.conf

# Restart daemon to apply changes
```

### 4. Verify Obfuscation

Check daemon logs for:
```
amneziawg: loading config from /etc/mesh/amneziawg.conf
amneziawg: obfuscation ENABLED - Jc=5 Jmin=50 Jmax=1000 S1=30 S2=40 H1=100 H2=200 H3=300 H4=400
```

Or capture packets to verify obfuscation:
```bash
# Find WireGuard port
sudo ss -unlp | grep tailscaled

# Capture packets (replace PORT with actual port)
sudo tcpdump -i any -n 'udp port PORT' -X -c 10
```

Standard WireGuard packets start with `01 00 00 00` (handshake init) or `04 00 00 00` (transport data).
AmneziaWG obfuscated packets will have custom message types (H1, H2, H3, H4 values).

## Troubleshooting

### DNS Issues

If internet stops working when daemon is running:

```bash
sudo ./meshcli set --accept-dns=false
```

### Peer Not Connecting

Ensure `TS_DEBUG_TRIM_WIREGUARD=false` is set:

```bash
export TS_DEBUG_TRIM_WIREGUARD=false
sudo -E ./tailscaled-amnezia ...
```

### Obfuscation Not Working

1. Check config file exists:
   ```bash
   ls -la /etc/mesh/amneziawg.conf
   ```

2. Check daemon logs for "obfuscation ENABLED" message

3. Verify config syntax (must have `[Interface]` section)

4. Restart daemon after config changes

### Direct Connection Not Established

Disable DERP relay to force direct connections:

```bash
export TS_DEBUG_ALWAYS_USE_DERP=false
sudo -E ./tailscaled-amnezia ...
```

## Technical Details

### Implementation

- **Dependency**: Replaced `golang.zx2c4.com/wireguard` with `github.com/amnezia-vpn/amneziawg-go`
- **Config Parser**: `MESH-Linux-client/mesh/amneziawg/config.go`
- **Integration**: `MESH-Linux-client/wgengine/userspace.go`
- **Compatibility**: GSO/GRO shims in `MESH-Linux-client/mesh/amneziawg/compat.go`

### Backward Compatibility

When all obfuscation parameters are set to default values (Jc=0, S1=0, S2=0, H1=1, H2=2, H3=3, H4=4), AmneziaWG produces identical packets to standard WireGuard, ensuring compatibility with:

- Standard WireGuard peers
- Existing WireGuard infrastructure

## Security Considerations

1. **Obfuscation ≠ Encryption**: AmneziaWG obfuscates traffic patterns but does not add additional encryption. WireGuard's encryption remains unchanged.

2. **Parameter Selection**: Choose parameters that balance DPI evasion with performance. Higher values provide better obfuscation but increase overhead.

3. **Network Compatibility**: Some networks may drop packets with unusual sizes. Test your configuration in your target environment.

4. **Peer Coordination**: All peers in a connection must use compatible AmneziaWG parameters. Mismatched parameters will prevent connection.

## References

- [AmneziaWG Specification](https://github.com/amnezia-vpn/amneziawg-go)
- [WireGuard Protocol](https://www.wireguard.com/protocol/)
- [MESH Project](https://github.com/barghest/mesh)

## License

MESH is licensed under the [GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later)](LICENSE).

Portions of this software are a derivative work of [Tailscale](https://github.com/tailscale/tailscale), which is licensed under the BSD 3-Clause License. The original Tailscale copyright and license are preserved in accordance with the BSD-3-Clause requirements. AmneziaWG/Wireguard code is licensed under MIT license. See `.licenses/` for details.

All modifications and additions by BARGHEST are Copyright (c) BARGHEST and licensed under AGPL-3.0-or-later.
