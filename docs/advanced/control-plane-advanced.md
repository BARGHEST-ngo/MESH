# Advanced control plane configuration

Advanced configuration options for MESH control plane deployments, including custom DERP servers and PostgreSQL database setup.

!!! warning "Advanced users only"
    These configurations are for advanced users with production deployments or specific requirements. Most users should use the default configuration from the [Getting started guide](../getting-started/control-plane.md).

## Prerequisites

Before implementing advanced configurations:

1. **Working deployment** - Complete the [Getting started guide](../getting-started/index.md)
2. **Basic understanding** - Familiar with [control plane configuration](../installation/control-plane.md)
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
# config/config.yaml
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

```bash
cat > config/derp.yaml << 'EOF'
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
EOF
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

Mount the custom DERP map in your `docker-compose.yml`:

```yaml
services:
  headscale:
    image: headscale/headscale:latest
    container_name: headscale
    restart: unless-stopped
    volumes:
      - ./config:/etc/headscale
      - ./data:/var/lib/headscale
      - ./derp.yaml:/etc/headscale/derp.yaml  # Add custom DERP map
    # ... rest of configuration
```

Restart Headscale to apply changes:

```bash
docker-compose restart headscale
```

### Verify DERP server connectivity

Test DERP server connectivity:

```bash
# Test STUN server
nc -u -v derp-eu.yourdomain.com 3478

# Test DERP relay (should show TLS handshake)
openssl s_client -connect derp-eu.yourdomain.com:443

# Check from Headscale logs
docker logs headscale | grep -i derp
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

## PostgreSQL database

For production deployments with high node counts or high availability requirements, use PostgreSQL instead of SQLite.

### When to use PostgreSQL

**Use PostgreSQL when:**

- ✅ Large deployments (1000+ nodes)
- ✅ High availability requirements
- ✅ Need database replication
- ✅ Multiple Headscale instances (clustering)
- ✅ Advanced backup and recovery needs

**Use SQLite when:**

- ✅ Small to medium deployments (< 1000 nodes)
- ✅ Single Headscale instance
- ✅ Simplicity is preferred
- ✅ Lower operational overhead

**Trade-offs:**

- ➕ Better performance at scale
- ➕ Support for replication and high availability
- ➕ Advanced backup and recovery options
- ➕ Better concurrent access handling
- ➖ More complex to deploy and maintain
- ➖ Additional infrastructure required
- ➖ Higher resource usage

### Configure control plane for PostgreSQL

Update `config/config.yaml`:

```yaml
database:
  type: postgres
  postgres:
    host: postgres          # PostgreSQL hostname
    port: 5432              # PostgreSQL port
    name: headscale         # Database name
    user: headscale         # Database user
    pass: your-secure-password  # Database password
    max_open_conns: 10      # Maximum open connections
    max_idle_conns: 5       # Maximum idle connections
    conn_max_idle_time: 3600s  # Connection max idle time
```

### Deploy PostgreSQL with Docker compose

Add PostgreSQL to your `docker-compose.yml`:

```yaml
services:
  headscale:
    image: headscale/headscale:latest
    container_name: headscale
    restart: unless-stopped
    depends_on:
      - postgres  # Wait for PostgreSQL to start
    volumes:
      - ./config:/etc/headscale
      - ./data:/var/lib/headscale
    networks:
      - mesh
    # ... rest of configuration

  postgres:
    image: postgres:15-alpine
    container_name: postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: headscale
      POSTGRES_USER: headscale
      POSTGRES_PASSWORD: your-secure-password
      POSTGRES_INITDB_ARGS: "-E UTF8"
    volumes:
      - ./postgres-data:/var/lib/postgresql/data
    networks:
      - mesh
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U headscale"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  mesh:
    driver: bridge
```

### Secure PostgreSQL password

Use Docker secrets or environment files instead of hardcoding passwords:

```bash
# Create .env file (add to .gitignore!)
cat > .env << 'EOF'
POSTGRES_PASSWORD=your-secure-password-here
EOF

# Update docker-compose.yml to use environment variable
```

```yaml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
```

### Initialize database

Start PostgreSQL and let Headscale initialize the schema:

```bash
# Start PostgreSQL first
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
docker-compose logs -f postgres
# Wait for: "database system is ready to accept connections"

# Start Headscale (will auto-create schema)
docker-compose up -d headscale

# Check Headscale logs
docker logs headscale
# Should see: "Database migrated successfully"
```

### Migrate from SQLite to PostgreSQL

If you have an existing SQLite database:

```bash
# 1. Backup SQLite database
cp data/db.sqlite data/db.sqlite.backup

# 2. Export data from SQLite
docker exec headscale headscale export > headscale-export.json

# 3. Update config.yaml to use PostgreSQL

# 4. Start PostgreSQL
docker-compose up -d postgres

# 5. Import data to PostgreSQL
docker exec -i headscale headscale import < headscale-export.json

# 6. Verify migration
docker exec headscale headscale nodes list
docker exec headscale headscale users list
```

!!! warning "Migration testing"
    Always test migration in a development environment first. Keep SQLite backup until PostgreSQL is verified working.

### PostgreSQL backup and recovery

Regular backups are critical for production deployments:

#### Automated backups

```bash
# Create backup script
cat > backup-postgres.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/headscale_$DATE.sql.gz"

mkdir -p "$BACKUP_DIR"

# Backup database
docker exec postgres pg_dump -U headscale headscale | gzip > "$BACKUP_FILE"

# Keep only last 30 days of backups
find "$BACKUP_DIR" -name "headscale_*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE"
EOF

chmod +x backup-postgres.sh

# Add to crontab (daily at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * /path/to/backup-postgres.sh") | crontab -
```

#### Restore from backup

```bash
# Stop Headscale
docker-compose stop headscale

# Restore database
gunzip -c /backups/postgres/headscale_20250101_020000.sql.gz | \
  docker exec -i postgres psql -U headscale headscale

# Start Headscale
docker-compose start headscale
```

### PostgreSQL performance tuning

For large deployments, tune PostgreSQL settings:

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:15-alpine
    command:
      - "postgres"
      - "-c"
      - "max_connections=200"
      - "-c"
      - "shared_buffers=256MB"
      - "-c"
      - "effective_cache_size=1GB"
      - "-c"
      - "maintenance_work_mem=64MB"
      - "-c"
      - "checkpoint_completion_target=0.9"
      - "-c"
      - "wal_buffers=16MB"
      - "-c"
      - "default_statistics_target=100"
      - "-c"
      - "random_page_cost=1.1"
      - "-c"
      - "effective_io_concurrency=200"
      - "-c"
      - "work_mem=2MB"
      - "-c"
      - "min_wal_size=1GB"
      - "-c"
      - "max_wal_size=4GB"
```

**Tuning parameters explained:**

- `max_connections` - Maximum concurrent connections (adjust based on node count)
- `shared_buffers` - Memory for caching data (25% of RAM recommended)
- `effective_cache_size` - Estimate of OS cache (50-75% of RAM)
- `work_mem` - Memory per query operation (adjust based on query complexity)

### PostgreSQL monitoring

Monitor database health and performance:

```bash
# Check database size
docker exec postgres psql -U headscale -c "SELECT pg_size_pretty(pg_database_size('headscale'));"

# Check active connections
docker exec postgres psql -U headscale -c "SELECT count(*) FROM pg_stat_activity;"

# Check slow queries
docker exec postgres psql -U headscale -c "SELECT query, calls, total_time, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"

# Check table sizes
docker exec postgres psql -U headscale -c "SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

### PostgreSQL high availability

For mission-critical deployments, implement PostgreSQL replication:

**Options:**

- **Streaming replication** - Built-in PostgreSQL replication
- **Patroni** - Automated failover and HA management
- **pgpool-II** - Connection pooling and load balancing
- **Managed services** - AWS RDS, Google Cloud SQL, Azure Database

!!! info "High availability"
    PostgreSQL HA setup is beyond the scope of this guide. Consult PostgreSQL documentation or use managed database services for production HA requirements.

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
   docker exec headscale cat /etc/headscale/derp.yaml
   ```

4. **Check Headscale logs:**

   ```bash
   docker logs headscale | grep -i derp
   ```

### PostgreSQL connection errors

**Symptoms:**

- Headscale can't connect to PostgreSQL
- "connection refused" errors
- Database migration fails

**Solutions:**

1. **Verify PostgreSQL is running:**

   ```bash
   docker-compose ps postgres
   docker logs postgres
   ```

2. **Check network connectivity:**

   ```bash
   # From Headscale container
   docker exec headscale ping postgres
   ```

3. **Verify credentials:**

   ```bash
   # Test connection manually
   docker exec postgres psql -U headscale -d headscale -c "SELECT 1;"
   ```

4. **Check Headscale configuration:**

   ```bash
   docker exec headscale cat /etc/headscale/config.yaml | grep -A 10 database
   ```

### Performance issues with PostgreSQL

**Symptoms:**

- Slow query responses
- High CPU usage on PostgreSQL
- Connection timeouts

**Solutions:**

1. **Check database size:**

   ```bash
   docker exec postgres psql -U headscale -c "SELECT pg_size_pretty(pg_database_size('headscale'));"
   ```

2. **Analyze slow queries:**

   ```bash
   # Enable query logging in docker-compose.yml
   command:
     - "postgres"
     - "-c"
     - "log_min_duration_statement=1000"  # Log queries > 1 second
   ```

3. **Vacuum database:**

   ```bash
   docker exec postgres psql -U headscale -c "VACUUM ANALYZE;"
   ```

4. **Check connection pool settings:**

   ```yaml
   # config.yaml
   database:
     postgres:
       max_open_conns: 10
       max_idle_conns: 5
   ```

## Next steps

- **[Exit node and PCAP](exit-node-pcap.md)** - Configure exit node and packet capture
- **[AmneziaWG](amneziawg.md)** - Configure obfuscation for censorship resistance
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions

---

← [Back to Advanced](index.md) | [Next: Exit node and PCAP](exit-node-pcap.md) →
