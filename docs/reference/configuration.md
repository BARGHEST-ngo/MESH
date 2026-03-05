# Control Plane Configuration

The Headscale configuration file allows you to control the key aspects of your control plane server. This file is created automatically when you run `task controlPlane` for the first time. You can also create it manually if you prefer.

```yaml
server_url: https://mesh.yourdomain.com
listen_addr: 0.0.0.0:8080
metrics_listen_addr: 127.0.0.1:9090
grpc_listen_addr: 127.0.0.1:50443
grpc_allow_insecure: false

noise:
  private_key_path: /var/lib/headscale/noise_private.key

prefixes:
  v4: 100.64.0.0/10
  v6: fd7a:115c:a1e0::/48
  allocation: sequential

derp:
  server:
    enabled: false
    region_id: 999
    region_code: "headscale"
    region_name: "Headscale Embedded DERP"
    verify_clients: true
    stun_listen_addr: "0.0.0.0:3478"
    private_key_path: /var/lib/headscale/derp_server_private.key
    automatically_add_embedded_derp_region: true
    ipv4: 198.51.100.1
    ipv6: 2001:db8::1

  urls:
   - https://controlplane.tailscale.com/derpmap/default

  paths: []
  auto_update_enabled: false
  update_frequency: 3h

disable_check_updates: false
ephemeral_node_inactivity_timeout: 30m

database:
  type: sqlite
  debug: false
  gorm:
    prepare_stmt: true
    parameterized_queries: true
    skip_err_record_not_found: true
    slow_threshold: 1000

  sqlite:
    path: /var/lib/headscale/db.sqlite
    write_ahead_log: true
    wal_autocheckpoint: 1000

log:
  level: info
  format: text

policy:
  mode: database
  path: ""

dns:
  magic_dns: true
  base_domain: mesh.local
  override_local_dns: true
  nameservers:
    global:
      - 1.1.1.1
      - 1.0.0.1
      - 2606:4700:4700::1111
      - 2606:4700:4700::1001
    split: {}

  search_domains: []
  extra_records: []

unix_socket: /var/run/headscale/headscale.sock
unix_socket_permission: "0770"

logtail:
  enabled: false

randomize_client_port: false
```

Replace `mesh.yourdomain.com` with your actual domain or IP.

## Configuration parameters explained

Full documentation for each of the configuration parameters can be found in the [example config provided in the MESH repository](https://github.com/BARGHEST-ngo/MESH/blob/main/control-plane/headscale/config.example.yaml). Most of the values should not be changed, but those you may want to change are discussed below.

### Server and network settings

#### `server_url`

```yaml
server_url: https://mesh.yourdomain.com
```

The public URL/IP where your control plane is accessible. This is the URL/IP that clients will use to connect.

**When to modify:**

- Always set this to your actual domain name or public IP
- Must match your reverse proxy or DNS configuration
- Must use HTTPS in production (HTTP only for local testing)

**Security implications:**

- Clients validate this URL during connection
- Using HTTP exposes authentication tokens to network eavesdropping
- Must be accessible from all networks where clients will connect

#### `prefixes`

```yaml
prefixes:
  v4: 100.64.0.0/10
  v6: fd7a:115c:a1e0::/48
```

The IP address ranges assigned to mesh nodes.

**When to modify:**

- Rarely needs modification
- Change only if these ranges conflict with existing networks
- Must use private/CGNAT ranges to avoid internet routing conflicts
- You can use these when you want to seperate IP ranges in cases.

**Security implications:**

- These addresses are only routable within the mesh
- Using public IP ranges would cause routing conflicts
- IPv6 range should be a Unique Local Address (ULA) range

### Database configuration

#### `database`

```yaml
database:
  type: sqlite
  sqlite:
    path: /var/lib/headscale/db.sqlite
```

Database backend for storing mesh state.

**When to modify:**

- SQLite is fine for most deployments (hundreds of nodes)
- Consider PostgreSQL for very large deployments (thousands of nodes)
- Ensure the path is on persistent storage (not tmpfs)

**Security implications:**

- Database contains all mesh configuration and node keys
- Ensure regular backups for persistent control plane deployments
- Protect database file with appropriate filesystem permissions
- Consider encrypting the filesystem or using encrypted storage

### DNS configuration

#### `dns`

```yaml
dns:
  magic_dns: true
  base_domain: mesh.local
  override_local_dns: true
  nameservers:
    global:
      - 1.1.1.1
      - 8.8.8.8
```

DNS settings for mesh nodes.

**Parameters:**

- `override_local_dns`: Whether to override client's local DNS settings
- `nameservers`: DNS servers to use for external queries
- `magic_dns`: Enable hostname resolution within the mesh (e.g., `android-device-1.mesh.local`)
- `base_domain`: Domain suffix for MagicDNS hostnames

**When to modify:**

- Change `nameservers` to use your preferred DNS servers such as Quad9.
- Set `override_local_dns: true` to force all DNS through the mesh (useful for monitoring)
- Change `base_domain` to match your organisation

**Security implications:**

- `override_local_dns: true` routes all DNS queries through the mesh
- Can be used to monitor or filter DNS queries
- MagicDNS makes it easier to identify nodes but reveals hostnames

## DERP server configuration

DERP (Designated Encrypted Relay for Packets) is a critical component for ensuring connectivity in challenging network environments.

### What is DERP?

DERP is an encrypted relay protocol that acts as a fallback when direct peer-to-peer connections fail. Think of it as a "last resort" relay that ensures your mesh network remains connected even in hostile network conditions.

**When DERP is used:**

- Both peers are behind restrictive NAT that prevents hole-punching
- UDP traffic is blocked by firewalls
- Direct peer-to-peer connection fails
- Networks with aggressive DPI that blocks WireGuard
- Censored networks that block VPN protocols

**How it works:**

1. Clients attempt direct peer-to-peer connection first (fastest)
2. If direct connection fails, traffic is relayed through DERP server
3. All traffic remains end-to-end encrypted (DERP server cannot decrypt)
4. DERP uses HTTPS, making it harder to block than UDP-based VPNs

### Built-in DERP server

```yaml
derp:
  server:
    enabled: false
    region_id: 999
    region_code: "MESH"
    region_name: "PEOPLES_REPUBLIC_OF_BARGHEST"
    stun_listen_addr: "0.0.0.0:3478"
```

MESH includes a built-in DERP server for convenience. By default, the control plane acting as the DERP relay is set to false. This is because it automatically uses the [Tailscale DERP relays](https://login.tailscale.com/derpmap/default).

!!! danger "Is DERP safe?"
    MESH'sarchitecture is end-to-end encrypted using WireGuard between your devices. That holds whether traffic is direct or relayed through DERP. Because of that DERP can see:

    - Source and destination IPs and ports (on the public Internet)
    - Which nodes are talking via that relay (public keys + node IDs)
    - Packet sizes and timing 
    - That packets belong to WireGuard/AmneziaWG flows
    
    DERP cannot see:
    
    - The packet payload (it’s WireGuard-encrypted)
    - Any inner IPs, TCP ports, or application-level data
    - Any authentication material (private keys never leave devices) 
    
    If a DERP operator tried to tamper with packets:
    
    - Modifying packets breaks MAC/authentication; endpoints drop them
    - Injecting new packets without keys is not possible in a way that passes WireGuard’s crypto
    
    So DERP can drop or delay traffic (availability), and can perform some netflow analysis, but it cannot decrypt or undetectably alter your data (confidentiality/integrity).

You can also host multiple DERP relays that will work with you MESH infastructure following this [guide](https://github.com/tailscale/tailscale/blob/main/cmd/derper/README.md).

**Parameters explained:**

- `enabled`: Whether to run the built-in DERP server
- `region_id`: Unique identifier for this DERP region (use 900+ for custom servers)
- `region_code`: Short code for this region (e.g., "mesh", "eu-west", "us-east")
- `region_name`: Human-readable name shown to users
- `stun_listen_addr`: STUN server address for NAT traversal (UDP port 3478)

**When to enable:**

- Enable if you want a self-contained deployment
- Enable for censored networks where external DERP servers may be blocked
- Disable if using external DERP servers only

**Network requirements:**

- DERP server must be publicly accessible on HTTPS (port 443)
- STUN server must be accessible on UDP port 3478
- Ensure firewall allows both TCP/443 and UDP/3478

**Security implications:**

- DERP server can see encrypted traffic metadata (source, destination, timing)
- Cannot decrypt traffic contents (end-to-end encrypted)
- Should be hosted in a trusted location
- Consider geographic placement for latency and jurisdiction

### External DERP servers

```yaml
derp:
  urls:
    - https://controlplane.tailscale.com/derpmap/default

  auto_update_enabled: false
  update_frequency: 3h
```

You can use external DERP servers in addition to or instead of the built-in server.

**Parameters explained:**

- `urls`: List of DERP map URLs to fetch
- `auto_update_enabled`: Automatically fetch updated DERP server lists
- `update_frequency`: How often to check for updates

**When to use external DERP servers:**

- Geographic distribution: Use servers closer to your users
- Redundancy: Multiple DERP servers provide failover
- Avoiding hosting costs: Use Tailscale's public DERP servers
- Censorship resistance: Distribute across multiple providers/jurisdictions

**When to use only built-in DERP:**

- Maximum privacy: No external dependencies
- Air-gapped or restricted networks
- Compliance requirements prohibit external relays
- Censored networks where external servers are blocked

**Security considerations:**

- External DERP servers can see connection metadata
- Tailscale's public DERP servers are well-maintained but operated by a third party
- For sensitive forensic work, consider using only your own DERP servers
- Multiple DERP servers provide redundancy but increase metadata exposure

### DERP server placement strategy

For optimal performance and reliability:

**Single region deployment:**

```yaml
derp:
  server:
    enabled: true
    region_id: 999
    region_code: "mesh"
    region_name: "MESH Primary"
    stun_listen_addr: "0.0.0.0:3478"

  urls: []  # Don't use external DERP servers
  auto_update_enabled: false
```

**Multi-region deployment:**

```yaml
derp:
  server:
    enabled: true
    region_id: 999
    region_code: "eu-west"
    region_name: "MESH Europe"
    stun_listen_addr: "0.0.0.0:3478"

  urls:
    - https://derp.us-east.yourdomain.com/derpmap
    - https://derp.asia.yourdomain.com/derpmap

  auto_update_enabled: true
  update_frequency: 3h
```

**Hybrid deployment (own + Tailscale's):**

```yaml
derp:
  server:
    enabled: true
    region_id: 999
    region_code: "mesh"
    region_name: "MESH Primary"
    stun_listen_addr: "0.0.0.0:3478"

  urls:
    - https://controlplane.tailscale.com/derpmap/default

  auto_update_enabled: true
  update_frequency: 3h
```
