# MESH Architecture

This document provides a comprehensive overview of MESH's architecture, components, and how they work together to enable secure, censorship-resistant forensic investigations.

## System Overview

MESH is built on a modified Tailscale protocol stack, replacing the proprietary Tailscale control plane with the open-source Headscale server. The system creates a peer-to-peer overlay network using WireGuard (or AmneziaWG for censorship resistance) with automatic NAT traversal and encrypted relay fallback.

```
┌─────────────────────────────────────────────────────────────┐
│                     MESH Network                             │
│                                                              │
│  ┌──────────────┐         ┌──────────────┐                 │
│  │   Analyst    │◄───────►│   Endpoint   │                 │
│  │   Client     │  P2P    │   (Android)  │                 │
│  │ (Linux/macOS)│ Tunnel  │              │                 │
│  └──────┬───────┘         └──────┬───────┘                 │
│         │                        │                          │
│         │  ┌──────────────────┐  │                          │
│         └─►│  Control plane   │◄─┘                          │
│            │   (Headscale)    │                             │
│            └────────┬─────────┘                             │
│                     │                                        │
│            ┌────────▼─────────┐                             │
│            │   DERP Relays    │                             │
│            │ (Fallback Path)  │                             │
│            └──────────────────┘                             │
└─────────────────────────────────────────────────────────────┘
```

## Core components

### 1. Control plane

The control plane is the coordination server that manages the MESH network. It's based on Headscale, an open-source implementation of the Tailscale control server.

**Responsibilities:**

- **Node Registration**: Authenticates and registers new nodes joining the mesh
- **Key Distribution**: Distributes WireGuard public keys between peers
- **ACL Enforcement**: Enforces access control policies
- **Network Coordination**: Coordinates NAT traversal and peer discovery
- **DERP Coordination**: Manages fallback relay connections

**Technology Stack:**

- Go-based HTTP/gRPC server
- SQLite or PostgreSQL database
- Docker containerized deployment
- REST API for management
- Web-based UI for graphical management

**Key features:**

- **Web UI**: Intuitive browser-based interface for managing the mesh
- Pre-authentication keys for automated node enrollment
- User and machine management
- Route advertisement and subnet routing
- Exit node configuration
- API key-based authentication
- ACL policy editor

### 2. Analyst client (Linux/macOS)

The analyst client is the forensic acquisition node used by investigators to connect to and analyse target devices.

**Responsibilities:**

- Connect to the MESH network
- Establish peer-to-peer tunnels to endpoint devices
- Run forensic tools (ADB, AndroidQF, MVT, libimobiledevice)
- Collect and analyse artifacts
- Manage network configuration

**Technology Stack:**

- Go-based CLI
- WireGuard/AmneziaWG for encryption
- TUN interface for virtual networking
- Integration with forensic tools
- Tailscale Daemon

**Key features:**

- Command-line interface for network management
- ADB-over-WiFi support for Android forensics
- libimobiledevice support for iOS forensics (planned)
- Route advertisement for LAN access
- Exit node capabilities

### 3. Endpoint client (Android/iOS)

The endpoint client runs on the target device being analysed. Currently supports Android, with iOS support planned.

**Responsibilities:**

- Join the MESH network
- Establish encrypted tunnel to analyst
- Guide user in ADB interfaces
- Split tunnel - remove app traffic from the mesh network
- Support ephemeral connections

**Technology Stack:**

- Android: Native app with WireGuard/AmneziaWG
- iOS: Planned native app (Q4 2026)
- VPN service integration
- Background service for persistent connection

**Key features:**

- One-tap connection to mesh
- ADB-over-WiFi guiudance
- Ephemeral mode (disconnect on app close)
- Battery optimization
- Network usage monitoring
- Setting of exit node for network monitoring

## Network architecture

### Peer-to-Peer connectivity

MESH prioritizes direct peer-to-peer connections between nodes using WireGuard tunnels. This provides:

- **Low latency**: Direct connections without intermediary hops
- **High throughput**: No bandwidth bottlenecks from relay servers
- **Privacy**: Traffic doesn't traverse third-party servers

### NAT traversal

MESH uses multiple techniques to establish direct connections through NAT:

1. **STUN**: Discovers public IP and port mappings
2. **UDP Hole Punching**: Establishes bidirectional UDP flows
3. **Port Prediction**: Predicts NAT port allocation patterns
4. **Hairpin Detection**: Detects and handles hairpin NAT scenarios

### DERP relay fallback

When direct P2P connections fail (blocked UDP, symmetric NAT, restrictive firewalls), MESH falls back to DERP (Detoured Encrypted Relay for Packets):

- **Encrypted Relay**: All traffic remains end-to-end encrypted
- **HTTPS Transport**: Relays use HTTPS to bypass firewall restrictions
- **Automatic Failover**: Seamless transition from P2P to relay
- **Global Distribution**: Multiple relay servers for redundancy

**DERP Protocol:**

```
Analyst ──[E2E Encrypted]──► DERP Relay ──[E2E Encrypted]──► Endpoint
         (HTTPS Tunnel)                    (HTTPS Tunnel)
```

The DERP relay only sees encrypted WireGuard packets and cannot decrypt traffic.

## Addressing and routing

### CGNAT address space

MESH uses CGNAT (Carrier-Grade NAT) address ranges for the overlay network:

- **IPv4**: 100.64.0.0/10 (RFC 6598)
- **IPv6**: fd7a:115c:a1e0::/48 (Unique Local Address)

These addresses are routable within the mesh but not on the public internet.

### Subnet routing

Nodes can advertise routes to their local networks:

```
Analyst (100.64.1.1) ──► Endpoint (100.64.2.1)
                              │
                              └──► LAN (192.168.1.0/24)
```

This enables access to devices on the endpoint's local network.

### Exit nodes

Nodes can act as exit nodes, routing internet traffic through the mesh:

```
Analyst ──► Exit Node ──► Internet
         (Encrypted)   (Cleartext)
```

Useful for monitoring endpoint's internet traffic or providing internet access.

## Deployment topologies

### Single analyst, single endpoint

Simplest deployment for one-on-one forensic analysis:

```
Analyst ◄──────► Control plane ◄──────► Endpoint
```

### Multiple analysts, multiple endpoints

Team-based forensic operations:

```
Analyst 1 ◄─┐
Analyst 2 ◄─┼──► Control plane ◄─┬──► Endpoint 1
Analyst 3 ◄─┘                     ├──► Endpoint 2
                                  └──► Endpoint 3
```

ACLs control which analysts can access which endpoints.

### Distributed DERP relays

For global operations with censorship resistance:

```
                    ┌──► DERP US
Analyst ◄──► Control plane ──► DERP EU ◄──► Endpoint
                    └──► DERP Asia
```

Multiple DERP servers provide redundancy and geographic distribution.
