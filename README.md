<div align="center">
  <img width="269" height="188" alt="MESH" src="https://github.com/user-attachments/assets/03d0b879-ca5e-4a20-af72-533ea3ae2b52" />

  <p>
    <a href="https://docs.meshforensics.org/"><img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" /></a>
    <a href="#"><img src="https://img.shields.io/badge/status-public%20beta-blueviolet?style=flat-square" alt="Public Beta Status" /></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/BARGHEST-ngo/MESH?style=flat-square" alt="License" /></a>
    <a href="https://github.com/BARGHEST-ngo/MESH/releases"><img src="https://img.shields.io/github/v/release/BARGHEST-ngo/MESH?include_prereleases&style=flat-square" alt="Latest Release" /></a>
  </p>
  <p><a href="https://docs.meshforensics.org/">DOCS</a> &nbsp;|&nbsp; <a href="https://Barghest.asia">BARGHEST</a></p>
</div>

**MESHforensics enables internet-routable wireless debugging for mobile devices over an encrypted, censorship-resistant peer-to-peer mesh network.**

Mobile devices are often trapped behind NAT, firewalls, or restrictive carrier networks that block direct inbound access. MESHforensics creates an encrypted peer-to-peer overlay and assigns each node a CGNAT-range address via a virtual TUN interface, so devices behave as if they share a private subnet, with no central VPN, port forwarding, or permanent infrastructure.

This lets analysts run remote mobile forensics over ADB Wireless Debugging and libimobiledevice, using tools like [WARD](https://github.com/BARGHEST-ngo/WARD), [MVT](https://github.com/mvt-project/mvt), and [AndroidQF](https://github.com/mvt-project/androidqf) as if the device were local. MESHforensics also supports remote network monitoring, including PCAP capture and Suricata intrusion detection over the same mesh.

Meshes are ephemeral and analyst-controlled: bring devices online, collect, then tear the network down, so nothing is left running to fingerprint.

> [!IMPORTANT]
> **Public Beta**: under active development. A full penetration test has been completed and major vulnerabilities patched. Breaking changes should be expected, and some technical expertise is required. Report bugs or security concerns via GitHub Issues.

## Getting started

Full documentation: https://docs.meshforensics.org/

Clone the repo first:

```
git clone https://github.com/BARGHEST-ngo/mesh.git
cd MESH
```

Then pick how your control plane is reached from the internet.

### Option A: Provisioner service (hosted reverse proxy)

Run the control plane on your workstation and expose it through the MESHforensics provisioner, a hosted reverse-proxy tunnel. Request a provisioner API key from BARGHEST (see https://docs.meshforensics.org/). The provisioner auto-assigns an ephemeral `*.tunnels.meshforensics.app` subdomain and issues your API credentials, so **`task apikey` is not needed** on this path.

```
# Choose the "Ephemeral" control plane when prompted
task build
task controlPlane

# Open the web UI at the provisioner-assigned URL and log in
```

### Option B: Self-hosted reverse proxy

Bring your own ingress (a reverse proxy on an internet-facing server, or a tunneling service such as [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/) or [ngrok](https://ngrok.com/)), then mint an API key locally.

```
# Choose the "Persistent" control plane when prompted, and your own domain
# (or the hostname your Cloudflare Tunnel / ngrok endpoint provides)
task build
task controlPlane
task apikey

# Open the web UI and log in with the API key
#    Local:  https://localhost
#    Remote: https://your-domain:8443/login
```

The web UI uses a self-signed certificate by default. Once logged in, your MESHforensics network is ready to accept nodes. See the docs for node enrollment and forensic workflows.

> [!IMPORTANT]
> The default ACL lets all nodes in a network talk to each other. Production deployments should use restrictive policies, set via the ACL tab.

## Architecture

MESHforensics is a heavily modified fork of [Tailscale](https://github.com/tailscale/tailscale) that does not require Tailscale infrastructure. Peers connect directly via end-to-end encrypted WireGuard using UDP hole punching. When UDP is blocked, traffic fails over to E2EE HTTPS [DERP](https://github.com/tailscale/tailscale/tree/main/derp) relays, which transport packets without visibility into their contents.

The control plane handles only peer discovery and key exchange, and forensic traffic flows directly between endpoints whenever possible. Operators can run their own DERP servers (the control plane embeds an implementation), making MESHforensics fully self-hostable. Otherwise it falls back to Tailscale's public DERP servers for connectivity.

## Key capabilities

- Peer-to-peer encrypted forensic subnets with CGNAT-assigned TUN interfaces
- Automatic WireGuard / AmneziaWG key management, with optional transport obfuscation against DPI
- Self-hostable control plane with ACL enforcement
- ADB-over-WiFi and libimobiledevice compatibility, plus AndroidQF + MVT integration
- Secure forensic artifact transfer, optional kill-switch containment
- Rapid mesh creation and teardown

## Why not a VPN like WireGuard?

MESHforensics is built on WireGuard, but a plain WireGuard VPN is the wrong shape for hostile networks. WireGuard speaks only UDP, and UDP is exactly what censored and restrictive carrier networks degrade or block first. When that happens a raw WireGuard tunnel simply fails, which is unacceptable when an analyst needs to reach a device in the field.

MESHforensics keeps WireGuard for its fast, modern, end-to-end encrypted transport, then adds the connectivity WireGuard lacks on its own. Peers attempt a direct UDP connection via hole punching, and when UDP is blocked they fail over to end-to-end encrypted HTTPS DERP relays. HTTPS relaying is the critical piece: it is indistinguishable from ordinary web traffic, so it survives DPI and UDP censorship that would strand a conventional VPN, while the relay still carries packets without visibility into their contents.

Traditional VPN and hub-and-spoke designs also add persistent infrastructure risk, centralized traffic-analysis points, single points of failure, and operational exposure. MESHforensics separates coordination from transport. The control plane carries no forensic traffic, peer connections are direct, relays are fallbacks rather than hubs, and meshes are disposable and task-scoped. It is optimized for transient, high-risk environments rather than permanent enterprise networking.

## Repository structure

- `android-client`: Android endpoint APK
- `control-plane`: Coordination server
- `analyst`: Analyst CLI client

## Developer notes

- Development happens on branches, merged via PRs. Releases are cut as versioned tags.
- GitHub Actions mirrors tagged releases to `mesh-analyst-client`.
- External Go projects should depend on explicit version tags, not `main`.

## Supported by

<img width="312" height="109" alt="image" src="https://github.com/user-attachments/assets/d9476c43-7414-46e0-b6a5-1956bd6ca1b5" />

MESHforensics is supported by the [Open Technology Fund](https://www.opentech.fund/) to support the consensual remote forensics of spyware and surveillance targeting Human Right Defenders.

## License

MESHforensics is licensed under the [GNU Affero General Public License v3.0 or later](LICENSE) (AGPL-3.0-or-later). All modifications and additions by BARGHEST are Copyright (c) BARGHEST.

Portions are a derivative work of [Tailscale](https://github.com/tailscale/tailscale) (BSD 3-Clause). AmneziaWG/WireGuard code is MIT-licensed. Original copyrights and licenses are preserved in `.licenses/`. WireGuard is a registered trademark of Jason A. Donenfeld.
