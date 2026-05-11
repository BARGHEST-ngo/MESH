# Advanced control plane configuration

Advanced configuration options for MESH control plane deployments, including custom DERP servers.

!!! warning "Advanced users only"
    These configurations are for advanced users with production deployments or specific requirements. Most users should use the default configuration from the [Getting started guide](../setup/control-plane.md).

## Prerequisites

Before implementing advanced configurations:

1. **Working deployment** - Complete the [Getting started guide](../setup/index.md)
2. **Basic understanding** - Familiar with [control plane configuration](../setup/control-plane.md)
3. **Production experience** - Tested MESH in development environment
4. **Technical skills** - Comfortable with Linux, networking, Docker, and YAML

## Custom DERP servers

Deploy your own DERP relay servers for maximum privacy, geographic distribution, or increased capacity.

### Why custom DERP servers?

**Use cases:**

- **Maximum privacy** - Full control over relay infrastructure, no third-party dependencies
- **Geographic distribution** - Deploy DERP servers closer to your users for lower latency
- **Increased capacity** - Handle more connections than public DERP servers
- **Censorship resistance** - Deploy in multiple jurisdictions to avoid blocking
- **Compliance** - Meet data sovereignty or regulatory requirements

**Trade-offs:**

- ➕ Full control over relay infrastructure
- ➕ Reduced latency with geographic distribution
- ➕ No dependency on third-party DERP servers
- ➕ Can scale to handle high connection volumes
- ➖ Additional infrastructure to deploy and maintain
- ➖ Increased operational complexity
- ➖ Requires public IP addresses and open ports
- ➖ Additional hosting costs

### Disable built-in DERP and external maps

First, configure Headscale to use only your custom DERP servers:

```yaml
# control-plane/headscale/config.yaml
derp:
  server:
    enabled: false  # Disable built-in DERP server
  
  urls: []  # Don't use Tailscale's DERP map
  
  paths:
    - /etc/headscale/derp.yaml  # Path to custom DERP map
  
  auto_update_enabled: false  # Don't auto-update (using custom map)
```

### Create custom DERP map

Create a custom DERP map file defining your relay servers:

```yaml
regions:
  900:
    regionid: 900
    regioncode: eu-west
    regionname: Europe West
    nodes:
      - name: eu-west-1
        regionid: 900
        hostname: derp-eu.yourdomain.com
        ipv4: 1.2.3.4
        stunport: 3478
        stunonly: false
        derpport: 443
  
  901:
    regionid: 901
    regioncode: us-east
    regionname: US East
    nodes:
      - name: us-east-1
        regionid: 901
        hostname: derp-us.yourdomain.com
        ipv4: 5.6.7.8
        stunport: 3478
        stunonly: false
        derpport: 443
  
  902:
    regionid: 902
    regioncode: asia-pacific
    regionname: Asia Pacific
    nodes:
      - name: asia-1
        regionid: 902
        hostname: derp-asia.yourdomain.com
        ipv4: 9.10.11.12
        stunport: 3478
        stunonly: false
        derpport: 443
```

### DERP map parameters

**Region parameters:**

- `regionid` - Unique identifier for this region (use 900+ for custom regions)
- `regioncode` - Short code for the region (e.g., "eu-west", "us-east")
- `regionname` - Human-readable name shown to users

**Node parameters:**

- `name` - Unique name for this DERP server
- `regionid` - Must match the parent region ID
- `hostname` - Public hostname of the DERP server (must have valid TLS certificate)
- `ipv4` - Public IPv4 address of the DERP server
- `stunport` - STUN server port (default: 3478, UDP)
- `stunonly` - If true, only provides STUN (NAT traversal), not DERP relay
- `derpport` - DERP relay port (default: 443, TCP with TLS)

### Deploy DERP server

Deploy a standalone DERP server on each host:

```bash
# On each DERP server host
docker run -d \
  --name derp \
  --restart unless-stopped \
  -p 443:443 \
  -p 3478:3478/udp \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e DERP_DOMAIN=derp-eu.yourdomain.com \
  -e DERP_CERT_DIR=/etc/letsencrypt/live/derp-eu.yourdomain.com \
  -e DERP_ADDR=:443 \
  -e DERP_STUN=true \
  -e DERP_STUN_PORT=3478 \
  ghcr.io/tailscale/derper:latest
```

### TLS certificates for DERP servers

DERP servers require valid TLS certificates. Use Let's Encrypt:

```bash
# Install certbot
sudo apt install -y certbot

# Obtain certificate
sudo certbot certonly --standalone \
  -d derp-eu.yourdomain.com \
  --agree-tos \
  --email admin@yourdomain.com

# Certificate will be at:
# /etc/letsencrypt/live/derp-eu.yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/derp-eu.yourdomain.com/privkey.pem
```

### Firewall configuration for DERP servers

DERP servers require specific ports to be open:

```bash
# Allow DERP relay (TCP/443)
sudo ufw allow 443/tcp

# Allow STUN (UDP/3478)
sudo ufw allow 3478/udp

# Enable firewall
sudo ufw enable
```

**Required ports:**

- **TCP/443** - DERP relay (HTTPS/TLS)
- **UDP/3478** - STUN server (NAT traversal)

### Update control plane configuration

Mount the custom DERP map in your `compose.yml`:

```yaml
services:
  headscale:
    image: docker.io/headscale/headscale:0.28.0
    restart: unless-stopped
    volumes:
      - ${LOCAL_WORKSPACE_FOLDER:-.}/control-plane/headscale:/etc/headscale
      - headscale-data:/var/lib/headscale
      - ./derp.yaml:/etc/headscale/derp.yaml  # Add custom DERP map
    # ... rest of configuration
```

Restart Headscale to apply changes:

```bash
docker compose restart headscale
```

### Verify DERP server connectivity

Test DERP server connectivity:

```bash
# Test STUN server
nc -u -v derp-eu.yourdomain.com 3478

# Test DERP relay (should show TLS handshake)
openssl s_client -connect derp-eu.yourdomain.com:443

# Check from Headscale logs
docker compose logs headscale | grep -i derp
```

### DERP server placement strategy

**Single region (simple):**

- One DERP server in your primary location
- Suitable for small deployments or single-region operations

**Multi-region (recommended):**

- DERP servers in multiple geographic locations
- Clients automatically use closest DERP server
- Provides redundancy and lower latency

**Example deployment:**

```
Region 900 (Europe): derp-eu.yourdomain.com
Region 901 (US East): derp-us.yourdomain.com
Region 902 (Asia):    derp-asia.yourdomain.com
```

### Monitoring DERP servers

Monitor DERP server health and usage:

```bash
# Check DERP server logs
docker logs derp

# Monitor active connections
docker exec derp netstat -an | grep :443

# Check STUN server
docker exec derp netstat -an | grep :3478
```

### Security considerations for custom DERP

**TLS certificates:**

- Always use valid TLS certificates (Let's Encrypt recommended)
- Never use self-signed certificates (clients will reject connection)
- Automate certificate renewal (certbot with cron)

**Network security:**

- Only expose required ports (443/tcp, 3478/udp)
- Use firewall to restrict access if needed
- Consider DDoS protection for public DERP servers

**Privacy:**

- DERP servers can see encrypted traffic metadata (source, destination, timing)
- Cannot decrypt traffic contents (end-to-end encrypted)
- Host DERP servers in trusted locations/jurisdictions
- Consider logging policies and data retention

**Capacity planning:**

- Each DERP server can handle hundreds of concurrent connections
- Monitor CPU and bandwidth usage
- Scale horizontally by adding more DERP servers
- Use multiple DERP servers per region for redundancy

## Troubleshooting advanced configurations

### Custom DERP servers not working

**Symptoms:**

- Clients can't connect through DERP
- Headscale logs show DERP errors
- Connections fail when direct P2P fails

**Solutions:**

1. **Verify DERP server is accessible:**

   ```bash
   # Test HTTPS connectivity
   curl -v https://derp-eu.yourdomain.com

   # Test STUN server
   nc -u -v derp-eu.yourdomain.com 3478
   ```

2. **Check TLS certificate:**

   ```bash
   # Verify certificate is valid
   openssl s_client -connect derp-eu.yourdomain.com:443 -servername derp-eu.yourdomain.com
   ```

3. **Verify DERP map syntax:**

   ```bash
   # Check YAML syntax
   docker compose exec headscale cat /etc/headscale/derp.yaml
   ```

4. **Check Headscale logs:**

   ```bash
   docker logs headscale | grep -i derp
   ```

## Next steps

- **[Exit node and PCAP](exit-node-pcap.md)** - Configure exit node and packet capture
- **[AmneziaWG](amneziawg.md)** - Configure obfuscation for censorship resistance
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions

---

← [Back to Advanced](index.md) | [Next: Exit node and PCAP](exit-node-pcap.md) →
