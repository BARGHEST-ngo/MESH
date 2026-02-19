<h1><div align="center">
  <img width="269" height="188" alt="MESH" src="https://github.com/user-attachments/assets/03d0b879-ca5e-4a20-af72-533ea3ae2b52" />
</div></h1>

<div align="center">
  <p>
    <a href="https://docs.meshforensics.org/">
      <img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" />
    </a>
    <a href="#">
      <img src="https://img.shields.io/badge/status-public%20alpha-orange?style=flat-square" alt="Public Alpha Status" />
    </a>
    <a href="https://discord.com/invite/">
      <img src="https://img.shields.io/discord/1161119546170687619?logo=discord&style=flat-square" alt="Chat" />
    </a>
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
> **Public Alpha**: MESH is currently in **public alpha** and under active development. It is **not production-ready**. A full penetration test is in progress. Until it is complete, do **not** use this project in production environments. Things may change and breaking changes should be expected. It currently requires some level of technical expertise. Please report bugs or security concerns via GitHub Issues.

**MESH enables remote mobile forensics by assigning CGNAT-range IP addresses to devices over an encrypted, censorship-resistant peer-to-peer mesh network.**

Mobile devices are often placed behind carrier-grade NAT (CGNAT), firewalls, or restrictive mobile networks that prevent direct inbound access. Traditional remote forensics typically requires centralized VPN servers or risky port-forwarding.

MESH solves this by creating an encrypted peer-to-peer overlay and assigning each node a CGNAT-range address via a virtual TUN interface. Devices appear as if they are on the same local subnet — even when geographically distant or behind multiple NAT layers.

This enables **remote mobile forensics** using ADB Wireless Debugging and [libimobiledevice](https://libimobiledevice.org/), allowing tools such as WARD, [MVT](https://github.com/mvt-project/), and [AndroidQF](https://github.com/mvt-project/androidqf) to operate remotely without exposing devices to the public internet.  

The mesh can also be used for **remote network monitoring**, including PCAP capture and Suricata-based intrusion detection over the encrypted overlay. Allowing for both immediate forensics capture and network capture.

MESH is designed specifically for civil society forensics & hardened for hostile/censored networks:

- Direct peer-to-peer WireGuard transport when available  
- Optional AmneziaWG to obfuscate WireGuard fingerprints to evade national firewalls or DPI inspection
- Automatic fallback to end-to-end encrypted HTTPS relays when UDP is blocked  

Meshes are ephemeral and analyst-controlled: bring devices online, collect evidence, and tear the network down immediately afterward. No complicated hub-and-spoke configurations.

## Architecture summary

MESH is a heavily modified fork of the [Tailscale protocol](https://github.com/tailscale/tailscale) but does **not** require Tailscale infrastructure.

Enhancements include, though not limited to:

- Self-hostable coordination server with UI for forensics operations
- Automatic WireGuard key distribution  
- Optional AmneziaWG transport obfuscation  
- Encrypted HTTPS relay fallback  

The control plane handles peer discovery and key exchange only.  
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

<img width="1920" height="1080" alt="1" src="https://github.com/user-attachments/assets/4b305c4c-c0f9-49bf-8e59-9ccecdfadb86" />

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

## Getting started

For full documentation:  
https://docs.meshforensics.org/

### 1. Clone the repository

```
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh/control-plane
```

### 2. Start control plane

```
docker-compose up -d
```

### 4. Access web UI

```
Local:  https://localhost:3000/login
Remote: https://your-domain:8443/login
```

The Web UI uses a self-signed certificate by default.

### 5. Create API key

```
docker exec headscale headscale apikeys create --expiration 90d
```

Use the generated key to authenticate in the Web UI.

> [!IMPORTANT]
> The default ACL allows nodes in each network talk to each other.
> Production deployments should use restrictive policies. Modify these via the ACL tab.

<img width="1035" height="868" alt="image" src="https://github.com/user-attachments/assets/52bb4020-18fe-4b49-9b33-516250055278" />

Your MESH network is now ready to accept nodes.  
See the documentation for node enrollment and forensic workflows.

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

### License

Forked components from Tailscale are licensed under BSD-3-Clause.  
Additional code is licensed under CC0.  
See `LICENSE` and `.licenses/` for details.

### Legal

WireGuard is a registered trademark of Jason A. Donenfeld.
