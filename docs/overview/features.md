# Features

MESH provides a comprehensive set of features designed specifically for secure mobile forensic investigations in challenging network environments.

## End-to-End Encryption (E2EE)

MESH ensures all communications are encrypted from endpoint to endpoint, protecting sensitive forensic data.

### Automatic key distribution

- **WireGuard/AmneziaWG keys** are automatically generated and distributed by the control plane
- No manual key exchange required
- No manual configuration changes to you VPN server
- Keys are rotated automatically for enhanced security
- Each peer-to-peer connection uses unique encryption keys

### Peer-to-Peer E2EE Tunnels

- Direct encrypted connections between analyst and endpoint devices
- Uses industry-standard cryptography:
  - **ChaCha20-Poly1305** for symmetric encryption
  - **Curve25519** for key exchange
  - **BLAKE2s** for hashing
- Perfect forward secrecy ensures past communications remain secure

### No plaintext exposure

- Even when using DERP relays, data remains encrypted end-to-end
- Control plane never sees any traffic
- Relay servers cannot decrypt mesh traffic
- All forensic data is protected in transit

!!! success "Security guarantee"
    MESH provides the same security guarantees as WireGuard: even if the control plane or relay servers are compromised, your forensic data remains encrypted and secure.

## Censorship resistance

MESH is designed to operate in hostile network environments with aggressive censorship and Deep packet inspection (DPI).

### AmneziaWG obfuscation

- **Evades Deep packet inspection (DPI)** by obfuscating WireGuard traffic
- Makes VPN traffic appear as regular HTTPS
- Effective against the Great Firewall of China and similar systems
- Configurable obfuscation parameters for different threat models

!!! info "This feature is still in development"

**How it works:**

AmneziaWG modifies WireGuard packet headers to avoid detection by DPI systems that specifically target VPN protocols. This allows MESH to operate in countries with aggressive internet censorship.

### HTTPS relay fallback

- Automatically falls back to HTTPS-wrapped relay when UDP is blocked
- DERP (Detoured Encrypted Relay for Packets) servers relay traffic over TCP/443
- Indistinguishable from regular HTTPS web traffic
- Works through GFWs and restrictive networks

!!! tip "Censorship resistance configuration"
    See the [AmneziaWG Configuration guide](../advanced/amneziawg.md) for detailed setup instructions.

## Mobile forensics

MESH is purpose-built for remote mobile device forensics, providing seamless access to Android and iOS devices.

### Android device access

- **ADB-over-WiFi** automatically acheivable on the MESH network due to CGNAT assignment
- Direct access to device shell, logcat, and system services and bugreports
- Supported APK for victim's to install and connect to the MESH
- No USB cable required for forensic collection
- Works from anywhere in the world over the mesh

### iOS device support

- Integration with **libimobiledevice** for iOS forensics
- Remote access to iOS diagnostic interfaces
- Support for backup extraction and analysis
- (Full iOS support coming Q4 2026)

!!! example "iOS forensics available in Q4 2026"

### CGNAT address assignment

- Each device gets a unique mesh IP address (100.64.0.0/10 range)
- Devices are reachable from anywhere in the mesh
- No port forwarding or NAT configuration required
- Persistent IP addresses for consistent access

### Forensic tool integration

MESH works seamlessly with industry-standard forensic tools:

- **AndroidQF** - Automated Android forensics artifact collection
- **MVT (Mobile Verification Toolkit)** - iOS and Android forensics
- **libimobiledevice** - iOS device communication
- **Custom scripts** - Any tool that works over ADB

!!! example "Example: Remote AndroidQF collection"
    ```bash
    # Connect to device over mesh
    adb connect 100.64.2.1:5555

    # Run AndroidQF spyware scan
    androidqf --adb 100.64.2.1:5555 --output ./artifacts/
    ```

## Forensic workflows

MESH provides specialized features to support complete forensic investigation workflows.

### Artifact collection

- **Bug reports** - Collect comprehensive device diagnostics
- **dumpsys output** - Extract system service information
- **System artifacts** - Pull files, databases, and logs
- **Package information** - List installed apps and permissions
- **Network configuration** - Capture network settings and connections

### Extended network analysis

- **LAN route advertisement** - Access devices on the endpoint's local network if needed
- **Subnet routing** - Investigate entire network segments remotely
- **Network packet capture** - Monitor traffic from remote locations easily
- **DNS analysis** - Inspect DNS queries and responses

### Exit node capabilities

- Route internet traffic through mesh nodes
- Monitor and analyse network behavior
- Collect PCAPs
- Detect malicious connections and C2 traffic
- Capture network forensics evidence

### Kill-Switch isolation

- Completely isolate device network traffic
- Force all traffic through the mesh
- Prevent data exfiltration during investigation
- Ensure forensic integrity

!!! info "Forensic Workflows Guide"
    See the [User guide](../user-guide/user-guide.md) for detailed forensic workflow documentation.

## Rapid deployment

MESH is designed for quick deployment in time-sensitive investigations.

### Ephemeral mesh networks

- Spin up a mesh network in minutes
- Tear down completely when investigation is complete
- No persistent infrastructure required
- Leave no trace on target devices (optional ephemeral mode)

### No complex VPN configuration

- No manual IP address assignment
- No routing table configuration
- No firewall rule management
- Everything is automatic

### Automatic NAT Traversal

- Works behind NAT without port forwarding
- Automatic STUN for NAT hole-punching
- Fallback to DERP relay when P2P fails
- No network administrator access required

### Self-Hostable control plane

- Deploy your own control plane in minutes
- Full control over your infrastructure
- No dependency on third-party services
- Docker-based deployment for easy management

!!! success "Quick deployment"
    From zero to first forensic collection in under 1 hour. See the [Getting started Guide](../getting-started/index.md).

## Security and privacy

### Access Control Lists (ACLs)

- Fine-grained control over which nodes can communicate
- User-based and tag-based access policies
- Restrict access to specific services and ports
- Audit logging for compliance

### Node isolation

- Isolate endpoint devices from each other
- Prevent lateral movement in the mesh
- Enforce least-privilege access
- Protect sensitive investigations

### Audit logging

- Complete logs of mesh activity
- Track node connections and disconnections
- Monitor access to endpoint devices
- Compliance with forensic evidence standards

## Performance

### Peer-to-Peer first

- Direct connections between nodes when possible
- Minimal latency for forensic operations
- No unnecessary relay overhead
- Automatic path optimization

### Efficient protocol

- WireGuard's high-performance cryptography
- Minimal CPU and battery impact on endpoints
- Efficient bandwidth usage
- Suitable for low-bandwidth environments

### Scalability

- Support for hundreds of nodes per mesh
- Multiple concurrent investigations
- Distributed DERP relay infrastructure
- Horizontal scaling of control plane

---

**Next:** Explore the [Architecture](architecture.md) to understand how these features are implemented →
