# User guide

This guide covers common forensic workflows and best practices for using MESH in real-world investigations.

## Forensic workflows

### Basic Android mobile forensics workflow

1. **Set up the mesh network**
   - Deploy control plane
   - Install analyst client
   - Install endpoint client on target device

2. **Connect devices**
   - Start MESH daemon on analyst workstation
   - Connect endpoint device to mesh
   - Verify connectivity with `meshcli status --peers`

3. **Establish ADB connection**

   ```bash
   adb connect 100.64.X.X:5555
   adb devices
   ```

4. **Collect artifacts & run analysis**

    You should refer to for a comprehensive guide on Android ADB forensics: <https://forensics.socialtic.org/how-tos/04-how-to-extract-with-androidqf/index.html>

5. **Tear down**

   ```bash
   adb disconnect
   sudo meshcli down
   ```

### Censored Network Workflow

For operations in censored environments:

1. **Enable AmneziaWG obfuscation**

   **Analyst:**

   ```bash
   sudo cat > /etc/mesh/amneziawg.conf << 'EOF'
   [Interface]
   Jc = 10
   Jmin = 50
   Jmax = 1000
   S1 = 100
   S2 = 150
   H1 = 1234567
   H2 = 2345678
   H3 = 3456789
   H4 = 4567890
   EOF
   sudo systemctl restart mesh
   ```

   **Endpoint:**
   - Open MESH app
   - Settings > AmneziaWG > Heavy
   - Reconnect

2. **Verify DERP fallback**

   ```bash
   sudo meshcli status --json | jq '.Peer[] | {name: .HostName, relay: .CurAddr}'
   ```

3. **Proceed with investigation**
   - Connection may be slower due to obfuscation/relay
   - Monitor connection stability
   - Have backup communication channel

### Team-based investigation

Multiple analysts investigating multiple devices:

1. **Set up ACLs**

   ```yaml
   groups:
     group:analysts:
       - analyst1
       - analyst2
     group:endpoints:
       - device1
       - device2
   
   acls:
     - action: accept
       src:
         - group:analysts
       dst:
         - group:endpoints:*
   ```

2. **Assign devices**
   - Use tags to organise devices
   - Document which analyst is investigating which device

3. **Coordinate**
   - Use separate communication channel for coordination
   - Share findings through secure channels
   - Avoid simultaneous ADB connections to same device

## Advanced Features

### Subnet Routing

Access devices on the endpoint's local network:

**On endpoint (Android):**

1. Open MESH app
2. Settings > Advertise LAN routes
3. Enable and select network interface

**On analyst:**

```bash
sudo meshcli up --accept-routes
```

**Access LAN devices:**

```bash
# Ping device on endpoint's LAN
ping 192.168.1.100

# SSH to device on endpoint's LAN
ssh user@192.168.1.100
```

### Kill switch

Block all non-MESH traffic on endpoint:

**On endpoint (Android):**

1. Settings > VPN > MESH Settings > Always-on VPN -> Block connections without VPN
2. Enable

**Use cases:**

- Prevent data leaks during investigation
- Isolate compromised device
- Ensure all traffic goes through mesh

### Exit nodes for Network Capture

Route internet traffic through an endpoint:

**On analyst:**

```bash
# Use endpoint as exit node
sudo meshcli up --exit-node=100.64.X.X

# Verify
curl ifconfig.me  # Should show endpoint's public IP

# Stop using exit node
sudo meshcli up --exit-node=
```

**On endpoint:**

1. Open MESH app
2. Settings > Use 100.64.x.x as exit node
3. Enable

**Use cases:**

- Monitor endpoint's internet traffic and capture traffic
- Access geo-restricted content
- Investigate network-level issues

### Operational

1. **Test before deployment**: Verify setup in controlled environment
2. **Document procedures**: Create runbooks for common tasks
3. **Have backup plans**: Alternative communication channels
4. **Coordinate with team**: Clear communication and task assignment
5. **Monitor performance**: Watch for connection issues
6. **Plan for disconnections**: Handle network interruptions gracefully
7. **Respect privacy**: Only collect necessary data
8. **Follow legal requirements**: Ensure compliance with local laws

### Performance

1. **Use P2P when possible**: Direct connections are faster
2. **Minimize obfuscation overhead**: Use lightest AmneziaWG preset that works
3. **Collect artifacts efficiently**: Avoid unnecessary data transfer
4. **Use compression**: Compress large artifacts before transfer
5. **Monitor bandwidth**: Be aware of network limitations
6. **Schedule intensive tasks**: Run large collections during off-peak hours

## Common Tasks

### Checking Connection Status

```bash
# Basic status
sudo meshcli status

# Detailed peer information
sudo meshcli status --peers

# JSON output for scripting
sudo meshcli status --json | jq

# Check if using DERP relay
sudo meshcli status --json | jq '.Peer[] | select(.CurAddr | contains("derp"))'
```

### Collecting comprehensive artifacts

```bash
#!/bin/bash
# comprehensive-collection.sh

DEVICE_IP="100.64.X.X"
OUTPUT_DIR="./artifacts/$(date +%Y%m%d-%H%M%S)"
mkdir -p "$OUTPUT_DIR"

# Connect
adb connect $DEVICE_IP:5555

# Bug report
adb bugreport "$OUTPUT_DIR/bugreport.zip"

# System dump
adb shell dumpsys > "$OUTPUT_DIR/dumpsys.txt"

# Logcat
adb logcat -d > "$OUTPUT_DIR/logcat.txt"

# Package list
adb shell pm list packages -f > "$OUTPUT_DIR/packages.txt"

# Running processes
adb shell ps -A > "$OUTPUT_DIR/processes.txt"

# Network connections
adb shell netstat > "$OUTPUT_DIR/netstat.txt"

# System properties
adb shell getprop > "$OUTPUT_DIR/properties.txt"

# AndroidQF
androidqf --adb $DEVICE_IP:5555 --output "$OUTPUT_DIR/androidqf/"

# Disconnect
adb disconnect

echo "Collection complete: $OUTPUT_DIR"
```

### Managing multiple devices

```bash
# List all connected devices
sudo meshcli status --peers | grep -E "100\.64\."

# Connect to all devices
for ip in $(sudo meshcli status --peers | grep -oE "100\.64\.[0-9]+\.[0-9]+"); do
    adb connect $ip:5555
done

# Check all ADB connections
adb devices

# Disconnect all
adb disconnect
```

### Automating investigations

```bash
#!/bin/bash
# auto-investigate.sh

CONTROL_PLANE="https://mesh.yourdomain.com"
PREAUTH_KEY="your-key-here"

# Start daemon
sudo systemctl start mesh

# Connect to control plane
sudo meshcli up --login-server=$CONTROL_PLANE --authkey=$PREAUTH_KEY --accept-dns=false

# Wait for peers
echo "Waiting for devices..."
while [ $(sudo meshcli status --peers | grep -c "100.64") -eq 0 ]; do
    sleep 5
done

# Get device IPs
DEVICES=$(sudo meshcli status --peers | grep -oE "100\.64\.[0-9]+\.[0-9]+")

# Investigate each device
for device in $DEVICES; do
    echo "Investigating $device..."
    ./comprehensive-collection.sh $device
done

# Disconnect
sudo meshcli down

echo "Investigation complete"
```

## Troubleshooting

See the [Troubleshooting guide](../reference/troubleshooting.md) for detailed solutions to common issues.

## Next steps

- **[CLI reference](../reference/cli-reference.md)** - Complete command documentation
- **[AmneziaWG Configuration](../advanced/amneziawg.md)** - Configure censorship resistance
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
