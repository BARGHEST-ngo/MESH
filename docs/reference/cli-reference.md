# CLI reference

Complete reference for the MESH command-line interface (`meshcli`).

## Synopsis

```bash
meshcli [global-options] <command> [command-options]
```

## Global Options

- `--socket <path>`: Path to tailscaled socket (default: `/var/run/mesh/tailscaled.sock`)
- `--help`, `-h`: Show help message
- `--version`, `-v`: Show version information

## Commands

### up

Connect to the mesh network.

**Synopsis:**

```bash
meshcli up [options]
```

**Options:**

- `--login-server <url>`: Control plane URL (required)
- `--authkey <key>`: Pre-authentication key
- `--accept-dns <bool>`: Accept DNS configuration (default: true)
- `--accept-routes <bool>`: Accept subnet routes from peers (default: false)
- `--advertise-routes <routes>`: Advertise subnet routes (comma-separated CIDRs)
- `--advertise-exit-node`: Advertise as exit node
- `--exit-node <ip>`: Use peer as exit node
- `--exit-node-allow-lan-access`: Allow LAN access when using exit node
- `--hostname <name>`: Set hostname for this node
- `--shields-up`: Block all incoming connections
- `--snat-subnet-routes`: SNAT traffic to advertised routes
- `--netfilter-mode <mode>`: Netfilter mode (on, off, nodivert)
- `--operator <user>`: Unix user to run as
- `--ssh`: Enable SSH server
- `--timeout <duration>`: Timeout for connection attempt

**Examples:**

```bash
# Basic connection
sudo meshcli up --login-server=https://mesh.example.com

# Connect with pre-auth key
sudo meshcli up \
  --login-server=https://mesh.example.com \
  --authkey=abc123def456

# Advertise subnet routes
sudo meshcli up \
  --login-server=https://mesh.example.com \
  --advertise-routes=192.168.1.0/24,10.0.0.0/8

# Act as exit node
sudo meshcli up \
  --login-server=https://mesh.example.com \
  --advertise-exit-node

# Use peer as exit node
sudo meshcli up \
  --login-server=https://mesh.example.com \
  --exit-node=100.64.1.5

# Disable DNS
sudo meshcli up \
  --login-server=https://mesh.example.com \
  --accept-dns=false
```

### down

Disconnect from the mesh network.

**Synopsis:**

```bash
meshcli down
```

**Examples:**

```bash
sudo meshcli down
```

### status

Show connection status and peer information.

**Synopsis:**

```bash
meshcli status [options]
```

**Options:**

- `--peers`: Show detailed peer information
- `--json`: Output in JSON format
- `--self <bool>`: Show self information (default: true)
- `--active`: Only show active peers

**Examples:**

```bash
# Basic status
sudo meshcli status

# Show peers
sudo meshcli status --peers

# JSON output
sudo meshcli status --json

# JSON with jq filtering
sudo meshcli status --json | jq '.Peer[] | {name: .HostName, ip: .TailscaleIPs[0]}'

# Only active peers
sudo meshcli status --peers --active
```

**Output fields:**

- `BackendState`: Connection state (Running, Stopped, etc.)
- `TailscaleIPs`: Mesh IP addresses
- `Self`: Information about this node
- `Peer`: Information about connected peers
  - `HostName`: Peer hostname
  - `TailscaleIPs`: Peer mesh IPs
  - `Online`: Whether peer is online
  - `CurAddr`: Current connection address (direct or DERP)
  - `RxBytes`, `TxBytes`: Bytes received/transmitted
  - `Created`: When peer was added
  - `LastSeen`: Last activity timestamp

### ping

Ping a peer using MESH protocol.

**Synopsis:**

```bash
meshcli ping <hostname-or-ip> [options]
```

**Options:**

- `--timeout <duration>`: Timeout for ping (default: 5s)
- `--c <count>`: Number of pings to send
- `--until-direct`: Ping until direct connection established
- `--verbose`: Verbose output

**Examples:**

```bash
# Ping by hostname
sudo meshcli ping android-device

# Ping by IP
sudo meshcli ping 100.64.2.1

# Ping 10 times
sudo meshcli ping 100.64.2.1 --c 10

# Ping until direct connection
sudo meshcli ping 100.64.2.1 --until-direct
```

### netcheck

Check network connectivity and NAT traversal capabilities.

**Synopsis:**

```bash
meshcli netcheck
```

**Examples:**

```bash
sudo meshcli netcheck
```

**Output:**

- UDP connectivity status
- DERP relay latencies
- NAT type
- Port mapping capabilities
- Preferred DERP region

### bugreport

Generate a bug report for troubleshooting.

**Synopsis:**

```bash
meshcli bugreport [options]
```

**Options:**

- `--diagnose`: Include diagnostic information

**Examples:**

```bash
sudo meshcli bugreport > mesh-bugreport.txt
sudo meshcli bugreport --diagnose > mesh-bugreport-full.txt
```

### logout

Log out from the control plane.

**Synopsis:**

```bash
meshcli logout
```

**Examples:**

```bash
sudo meshcli logout
```

### set

Configure MESH settings.

**Synopsis:**

```bash
meshcli set [options]
```

**Options:**

- `--accept-dns <bool>`: Accept DNS configuration
- `--accept-routes <bool>`: Accept subnet routes
- `--advertise-routes <routes>`: Advertise subnet routes
- `--advertise-exit-node <bool>`: Advertise as exit node
- `--exit-node <ip>`: Use peer as exit node
- `--shields-up <bool>`: Block incoming connections
- `--hostname <name>`: Set hostname

**Examples:**

```bash
# Disable DNS
sudo meshcli set --accept-dns=false

# Enable route acceptance
sudo meshcli set --accept-routes=true

# Change exit node
sudo meshcli set --exit-node=100.64.1.5

# Remove exit node
sudo meshcli set --exit-node=

# Enable shields up
sudo meshcli set --shields-up=true
```

### ip

Show mesh IP addresses.

**Synopsis:**

```bash
meshcli ip [options]
```

**Options:**

- `-4`: Show only IPv4
- `-6`: Show only IPv6
- `-1`: Show only first IP

**Examples:**

```bash
# Show all IPs
sudo meshcli ip

# IPv4 only
sudo meshcli ip -4

# First IP only
sudo meshcli ip -1
```

### routes

Show and manage routes.

**Synopsis:**

```bash
meshcli routes
```

**Examples:**

```bash
sudo meshcli routes
```

### version

Show version information.

**Synopsis:**

```bash
meshcli version
```

**Examples:**

```bash
meshcli version
```

## Environment Variables

- `TS_DEBUG_TRIM_WIREGUARD`: Prevent peer trimming (set to `false`)
- `TS_DEBUG_ALWAYS_USE_DERP`: Force DERP relay usage (set to `true`)
- `TS_DEBUG_FIREWALL_MODE`: Firewall mode (auto, on, off)

**Examples:**

```bash
# Prevent peer trimming
export TS_DEBUG_TRIM_WIREGUARD=false
sudo -E meshcli up --login-server=https://mesh.example.com

# Force DERP relay
export TS_DEBUG_ALWAYS_USE_DERP=true
sudo -E meshcli up --login-server=https://mesh.example.com
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Connection error
- `3`: Authentication error

## Configuration Files

### State File

Location: `/var/lib/mesh/tailscaled.state`

Contains:

- Node private key
- Control plane URL
- Authentication state

**Backup:**

```bash
sudo cp /var/lib/mesh/tailscaled.state /var/lib/mesh/tailscaled.state.backup
```

### AmneziaWG Config

Location: `/etc/mesh/amneziawg.conf`

See [AmneziaWG Configuration](../advanced/amneziawg.md) for details.

## Scripting Examples

### Check if connected

```bash
#!/bin/bash
if sudo meshcli status | grep -q "Running"; then
    echo "Connected"
else
    echo "Disconnected"
fi
```

### Get mesh IP

```bash
#!/bin/bash
MESH_IP=$(sudo meshcli ip -4 -1)
echo "Mesh IP: $MESH_IP"
```

### List online peers

```bash
#!/bin/bash
sudo meshcli status --json | jq -r '.Peer[] | select(.Online == true) | .HostName'
```

### Auto-reconnect on disconnect

```bash
#!/bin/bash
while true; do
    if ! sudo meshcli status | grep -q "Running"; then
        echo "Disconnected, reconnecting..."
        sudo meshcli up --login-server=https://mesh.example.com --authkey=abc123
    fi
    sleep 60
done
```

### Monitor connection type

```bash
#!/bin/bash
sudo meshcli status --json | jq -r '.Peer[] | "\(.HostName): \(.CurAddr)"'
```

## Next steps

- **[User guide](../user-guide/user-guide.md)** - Learn forensic workflows
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Analyst client](../installation/analyst-client.md)** - Analyst client documentation
