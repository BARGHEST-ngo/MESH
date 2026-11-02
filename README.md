<div align="center">
  <img width="269" height="188" alt="MESH" src="https://github.com/user-attachments/assets/03d0b879-ca5e-4a20-af72-533ea3ae2b52" />

  <p>
    <a href="https://docs.meshforensics.org/"><img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" /></a>
    <a href="#"><img src="https://img.shields.io/badge/status-public%20alpha-orange?style=flat-square" alt="Public Alpha Status" /></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/BARGHEST-ngo/MESH?style=flat-square" alt="License" /></a>
    <a href="https://github.com/BARGHEST-ngo/MESH/releases"><img src="https://img.shields.io/github/v/release/BARGHEST-ngo/MESH?include_prereleases&style=flat-square" alt="Latest Release" /></a>
  </p>
  <p><a href="https://docs.meshforensics.org/">DOCS</a> &nbsp;|&nbsp; <a href="https://Barghest.asia">BARGHEST</a></p>
</div>

**MESH enables internet-routable wireless debugging for mobile devices over an encrypted, censorship-resistant peer-to-peer mesh network.**

Mobile devices are often trapped behind NAT, firewalls, or restrictive carrier networks that block direct inbound access. MESH creates an encrypted peer-to-peer overlay and assigns each node a CGNAT-range address via a virtual TUN interface, so devices behave as if they share a private subnet — no central VPN, port forwarding, or permanent infrastructure.

This lets analysts run remote mobile forensics over ADB Wireless Debugging and libimobiledevice, using tools like [WARD](https://github.com/BARGHEST-ngo/WARD), [MVT](https://github.com/mvt-project/mvt), and [AndroidQF](https://github.com/mvt-project/androidqf) as if the device were local. MESH also supports remote network monitoring, including PCAP capture and Suricata intrusion detection over the same mesh.

Meshes are ephemeral and analyst-controlled: bring devices online, collect, then tear the network down — nothing is left running to fingerprint.

> [!IMPORTANT]
> **Public Alpha** — under active development. A full penetration test has been completed and major vulnerabilities patched. Breaking changes should be expected, and some technical expertise is required. Report bugs or security concerns via GitHub Issues.

## Getting started

Full documentation: https://docs.meshforensics.org/

```
# 1. Clone
git clone https://github.com/BARGHEST-ngo/mesh.git
cd MESH

# 2. Build, start the control plane, and mint an API key
task build
task controlPlane
task apikey

# 3. Open the web UI and log in with the API key
#    Local:  https://localhost
#    Remote: https://your-domain:8443/login
```

The web UI uses a self-signed certificate by default. Once logged in, your MESH network is ready to accept nodes — see the docs for node enrollment and forensic workflows.

> [!IMPORTANT]
> The default ACL lets all nodes in a network talk to each other. Production deployments should use restrictive policies, set via the ACL tab.

## Architecture

MESH is a heavily modified fork of [Tailscale](https://github.com/tailscale/tailscale) that does not require Tailscale infrastructure. Peers connect directly via end-to-end encrypted WireGuard using UDP hole punching. When UDP is blocked, traffic fails over to E2EE HTTPS [DERP](https://github.com/tailscale/tailscale/tree/main/derp) relays, which transport packets without visibility into their contents.

The control plane handles only peer discovery and key exchange; forensic traffic flows directly between endpoints whenever possible. Operators can run their own DERP servers (the control plane embeds an implementation), making MESH fully self-hostable; otherwise it falls back to Tailscale's public DERP servers for connectivity.

## Key capabilities

- Peer-to-peer encrypted forensic subnets with CGNAT-assigned TUN interfaces
- Automatic WireGuard / AmneziaWG key management, with optional transport obfuscation against DPI
- Self-hostable control plane with ACL enforcement
- ADB-over-WiFi and libimobiledevice compatibility; AndroidQF + MVT integration
- Secure forensic artifact transfer, optional kill-switch containment
- Rapid mesh creation and teardown

## Why not a VPN?

Traditional VPN and hub-and-spoke designs add persistent infrastructure risk, centralized traffic-analysis points, single points of failure, and operational exposure. MESH separates coordination from transport: the control plane carries no forensic traffic, peer connections are direct, relays are fallbacks rather than hubs, and meshes are disposable and task-scoped. It is optimized for transient, high-risk environments rather than permanent enterprise networking.

## Repository structure

- `android-client` — Android endpoint APK
- `control-plane` — Coordination server
- `analyst` — Analyst CLI client

## Developer notes

- Development happens on branches, merged via PRs; releases are cut as versioned tags.
- GitHub Actions mirrors tagged releases to `mesh-analyst-client`.
- External Go projects should depend on explicit version tags, not `main`.

## Supported by

<img width="312" height="109" alt="image" src="https://github.com/user-attachments/assets/d9476c43-7414-46e0-b6a5-1956bd6ca1b5" />

MESH is supported by the [Open Technology Fund](https://www.opentech.fund/) to support the consensual remote forensics of spyware and surveillance targeting Human Right Defenders.

## License

MESH is licensed under the [GNU Affero General Public License v3.0 or later](LICENSE) (AGPL-3.0-or-later). All modifications and additions by BARGHEST are Copyright (c) BARGHEST.

Portions are a derivative work of [Tailscale](https://github.com/tailscale/tailscale) (BSD 3-Clause); AmneziaWG/WireGuard code is MIT-licensed. Original copyrights and licenses are preserved in `.licenses/`. WireGuard is a registered trademark of Jason A. Donenfeld.
