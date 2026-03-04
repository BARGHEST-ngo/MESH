# Control plane setup

The MESH control plane coordinates the mesh network, distributes keys, enforces ACLs, and manages HTTPS DERP relays for when UDP is not available.

## Overview

The control plane is responsible for:

- **Node Registration**: Authenticating and registering nodes
- **Key Distribution**: Distributing WireGuard public keys between peers
- **ACL Enforcement**: Enforcing access control policies
- **NAT Traversal**: Coordinating peer discovery and connection establishment
- **DERP Coordination**: Managing fallback relay connections

### Configuration

The configuration file allows you to control the key aspects of your control plane server.
If you haven't already created the file following the 'Getting Started' guide, you can do so like this:

```bash
cat > config/config.yaml << 'EOF'
server_url: https://mesh.yourdomain.com
listen_addr: 0.0.0.0:8080
metrics_listen_addr: 0.0.0.0:9090

grpc_listen_addr: 0.0.0.0:50443
grpc_allow_insecure: false

private_key_path: /var/lib/headscale/private.key
noise_private_key_path: /var/lib/headscale/noise_private.key
base_domain: mesh.yourdomain.com

ip_prefixes:
  - 100.64.0.0/10
  - fd7a:115c:a1e0::/48

derp:
  server:
    enabled: true
    region_id: 999
    region_code: "mesh"
    region_name: "MESH DERP"
    stun_listen_addr: "0.0.0.0:3478"

  urls:
    - https://controlplane.tailscale.com/derpmap/default

  auto_update_enabled: true
  update_frequency: 24h

database:
  type: sqlite3
  sqlite:
    path: /var/lib/headscale/db.sqlite

log:
  level: info
  format: text

dns_config:
  override_local_dns: false
  nameservers:
    - 1.1.1.1
    - 8.8.8.8
  magic_dns: true
  base_domain: mesh.local

acl_policy_path: /etc/headscale/acl.yaml
EOF
```

Replace `mesh.yourdomain.com` with your actual domain or IP.

## Configuration parameters explained

Understanding each configuration parameter is essential for properly securing and operating your MESH control plane.

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

#### `listen_addr`

```yaml
listen_addr: 0.0.0.0:8080
```

The address and port where Headscale's HTTP API listens internally.

**When to modify:**

- Change port if 8080 conflicts with other services
- Use `127.0.0.1:8080` if only accessing via reverse proxy on same host
- Keep `0.0.0.0` if reverse proxy is on different host

**Security implications:**

- Should not be directly exposed to the internet
- Always use a reverse proxy with HTTPS in front of this
- Firewall this port to only allow reverse proxy access

#### `metrics_listen_addr`

```yaml
metrics_listen_addr: 0.0.0.0:9090
```

Prometheus metrics endpoint for monitoring.

**When to modify:**

- Change port if 9090 conflicts
- Use `127.0.0.1:9090` to restrict to localhost only

**Security implications:**

- Metrics may contain sensitive information about your network topology
- Should not be publicly accessible
- Firewall appropriately or bind to localhost

#### `grpc_listen_addr` and `grpc_allow_insecure`

```yaml
grpc_listen_addr: 0.0.0.0:50443
grpc_allow_insecure: false
```

The gRPC API used by the `headscale` CLI tool and web UI.

**When to modify:**

- Change port if 50443 conflicts
- **Never** set `grpc_allow_insecure: true` in production

**Security implications:**

- `grpc_allow_insecure: false` requires TLS for gRPC connections
- Setting to `true` exposes administrative API without encryption
- This API has full control over your mesh network
- Firewall this port to only allow administrative access

!!! danger "Critical security setting"
    **Never** set `grpc_allow_insecure: true` in production. This would allow anyone with network access to fully control your mesh network without authentication or encryption.

#### `ip_prefixes`

```yaml
ip_prefixes:
  - 100.64.0.0/10      # IPv4 CGNAT range
  - fd7a:115c:a1e0::/48  # IPv6 ULA range
```

The IP address ranges assigned to mesh nodes.

**When to modify:**

- Rarely needs modification
- Change only if these ranges conflict with existing networks
- Must use private/CGNAT ranges to avoid internet routing conflicts
- You can use these when want to seperate IP ranges in cases.

**Security implications:**

- These addresses are only routable within the mesh
- Using public IP ranges would cause routing conflicts
- IPv6 range should be a Unique Local Address (ULA) range

### Database configuration

#### `database`

```yaml
database:
  type: sqlite3
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
- Ensure regular backups
- Protect database file with appropriate filesystem permissions
- Consider encrypting the filesystem or using encrypted storage

### DNS configuration

#### `dns_config`

```yaml
dns_config:
  override_local_dns: false
  nameservers:
    - 1.1.1.1
    - 8.8.8.8
  magic_dns: true
  base_domain: mesh.local
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

MESH includes a built-in DERP server for convenience.By default, the controlplane acting as the DERP relay is set to false. This is because it automatically uses the [Tailscale DERP relays](https://login.tailscale.com/derpmap/default).

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

  auto_update_enabled: true
  update_frequency: 24h
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
  update_frequency: 24h
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
  update_frequency: 24h
```

## Access Control Lists (ACLs)

ACLs are the **most critical security component** of your MESH deployment. They control which nodes can communicate with each other and on which ports.

!!! danger "Critical security warning"
    **Never use permissive ACLs (`"*"` rules) in production.** This allows any compromised endpoint to access all other devices in your mesh, including other forensic targets and analyst workstations. Always implement the principle of least privilege.

### Why ACLs are critical for MESH

In a forensic mesh network, security isolation is paramount:

**Threat model:**

- **Compromised endpoints**: Target devices may be infected with malware
- **Lateral movement**: Malware could attempt to spread to other devices
- **Data exfiltration**: Compromised devices could try to access other forensic targets
- **Analyst protection**: Analyst workstations must be protected from compromised endpoints

**Security goals:**

- Endpoints should **only** be accessible by authorised analysts using a secure collection node
- Endpoints should **never** be able to access each other
- Endpoints should **never** be able to access analyst workstations
- Analysts should have controlled access to specific endpoints
- All access should be logged and auditable

### ACL syntax and structure

ACLs are defined in YAML format with the following structure:

```yaml
# Define groups of users/nodes
groups:
  group:name:
    - user1
    - user2

# Define access control rules
acls:
  - action: accept  # or deny
    src:
      - source_group_or_user
    dst:
      - destination_group_or_user:port
```

**Components:**

**`groups`**: Logical groupings of users or nodes

- Format: `group:name`
- Contains list of usernames or node names
- Makes ACL rules more maintainable
- Can be nested or referenced in rules

**`acls`**: List of access control rules

- Evaluated in order (first match wins)
- Each rule has `action`, `src`, and `dst`

**`action`**: What to do with matching traffic

- `accept`: Allow the traffic
- `deny`: Block the traffic

**`src`**: Source of the traffic

- Can be: username, node name, group, or `"*"` (any)
- Multiple sources can be listed

**`dst`**: Destination of the traffic

- Format: `target:port` or `target:*` (all ports)
- Can be: username, node name, group, or `"*"` (any)
- Port can be specific (e.g., `5555` for ADB) or `*` (all ports)

### Testing vs. production ACLs

#### Testing ACL (permissive - DO NOT USE IN PRODUCTION)

```bash
cat > config/acl.yaml << 'EOF'
# WARNING: This is for testing only!
# Allows all traffic between all nodes
acls:
  - action: accept
    src:
      - "*"
    dst:
      - "*:*"
EOF
```

!!! warning "Testing only"
    This permissive ACL is **only** for initial testing and development. It provides **no security** and should **never** be used in production or with real forensic targets.

**When to use:**

- Initial setup and testing
- Verifying connectivity between nodes
- Troubleshooting connection issues

**When to remove:**

- Before connecting any real forensic targets
- Before deploying in production
- Before handling sensitive data

#### Production ACL (restrictive - RECOMMENDED)

```yaml
# Production ACL - Principle of least privilege
groups:
  group:analysts:
    - analyst1
    - analyst2
  group:endpoints:
    - android-device-1
    - android-device-2

acls:
  # Analysts can access all endpoints on ADB port
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:5555  # ADB port only

  # Analysts can access all endpoints on all ports (for full forensics)
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*

  # Endpoints CANNOT access each other (critical for isolation)
  - action: deny
    src:
      - group:endpoints
    dst:
      - group:endpoints:*

  # Endpoints CANNOT access analysts (prevent reverse connections)
  - action: deny
    src:
      - group:endpoints
    dst:
      - group:analysts:*

  # Analysts can access each other (for collaboration)
  - action: accept
    src:
      - group:analysts
    dst:
      - group:analysts:*

  # Deny all other traffic (default deny)
  - action: deny
    src:
      - "*"
    dst:
      - "*:*"
```

### Real-world ACL examples

#### Example 1: Single analyst, multiple endpoints

```yaml
# Simple forensic setup
groups:
  group:analysts:
    - forensic-workstation

  group:endpoints:
    - suspect-phone-1
    - suspect-phone-2
    - suspect-phone-3

acls:
  # Analyst can access all endpoints
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*

  # Endpoints are completely isolated
  - action: deny
    src:
      - group:endpoints
    dst:
      - "*:*"
```

#### Example 2: Multi-team deployment

```yaml
# Multiple analyst teams with separate endpoints
groups:
  group:team-alpha-analysts:
    - analyst-alice
    - analyst-bob

  group:team-alpha-endpoints:
    - case-123-phone-1
    - case-123-phone-2

  group:team-bravo-analysts:
    - analyst-charlie
    - analyst-diana

  group:team-bravo-endpoints:
    - case-456-phone-1
    - case-456-phone-2

acls:
  # Team Alpha: analysts can access their endpoints
  - action: accept
    src:
      - group:team-alpha-analysts
    dst:
      - group:team-alpha-endpoints:*

  # Team Bravo: analysts can access their endpoints
  - action: accept
    src:
      - group:team-bravo-analysts
    dst:
      - group:team-bravo-endpoints:*

  # Team Alpha: analysts can collaborate
  - action: accept
    src:
      - group:team-alpha-analysts
    dst:
      - group:team-alpha-analysts:*

  # Team Bravo: analysts can collaborate
  - action: accept
    src:
      - group:team-bravo-analysts
    dst:
      - group:team-bravo-analysts:*

  # Teams CANNOT access each other's endpoints
  - action: deny
    src:
      - group:team-alpha-analysts
    dst:
      - group:team-bravo-endpoints:*

  - action: deny
    src:
      - group:team-bravo-analysts
    dst:
      - group:team-alpha-endpoints:*

  # Endpoints are completely isolated
  - action: deny
    src:
      - group:team-alpha-endpoints
    dst:
      - "*:*"

  - action: deny
    src:
      - group:team-bravo-endpoints
    dst:
      - "*:*"

  # Default deny
  - action: deny
    src:
      - "*"
    dst:
      - "*:*"
```

#### Example 3: Port-specific access (ADB only)

```yaml
# Restrict analysts to only ADB access
groups:
  group:analysts:
    - forensic-workstation

  group:endpoints:
    - android-device-1
    - android-device-2

acls:
  # Analysts can ONLY access ADB port (5555)
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:5555

  # Deny all other ports
  - action: deny
    src:
      - group:analysts
    dst:
      - group:endpoints:*

  # Endpoints completely isolated
  - action: deny
    src:
      - group:endpoints
    dst:
      - "*:*"
```

### ACL best practices for forensic deployments

#### 1. Principle of least privilege

**Always grant the minimum access required:**

- Analysts should only access endpoints they're investigating
- Limit port access to only what's needed (e.g., ADB port 5555)
- Use time-limited access when possible (rotate keys regularly)

#### 2. Endpoint isolation

**Endpoints must never access each other:**

```yaml
# CRITICAL: Always include this rule
- action: deny
  src:
    - group:endpoints
  dst:
    - group:endpoints:*
```

**Why this matters:**

- Prevents malware from spreading between forensic targets
- Stops compromised devices from accessing other evidence
- Maintains chain of custody integrity

#### 3. Prevent reverse connections

**Endpoints should not initiate connections to analysts:**

```yaml
# Prevent endpoints from connecting to analysts
- action: deny
  src:
    - group:endpoints
  dst:
    - group:analysts:*
```

**Why this matters:**

- Prevents malware from attacking analyst workstations
- Stops data exfiltration from endpoints
- Maintains unidirectional trust model

#### 4. Use groups for maintainability

**Bad (hard to maintain):**

```yaml
acls:
  - action: accept
    src:
      - analyst1
      - analyst2
      - analyst3
    dst:
      - phone1:*
      - phone2:*
      - phone3:*
```

**Good (easy to maintain):**

```yaml
groups:
  group:analysts:
    - analyst1
    - analyst2
    - analyst3
  group:endpoints:
    - phone1
    - phone2
    - phone3

acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*
```

#### 5. Default deny at the end

**Always end with a default deny rule:**

```yaml
acls:
  # ... your specific rules ...

  # Default deny (catches anything not explicitly allowed)
  - action: deny
    src:
      - "*"
    dst:
      - "*:*"
```

### Testing and validating ACLs

After creating or modifying ACLs, always test them:

#### 1. Apply the new ACL configuration

```bash
# Copy new ACL file
cp acl.yaml config/acl.yaml

# Restart Headscale to apply
docker-compose restart headscale
```

#### 2. Test from analyst workstation

```bash
# Should succeed: Analyst accessing endpoint
ping 100.64.2.1  # Endpoint mesh IP

# Should succeed: ADB connection
adb connect 100.64.2.1:5555
```

#### 3. Test endpoint isolation

```bash
# On endpoint device (via ADB shell)
adb shell

# Should FAIL: Endpoint trying to access another endpoint
ping 100.64.2.2  # Another endpoint's mesh IP

# Should FAIL: Endpoint trying to access analyst
ping 100.64.1.1  # Analyst mesh IP
```

#### 4. Check logs

```bash
# View ACL enforcement logs
docker logs headscale | grep -i acl

# Look for denied connections
docker logs headscale | grep -i denied
```

### Common ACL mistakes and how to avoid them

#### Mistake 1: Using `"*"` in production

**Wrong:**

```yaml
acls:
  - action: accept
    src:
      - "*"
    dst:
      - "*:*"
```

**Right:**

```yaml
acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*
```

#### Mistake 2: Forgetting endpoint isolation

**Wrong (endpoints can access each other):**

```yaml
acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*
```

**Right (endpoints explicitly denied):**

```yaml
acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*

  - action: deny
    src:
      - group:endpoints
    dst:
      - group:endpoints:*
```

#### Mistake 3: Allowing reverse connections

**Wrong (endpoints can connect to analysts):**

```yaml
acls:
  - action: accept
    src:
      - "*"
    dst:
      - "*:*"
```

**Right (only analysts initiate connections):**

```yaml
acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*

  - action: deny
    src:
      - group:endpoints
    dst:
      - group:analysts:*
```

#### Mistake 4: Overly broad port access

**Wrong (all ports open):**

```yaml
acls:
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:*
```

**Better (specific ports only):**

```yaml
acls:
  # ADB access only
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:5555

  # SSH access only (if needed)
  - action: accept
    src:
      - group:analysts
    dst:
      - group:endpoints:22
```

### Security considerations

#### Pre-authentication key rotation

Pre-auth keys should be rotated regularly:

```bash
# List existing keys
docker exec headscale headscale preauthkeys list

# Expire old keys
docker exec headscale headscale preauthkeys expire --user default KEY_ID

# Create new key
docker exec headscale headscale preauthkeys create --user default --reusable --expiration 24h
```

**Best practices:**

- Use short expiration times (24-48 hours)
- Create new keys for each investigation
- Expire keys immediately after use if not reusable
- Never share keys between teams or cases

#### Node removal

Remove nodes when investigations complete using the UI or via the backend via:

```bash
# List nodes
docker exec headscale headscale nodes list

# Remove a node
docker exec headscale headscale nodes delete --identifier NODE_ID
```

**When to remove nodes:**

- Investigation completed
- Device returned to owner
- Node compromised or suspected compromise
- Regular cleanup of old investigations

#### Audit logging

Monitor ACL enforcement:

```bash
# Enable detailed logging in config.yaml
log:
  level: debug  # or trace for maximum detail
  format: json  # easier to parse

# Monitor logs in real-time
docker logs -f headscale | grep -i acl
```

**What to monitor:**

- Denied connection attempts (potential compromise)
- Unusual access patterns
- Failed authentication attempts
- ACL rule changes

### Applying ACL changes

After modifying ACLs, apply the changes:

```bash
# Method 1: Restart Headscale (applies immediately)
docker-compose restart headscale

# Method 2: Reload ACLs without restart (if supported)
docker exec headscale headscale policy reload
```

**Important notes:**

- ACL changes apply to new connections immediately
- Existing connections may remain until they reconnect
- Test ACL changes in a development environment first
- Keep backups of working ACL configurations

### Restart Services

```bash
docker-compose restart
```

## Exposing the control plane

The control plane needs to be accessible from the internet for nodes to connect. You have several options:

### Option 1: Reverse Proxy (Recommended)

Use Nginx or Caddy as a reverse proxy with automatic HTTPS:

#### Caddy (Easiest)

```bash
# Install Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy

# Create Caddyfile
sudo cat > /etc/caddy/Caddyfile << 'EOF'
mesh.yourdomain.com {
    reverse_proxy localhost:8080
}
EOF

# Restart Caddy
sudo systemctl restart caddy
```

#### Nginx

```bash
# Install Nginx and Certbot
sudo apt install nginx certbot python3-certbot-nginx

# Create Nginx config
sudo cat > /etc/nginx/sites-available/mesh << 'EOF'
server {
    listen 80;
    server_name mesh.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/mesh /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Get SSL certificate
sudo certbot --nginx -d mesh.yourdomain.com
```

### Option 2: Ngrok (Ephemeral/Development/Testing)

For a quick ephemeral session and/or when you don't have a public IP you can use services like NGROK:

```bash
# Install ngrok
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
sudo apt update
sudo apt install ngrok

# Authenticate (get token from ngrok.com)
ngrok config add-authtoken YOUR_TOKEN

# Expose Headscale
ngrok http 8080
```

Use the ngrok URL (e.g., `https://abc123.ngrok.io`) as your `server_url` in the Headscale config.

## Backup and Restore

### Backup

```bash
# Stop Headscale
docker-compose stop headscale

# Backup database and config
tar -czf mesh-backup-$(date +%Y%m%d).tar.gz config/ data/

# Restart Headscale
docker-compose start headscale
```

### Restore

```bash
# Stop Headscale
docker-compose stop headscale

# Restore from backup
tar -xzf mesh-backup-20240101.tar.gz

# Restart Headscale
docker-compose start headscale
```

## Troubleshooting

### Headscale won't start

**Check logs:**

```bash
docker logs headscale
```

**Common issues:**

- Port already in use: `sudo lsof -i :8080`
- Invalid config: Check `config/config.yaml` syntax
- Permission issues: `sudo chown -R 1000:1000 data/`

### Nodes can't connect

**Verify the control plane is accessible:**

```bash
curl https://mesh.yourdomain.com/health
```

**Check firewall:**

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 3478/udp
```

**Check DNS:**

```bash
nslookup mesh.yourdomain.com
```

### DERP Relay Not Working

**Check STUN port:**

```bash
sudo netstat -ulnp | grep 3478
```

**Test STUN:**

```bash
# Install stun client
sudo apt install stun-client

# Test STUN
stun mesh.yourdomain.com 3478
```

## Security best practices

1. **Use HTTPS**: Always use a reverse proxy with SSL/TLS
2. **Restrict API access**: Don't expose port 8080 publicly
3. **Use strong ACLs**: Implement least-privilege access control
4. **Rotate API keys**: Regularly expire and recreate API keys
5. **Monitor logs**: Watch for suspicious activity
6. **Backup regularly**: Automate database backups
7. **Update frequently**: Keep Headscale and Docker updated
8. **Use firewall**: Only allow necessary ports

## Advanced configuration

For advanced control plane configurations, see the [Advanced control plane documentation](../advanced/control-plane-advanced.md):

- **[Custom DERP servers](../advanced/control-plane-advanced.md#custom-derp-servers)** - Deploy your own relay infrastructure for maximum privacy and geographic distribution
- **[PostgreSQL database](../advanced/control-plane-advanced.md#postgresql-database)** - Use PostgreSQL for large deployments and high availability
- **[Exit node and packet capture](../advanced/exit-node-pcap.md)** - Configure exit node routing and capture network traffic for forensic analysis

These advanced features are recommended for:

- Production deployments with high node counts
- Geographic distribution requirements
- Maximum privacy and control
- Network forensics and traffic analysis

## Next steps

- **[Getting started](../getting-started/index.md)** - Set up your first mesh network
- **[Analyst client](analyst-client.md)** - Install the analyst client
- **[Endpoint client](endpoint-clients.md)** - Deploy to Android devices
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
