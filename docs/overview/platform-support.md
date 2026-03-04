# Platform support

MESH is designed to work across multiple platforms to support diverse forensic investigation scenarios. This page details current support status and planned features.

## Current support (public alpha)

The following platforms are currently supported in the public alpha release:

| Platform | Status | Role | Minimum Version |
|----------|--------|------|-----------------|
| **Android** | ✅ Supported | Endpoint client | Android 8.0 (Oreo) |
| **Linux** | ✅ Supported | Analyst client | Ubuntu 20.04+, Debian 11+ |
| **macOS** | ✅ Supported | Analyst client | macOS 11 (Big Sur) |
| **Control plane** | ✅ Supported | Self-hosted Server | Docker 20.10+ |

### Android (endpoint client)

**Supported Versions:**

- Android 8.0 (Oreo) and later
- Tested on Android 8.0 through Android 14

**Features:**

- WireGuard VPN
- Guided session for ADB-over-WiFi enablement
- Background operation
- Battery optimization
- VPN kill-switch
- Connection status UI

**Requirements:**

- WiFi connectivity
- VPN permission (granted during setup)
- At least 50MB free storage

**Installation:**

- Build from source (see [Endpoint client guide](../installation/endpoint-clients.md))
- Pre-built APKs available on [GitHub Releases](https://github.com/BARGHEST-ngo/mesh/releases)

### Linux (analyst client)

**Supported Distributions:**

- Ubuntu 20.04 LTS and later
- Debian 11 (Bullseye) and later
- Fedora 35 and later
- RHEL/CentOS 8 and later
- Arch Linux (latest)

**Features:**

- WireGuard/AmneziaWG VPN
- meshcli command-line interface
- ADB integration
- Route advertisement
- Exit node support
- Subnet routing

**Requirements:**

- Go 1.21 or later (for building from source)
- Root/sudo access
- ADB (Android Debug Bridge)

**Installation:**

- Build from source (see [Analyst client guide](../installation/analyst-client.md))
- Pre-built binaries coming soon

### macOS (analyst client)

**Supported Versions:**

- macOS 11 (Big Sur) and later
- Tested on macOS 11, 12 (Monterey), 13 (Ventura), 14 (Sonoma)

**Features:**

- WireGuard/AmneziaWG VPN
- meshcli command-line interface
- ADB integration
- Route advertisement
- Exit node support
- Subnet routing

**Requirements:**

- Go 1.21 or later (for building from source)
- Xcode Command Line Tools
- ADB (Android Platform Tools)

**Installation:**

- Build from source (see [Analyst client guide](../installation/analyst-client.md))
- Pre-built binaries coming soon

### Control plane (self-Hosted Server)

**Supported Platforms:**

- Any Linux distribution with Docker support
- Tested on Ubuntu 20.04+, Debian 11+, RHEL 8+

**Features:**

- Headscale coordination server
- Web-based management UI
- API for automation
- SQLite or PostgreSQL database
- DERP relay server
- ACL policy enforcement
- Multi-user support

**Requirements:**

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later
- 1GB RAM minimum (2GB recommended)
- 10GB disk space
- Public IP or domain name (recommended)

**Installation:**

- Docker Compose deployment (see [Control plane guide](../installation/control-plane.md))

## Planned Support (2026)

The following platforms are planned for future releases:

| Platform | Expected | Role | Status |
|----------|----------|------|--------|
| **iOS** | Q4 2026 | Endpoint bridge client | In development |
| **Windows** | Q4 2026 | Analyst client | Planned |

### iOS (endpoint bridge client) - Q4 2026

**Planned Features:**

- WireGuard/AmneziaWG VPN
- libimobiledevice integration
- Background operation
- Connection status UI
- iOS diagnostic interface access

**Requirements (Estimated):**

- iOS 14.0 or later
- WiFi or cellular connectivity
- VPN permission

**Status:**

- Currently in early development
- Prototype testing in progress
- Expected release: Q4 2026

!!! info "iOS Development"
    iOS support is actively being developed. Follow progress on [GitHub](https://github.com/BARGHEST-ngo/mesh) or join the [Discord community](https://discord.com/invite/) for updates.

### Windows (analyst client) - Q4 2026

**Planned Features:**

- WireGuard/AmneziaWG VPN
- meshcli command-line interface
- ADB integration
- Route advertisement
- Exit node support

**Requirements (Estimated):**

- Windows 10 version 1809 or later
- Windows 11 (all versions)
- Administrator access

**Status:**

- Planned for Q4 2026
- Development not yet started
- Community contributions welcome

## Community contributions

MESH is open source and welcomes community contributions for platform support:

See the [Contributing Guide](https://github.com/BARGHEST-ngo/mesh/blob/main/CONTRIBUTING.md) for more information.

## Compatibility Notes

### Android Compatibility

**Known Issues:**

- Some Android OEMs (Xiaomi, Huawei) may require additional battery optimization exemptions
- ADB-over-WiFi may be disabled by some custom ROMs
- VPN permission may conflict with other VPN apps

**Workarounds:**

- Disable battery optimization for MESH app
- Manually enable ADB-over-WiFi in Developer Options
- Disconnect other VPN apps before using MESH

### Linux Compatibility

**Known Issues:**

- Some distributions require manual kernel module loading for WireGuard
- SELinux may block VPN operations on RHEL/CentOS
- AppArmor profiles may need adjustment on Ubuntu

**Workarounds:**

- Install WireGuard kernel module: `sudo modprobe wireguard`
- Configure SELinux policies or set to permissive mode
- Adjust AppArmor profiles as needed

### macOS Compatibility

**Known Issues:**

- macOS Firewall may prompt for network access
- System Integrity Protection (SIP) may interfere with VPN operations
- Gatekeeper may block unsigned binaries

**Workarounds:**

- Allow network access when prompted
- SIP should not need to be disabled for normal operation
- Build from source or wait for signed binaries

---

**Next:** Learn [About BARGHEST](about.md) and the team behind MESH →
