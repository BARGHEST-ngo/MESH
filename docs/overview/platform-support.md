# Platform support

MESH is designed to work across multiple platforms to support diverse forensic investigation scenarios. This page details current support status and planned features.

## Current support (public alpha)

The following platforms are currently supported in the public alpha release:

| Component | Android | iOS | Linux | macOS | Windows |
|-----------|---------|-----|-------|-------|---------|
| **Control plane** | ❌ | ❌ | ✅ | ✅ | ✅ |
| **Analyst client** | ❌ | ❌ | ✅ | ✅ | ✅ |
| **Endpoint client** | ✅ Android 8.0+ | ❌ Coming Q4 2026 | ❌ | ❌ | ❌ |

### Android (endpoint client)

**Supported Versions:**

- Android 8.0 (Oreo) and later
- Tested on Android 8.0 through Android 16

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

- Build from source (see [Endpoint client guide](../setup/endpoint-clients.md))
- Pre-built APKs available on [GitHub Releases](https://github.com/BARGHEST-ngo/mesh/releases)

### Analyst client

**Supported Platforms:**

- Any operating system with Docker support
- Tested on Ubuntu 20.04+, Debian 11+, RHEL 8+, macOS 11+, Windows 10+

**Requirements:**

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later

**Installation:**

- Docker Compose deployment (see [Analyst client guide](../setup/analyst-client.md))
- Pre-built container images coming soon

### Control plane

**Supported Platforms:**

- Any operating system with Docker support
- Tested on Ubuntu 20.04+, Debian 11+, RHEL 8+, macOS 11+, Windows 10+

**Requirements:**

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later
- 1GB RAM minimum (2GB recommended)
- 10GB disk space
- Public IP or domain name (persistent deployment) or tunneling service (ephemeral deployment)

**Installation:**

- Docker Compose deployment (see [Control plane guide](../setup/control-plane.md))
- Pre-built container images coming soon

## Planned Support (2026)

The following platforms are planned for future releases:

| Platform | Expected | Role | Status |
|----------|----------|------|--------|
| **iOS** | Q4 2026 | Endpoint client | In development |

### iOS (endpoint client) - Q4 2026

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
    iOS support is actively being developed. Follow progress on [GitHub](https://github.com/BARGHEST-ngo/mesh) for updates.

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
