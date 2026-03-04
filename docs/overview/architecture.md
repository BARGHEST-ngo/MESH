# Architecture Overview

MESH consists of three main components that work together to create a secure, censorship-resistant forensic network.

## System Components

```
┌─────────────────────────────────────────────────────────────┐
│                      MESH network                           │
│                                                             │
│  ┌──────────────┐         ┌──────────────┐                  │
│  │   Analyst    │◄───────►│   Endpoint   │                  │
│  │   Client     │  P2P    │   Client     │                  │
│  │ (Linux/macOS)│ Tunnel  │  (Android)   │                  │
│  └──────────────┘         └──────────────┘                  │
│         |                        │                          │
│              ┌──────────────┐      Key exchange and ACL     │
│         └ ──►│   Control    │◄─ ─┘ HTTPS relay fallback     │
│              │    Plane     │                               │
│              │ (Headscale)  │                               │
│              └──────────────┘                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 1. Control plane

The **Control plane** is a self-hosted coordination server that manages your mesh network.

### What It Does

- **Manages mesh network membership** - Tracks which nodes are part of the mesh
- **Distributes WireGuard keys automatically** - Handles all cryptographic key exchange
- **Enforces Access Control Lists (ACLs)** - Controls which nodes can communicate
- **Coordinates NAT traversal** - Helps nodes establish direct P2P connections
- **Operates DERP relays** - Provides fallback relay when P2P fails

### Technology

- Based on **Headscale** (open-source Tailscale Protocol control server)
- Written in **Go** for performance and reliability
- Uses **SQLite** or **PostgreSQL** for data storage
- Deployed via **Docker** for easy management
- Includes **Web UI** for graphical management

### Key features

- **Self-hostable** - You control your own infrastructure
- **API-driven** - Programmatic access for automation
- **Multi-user support** - Organise nodes by user or team
- **Pre-authentication keys** - Automated node enrollment
- **Route management** - Control subnet routing and exit nodes

!!! info "Control plane setup"
    See the [Control plane documentation](../installation/control-plane.md) for detailed setup and configuration.

## 2. Analyst client

The **Analyst client** runs on your forensic collection node (Linux or macOS) and provides the tools to manage connections and perform investigations.

### What It Does

- **Connects to target devices** over the encrypted mesh
- **Runs forensic tools** like ADB, AndroidQF, MVT, and libimobiledevice
- **Can flip to an exist node allowing network monitoring/capture** from remote devices
- **Manages mesh network configuration** via CLI
- **Provides secure access** to endpoint diagnostic interfaces

### Technology

- Based on **Tailscale Daemon** with AmneziaWG modifications
- Written in **Go** with platform-specific integrations
- Uses **WireGuard/AmneziaWG** for encrypted tunnels
- Includes **meshcli** command-line interface
- Runs **tailscaled-amnezia** daemon in background

### Key features

- **ADB integration** - Direct access to Android devices
- **Automated AndroidQF collection** - Connect and collect instantly
- **Automatic reconnection** - Maintains persistent connections
- **Route advertisement** - Share local network access
- **Exit node support** - Route traffic through the mesh

### Supported platforms

- **Linux** - Ubuntu 20.04+, Debian 11+, Fedora 35+, and similar
- **macOS** - macOS 11 (Big Sur) and later

!!! info "Analyst client Setup"
    See the [Analyst client documentation](../installation/analyst-client.md) for installation and usage.

## 3. Endpoint client

The **Endpoint client** runs on the target device being analysed (right now we only support Android).

### What it does

- **Joins the forensic mesh network** using pre-authentication keys
- **Guides user to enable ADB** to be connected to over the mesh
- **Enables remote forensic collection** without physical access
- **Supports ephemeral connections** that leave no trace
- **Provides secure access** to device internals
- **Can be deployed in seconds**

### Technology (Android)

- Written in **Kotlin/Java** for Android platform
- Uses **WireGuard/AmneziaWG** for encrypted tunnels
- Integrates with **Android VPN API**
- Guides user to enable **ADB-over-WiFi** on mesh interface

!!! experimental "AmneziaWG is currently in development on the APK"

### Key features

- **VPN-based networking** - Uses Android VPN API for mesh interface
- **Background operation** - Maintains connection when app is closed
- **Battery efficient** - Minimal impact on device battery life
- **User-friendly UI** - Simple connection and status display

### Supported Platforms

- **Android** - Android 8.0 (Oreo) and later
- **iOS** - Planned for Q4 2026

!!! info "Endpoint client Setup"
    See the [Endpoint client documentation](../installation/endpoint-clients.md) for installation and configuration.

## How they work together

### 1. Initial setup

1. **Deploy Control plane** - Set up the coordination server
2. **Create Users** - Organise nodes by user or team
3. **Generate Pre-Auth Keys** - Create keys for automated enrollment

### 2. Node enrollment

1. **Analyst client** connects to control plane using pre-auth key
2. **Endpoint client** connects to control plane using pre-auth key
3. **Control plane** authenticates nodes and assigns mesh IP addresses
4. **WireGuard keys** are automatically exchanged between nodes

### 3. Connection establishment

1. **NAT traversal** - Nodes attempt direct P2P connection via STUN
2. **Firewall traversal** - Nodes negotiate through NAT and firewalls
3. **DERP fallback** - If P2P fails, traffic routes through HTTPS DERP relay
4. **Encrypted tunnel** - Secure WireGuard tunnel established

### 4. Forensic operations

1. **Analyst** connects to endpoint via mesh IP address
2. **ADB/diagnostic tools** access device over encrypted tunnel
3. **Artifacts collected** and transferred securely
4. **Investigation completed** - Mesh can be torn down

## Network architecture

### Peer-to-Peer First

MESH prioritizes direct peer-to-peer connections for:

- **Lower latency** - Direct connections are faster
- **Better performance** - No relay overhead
- **Reduced bandwidth** - No unnecessary relay traffic
- **Enhanced privacy** - Traffic doesn't pass through third parties

### DERP Relay Fallback

When P2P connections fail (restrictive NAT, blocked UDP), MESH automatically falls back to DERP relays:

- **Encrypted relay** - End-to-end encryption maintained
- **TCP/443 transport** - Works through firewalls
- **Automatic failover** - Seamless transition from P2P to relay
- **Multiple relay servers** - Distributed infrastructure for reliability
- **Self hostable relays** - You can host your own relays or use your control plane as one

### NAT Traversal

MESH uses multiple techniques to establish connections through NAT:

- **STUN** - Discovers public IP and port mappings
- **Hole-punching** - Establishes direct connections through NAT
- **UPnP/NAT-PMP** - Automatic port forwarding when available
- **DERP relay** - Fallback when direct connection impossible

## Security model

### Zero trust architecture

- **No implicit trust** - All connections must be authenticated
- **Mutual authentication** - Both peers verify each other
- **Encrypted by default** - All traffic is encrypted end-to-end
- **Least privilege** - ACLs enforce minimal necessary access

### Encryption

- **WireGuard protocol** - Industry-standard VPN encryption
- **ChaCha20-Poly1305** - Authenticated encryption
- **Curve25519** - Elliptic curve key exchange
- **BLAKE2s** - Cryptographic hashing

### Access control

- **ACL policies** - Define which nodes can communicate
- **User-based access** - Organise nodes by user or team
- **Tag-based policies** - Flexible access control rules
- **Port-level restrictions** - Control access to specific services

!!! info "Detailed architecture"
    For a deep dive into MESH architecture, see the [Architecture documentation](architecture.md).

## Deployment models

### Single investigator

- One control plane
- One analyst nodes
- Multiple endpoint devices
- Simple ACL policy (allow only ADB traffic)

### Team investigation

- One control plane
- Multiple analyst nodes
- Multiple endpoint devices
- User-based ACL policies that segregate networks

### Multi-site deployment

- Multiple control planes (regional)
- Distributed analyst teams
- Geographically dispersed endpoints
- Complex ACL policies with tags

---

**Next:** Learn about real-world [Use cases](use-cases.md) →
