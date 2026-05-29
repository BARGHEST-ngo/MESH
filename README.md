<h1><div align="center">
  <img width="269" height="188" alt="MESH" src="https://github.com/user-attachments/assets/03d0b879-ca5e-4a20-af72-533ea3ae2b52" />
</div></h1>

<div align="center">
  <p>
    <a href="https://docs.meshforensics.org/"><img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" /></a>
    <a href="#"><img src="https://img.shields.io/badge/status-public%20alpha-orange?style=flat-square" alt="Public Alpha Status" /></a>
    <a href="https://deepwiki.com/BARGHEST-ngo/MESH"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/BARGHEST-ngo/MESH?style=flat-square" alt="License" /></a>
    <a href="https://github.com/BARGHEST-ngo/MESH/releases"><img src="https://img.shields.io/github/v/release/BARGHEST-ngo/MESH?include_prereleases&style=flat-square" alt="Latest Release" /></a>
  </p>
  <p>
    <a href="https://play.google.com/store/apps/details?id=com.barghest.mesh"><img height="60" alt="Get it on Google Play" src="https://play.google.com/intl/en_us/badges/static/images/badges/en_badge_web_generic.png" /></a>
    <a href="https://f-droid.org/en/packages/com.barghest.mesh/"><img height="60" alt="Get it on F-Droid" src="https://fdroid.gitlab.io/artwork/badge/get-it-on.png" /></a>
  </p>
</div>

<h3><div align="center">
  <h3>
    <a href="https://docs.meshforensics.org/">
      DOCS
    </a>
    <span>  |  </span>
    <a href="https://Barghest.asia">
      BARGHEST
    </a>
  </h3>
</div>
<br/></h3>

> [!IMPORTANT]
> **Public Alpha**: Currently in **public alpha** and under active development. A full penetration has been completed and we have patched all major vulnerabilities. Things may change and breaking changes should be expected. It currently requires some level of technical expertise. Please report bugs or security concerns via GitHub Issues."

**MESH enables internet-routable wireless debugging for mobile devices over an encrypted, censorship-resistant peer-to-peer mesh network.**

Mobile devices are often trapped behind NAT, firewalls, carrier-grade NAT, or restrictive mobile networks that prevent direct inbound access. Traditional remote forensics usually depends on centralized VPN servers, brittle hub-and-spoke setups, or risky port forwarding.

MESH solves this by creating an encrypted peer-to-peer overlay network and assigning each node a CGNAT-range address through a virtual TUN interface. Devices appear as if they are on the same private subnet, even when they are geographically distant or hidden behind multiple NAT layers.

This allows analysts to run remote mobile forensics workflows over ADB Wireless Debugging and libimobiledevice without exposing target devices to the public internet. Tools such as WARD, MVT, AndroidQF, and other ADB/iOS tooling can operate across the encrypted overlay as if the device were locally reachable.

MESH also supports remote network monitoring, including PCAP capture and Suricata-based intrusion detection over the same encrypted mesh. This makes it possible to perform both immediate forensic acquisition and live network capture from devices in hostile or restricted environments.

MESH is designed for civil society forensics and hardened for censored or adversarial networks:

- Direct peer-to-peer WireGuard transport when available
- Optional AmneziaWG support to obfuscate WireGuard fingerprints against DPI and national firewall blocking
- Automatic fallback to end-to-end encrypted HTTPS relays when UDP is blocked

Meshes are ephemeral and analyst-controlled. Bring devices online, collect evidence, monitor traffic, and tear the network down immediately afterward — without maintaining permanent infrastructure or complex VPN configurations.

## Quick start

For full documentation:  
https://docs.meshforensics.org/

### 1. Clone the repository

```
git clone https://github.com/BARGHEST-ngo/mesh.git
cd MESH
```

### 2. Start control plane and get an API key

```
task build
task controlPlane
task apikey
```

### 3. Access web UI with API key

```
Local:  https://localhost
Remote: https://your-domain:8443/login
```

The Web UI uses a self-signed certificate by default.

> [!IMPORTANT]
> The default ACL allows nodes in each network talk to each other.
> Production deployments should use restrictive policies. Modify these via the ACL tab.

<img width="1035" height="868" alt="image" src="https://github.com/user-attachments/assets/52bb4020-18fe-4b49-9b33-516250055278" />

Your MESH network is now ready to accept nodes.  
See the documentation for node enrollment and forensic workflows.

## Architecture summary

MESH is a heavily modified fork of the [Tailscale protocol](https://github.com/tailscale/tailscale), but does not require Tailscale infrastructure.

To establish peer-to-peer, end-to-end encrypted channel is created using UDP hole punching. If UDP is unavailable or blocked, it will fail over to E2EE HTTPs relays called DERP relays. The DERP protocol [DERP (Designated Encrypted Relay for Packets)](https://github.com/tailscale/tailscale/tree/main/derp) servers relay traffic between nodes when a direct peer-to-peer connection cannot be established.

MESH follows the same model. By default, if an operator has not configured their own DERP infrastructure (which can be done using MESH's control plane), MESH uses Tailscale’s public DERP servers to ensure reliable connectivity, particularly in restrictive network environments. However, MESH does not require Tailscale infrastructure: operators can deploy and use their own DERP servers via the control plane, which includes an embedded DERP implementation. This makes MESH fully self-hostable when desired.

DERP servers act purely as transport relays. They facilitate connectivity between devices but do not have visibility into the data exchanged, which remains end-to-end encrypted.

Enhancements include (but are not limited to):

- Self-hostable coordination server with a UI tailored for forensic operations
- Automatic WireGuard key distribution
- Optional AmneziaWG-based transport obfuscation
- Encrypted HTTPS relay fallback

The control plane is responsible only for peer discovery and key exchange.
Forensic traffic flows directly between endpoints whenever possible.

## Key capabilities

- Peer-to-peer encrypted forensic subnets
- Automatic WireGuard / AmneziaWG key management
- Self-hostable control plane with ACL enforcement
- CGNAT-assigned virtual TUN interfaces
- ADB-over-WiFi & libimobiledevice compatibility
- AndroidQF + MVT integration
- Secure transfer of forensic artifacts
- Optional kill-switch containment
- Rapid mesh creation and teardown

<img width="1920" height="1080" alt="1" src="https://github.com/user-attachments/assets/2e539f33-d46a-4396-b25e-43c23e9e4040" />

<img width="1920" height="1080" alt="2" src="https://github.com/user-attachments/assets/31947724-d826-4807-b984-9b87473e6847" />

<img width="1920" height="1080" alt="3" src="https://github.com/user-attachments/assets/1d95e6ff-2361-42b8-af8a-9d01c18c6b3b" />

## Why not a VPN?

Traditional VPN and hub-and-spoke architectures introduce:

- Persistent infrastructure risk
- Centralized traffic analysis points
- Single points of failure
- Increased operational exposure

MESH separates coordination from data transport:

- The control plane does not carry forensic traffic
- Peer connections are direct whenever possible
- Relays are transport fallbacks, not architectural hubs
- Meshes are disposable and task-scoped

MESH is optimized for transient, high-risk environments rather than permanent enterprise networking.

### Repository structure

- `android-client` — Android endpoint APK
- `control-plane` — Coordination server
- `analyst` — Analyst CLI client

### Developer notes

**Workflow:**

- Development happens on branches and is merged via PRs.
- Releases are cut as versioned tags.
- GitHub Actions mirrors tagged releases to `mesh-analyst-client`.
- External Go projects should depend on explicit version tags, not `main`.

## License

MESH is licensed under the [GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later)](LICENSE).

Portions of this software are a derivative work of [Tailscale](https://github.com/tailscale/tailscale), which is licensed under the BSD 3-Clause License. The original Tailscale copyright and license are preserved in accordance with the BSD-3-Clause requirements. AmneziaWG/Wireguard code is licensed under MIT license. See `.licenses/` for details.

All modifications and additions by BARGHEST are Copyright (c) BARGHEST and licensed under AGPL-3.0-or-later.

### Legal

WireGuard is a registered trademark of Jason A. Donenfeld.
