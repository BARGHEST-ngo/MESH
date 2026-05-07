# Troubleshooting

Common issues and solutions for MESH deployments.

## Connection issues

### Can't Connect to Control plane

**Symptoms:**

- `meshcli up` fails with connection error
- "Unable to reach control plane" message
- Timeout errors

**Solutions:**

1. **Verify control plane URL**

   ```bash
   curl https://mesh.yourdomain.com/health
   # Should return: {"status":"ok"}
   ```

2. **Check DNS resolution**

   ```bash
   nslookup mesh.yourdomain.com
   dig mesh.yourdomain.com
   ```

3. **Check firewall**

   ```bash
   # Test HTTPS connectivity
   telnet mesh.yourdomain.com 443
   
   # Check if port is open
   nc -zv mesh.yourdomain.com 443
   ```

4. **Verify control plane is running**

   ```bash
   docker compose ps | grep headscale
   docker compose logs headscale
   ```

5. **Check reverse proxy**

   ```bash
   # Nginx
   sudo nginx -t
   sudo systemctl status nginx
   
   # Caddy
   sudo systemctl status caddy
   ```

6. **Try different network**
   - Switch from WiFi to mobile data (or vice versa)
   - Some networks block VPN traffic

### Authentication Fails

**Symptoms:**

- "Invalid auth key" error
- "Auth key expired" error
- Login prompt doesn't appear

**Solutions:**

1. **Verify pre-auth key**

   ```bash
   # On control plane
   docker compose exec headscale headscale preauthkeys list
   ```

2. **Generate new key**

   Pre-auth keys allow nodes to join the mesh without interactive authentication. See the [Create a Pre-authentication key](../setup/control-plane.md#step-6-create-a-pre-authentication-key) instructions for details.

3. **Check key expiration**
   - Pre-auth keys expire after the specified time
   - Generate a new key if expired

### No Mesh IP Assigned

**Symptoms:**

- `meshcli status` doesn't show a mesh IP address
- Node appears offline
- Can't communicate with other nodes

**Solutions:**

1. **Verify pre-auth key is valid**

   ```bash
   # On control plane, check if key is valid and not expired
   docker compose exec headscale headscale preauthkeys list
   ```

2. **Check control plane logs**

   ```bash
   docker compose logs headscale
   # Look for registration errors
   ```

3. **Ensure MESH network exists**

   ```bash
   # List networks
   docker compose exec headscale headscale users list
   ```

If the network does not exist, it can be created by following the [Create a new network](../setup/control-plane.md#step-5-create-a-new-network) instructions.

4. **Try reconnecting**

   ```bash
   docker compose restart analyst
   ```

5. **Check node registration**

   ```bash
   # On control plane
   docker compose exec headscale headscale nodes list
   # Verify your node appears in the list
   ```

## Connectivity issues

### No Peer-to-Peer Connection (Using DERP Relay)

**Symptoms:**

- `meshcli status` shows DERP address in `CurAddr`
- Slow connection speeds
- High latency

**Solutions:**

1. **Check NAT type**

   ```bash
   meshcli netcheck
   ```

2. **Check UDP connectivity**

   ```bash
   # Test if UDP is blocked
   nc -u -v mesh.yourdomain.com 3478
   ```

3. **Check firewall rules**

   ```bash
   # Linux
   sudo iptables -L -n -v
   sudo ufw status
   
   # Allow WireGuard port (usually random high port)
   sudo ufw allow 41641/udp
   ```

4. **Enable UPnP/NAT-PMP on router**
   - Check router settings
   - Enable UPnP or NAT-PMP for automatic port forwarding

5. **Try different network**
   - Some networks (corporate, public WiFi) block UDP
   - Try mobile data or different WiFi network

6. **Force direct connection**

   ```bash
   export TS_DEBUG_ALWAYS_USE_DERP=false
   sudo -E meshcli up --login-server=https://mesh.yourdomain.com
   ```

### Peers Not Visible

**Symptoms:**

- `meshcli status` shows no peers
- Devices don't appear in mesh
- Can't ping peers

**Solutions:**

1. **Verify both devices are connected**

   ```bash
   meshcli status
   # Should show "Running"
   ```

2. **Check ACLs**

   ```bash
   # On control plane
   docker compose exec headscale headscale nodes list

   # Print the active policy
   docker compose exec headscale headscale policy get
   ```

   The policy can also be inspected and edited from the CUSTOMIZE ACL page in the control plane web UI.

3. **Verify same control plane**
   - Ensure all devices connect to the same control plane URL

4. **Check user assignment**

   ```bash
   docker compose exec headscale headscale nodes list
   # Verify nodes are in correct user/group
   ```

5. **Restart daemon**

   ```bash
   docker compose restart analyst
   ```

### Connection Drops Frequently

**Symptoms:**

- Intermittent connectivity
- Frequent reconnections
- "Peer offline" messages

**Solutions:**

1. **Check network stability**

   ```bash
   ping -c 100 8.8.8.8
   # Look for packet loss
   ```

2. **Disable power saving (mobile)**
   - Android: Settings > Battery > Optimize battery usage > MESH > Don't optimize
   - Disable "Adaptive battery" for MESH app

3. **Increase keepalive**

   ```bash
   # Edit WireGuard config to add:
   PersistentKeepalive = 25
   ```

4. **Check DERP relay stability**

   ```bash
   meshcli netcheck
   # Check DERP latencies
   ```

5. **Monitor logs**

   ```bash
   docker compose logs -f analyst
   # Look for errors or warnings
   ```

## ADB Issues

### ADB Connection Fails

**Symptoms:**

- `meshcli adbpair` fails
- "Connection refused" error
- "Unable to connect" error

**Solutions:**

1. **Verify device is online**

   ```bash
   meshcli status
   ping 100.64.X.X
   ```

2. **Check ADB port**

   ```bash
   # On Android device (via USB)
   adb shell getprop service.adb.tcp.port
   # Should return a valid port number
   ```

3. **Enable ADB-over-WiFi**

   ```bash
   # Via USB
   adb tcpip 5555
   adb connect 100.64.X.X:5555
   ```

4. **Restart ADB**

   ```bash
   adb kill-server
   adb start-server
   adb connect 100.64.X.X:5555
   ```

5. **Check firewall on Android**
   - Some security apps block ADB
   - Temporarily disable and test

6. **Verify ADB is enabled in MESH app**
   - Open MESH app
   - Settings > ADB-over-WiFi > Enabled

### ADB Disconnects Frequently

**Symptoms:**

- ADB connection drops during operations
- "device offline" errors
- Need to reconnect frequently

**Solutions:**

1. **Disable power saving**
   - Settings > Battery > MESH > Don't optimize

2. **Keep MESH app in foreground**
   - Don't minimize the app during ADB operations

3. **Use persistent connection**

   ```bash
   # Keep pinging to maintain connection
   ping 100.64.X.X &
   adb connect 100.64.X.X:5555
   ```

4. **Check network stability**
   - Ensure stable WiFi/mobile connection
   - Avoid switching networks during operation

## Performance issues

### Slow Transfer Speeds

**Symptoms:**

- Slow file transfers
- High latency
- Timeouts during large operations

**Solutions:**

1. **Check connection type**

   ```bash
   meshcli status
   # Direct connection is faster than DERP relay
   ```

2. **Disable AmneziaWG obfuscation**

   ```bash
   # If not needed for censorship resistance
   docker compose exec analyst rm -f /etc/mesh/amneziawg.conf
   docker compose restart analyst
   ```

3. **Use lighter obfuscation**
   - Switch from Heavy to Balanced or Light preset
   - See [AmneziaWG Configuration](../advanced/amneziawg.md)

4. **Check network bandwidth**

   ```bash
   # Test bandwidth
   iperf3 -s  # On one peer
   iperf3 -c 100.64.X.X  # On other peer
   ```

5. **Compress data**

   ```bash
   # Compress before transfer
   adb shell tar -czf /sdcard/data.tar.gz /data/local/tmp/
   adb pull /sdcard/data.tar.gz
   ```

### High CPU Usage

**Symptoms:**

- High CPU usage by mesh
- Device heating up
- Battery drain (mobile)

**Solutions:**

1. **Disable obfuscation**
   - AmneziaWG adds CPU overhead
   - Disable if not needed

2. **Limit peer count**
   - Use ACLs to limit visible peers
   - Only connect to necessary devices

## Control plane Issues

### Headscale Won't Start

**Symptoms:**

- Docker container exits immediately
- "Address already in use" errors
- Database errors

**Solutions:**

1. **Check logs**

   ```bash
   docker compose logs headscale
   ```

2. **Check port conflicts**

   ```bash
   sudo lsof -i :8080
   sudo lsof -i :9090
   sudo lsof -i :3478
   ```

3. **Verify config syntax**

   ```bash
   cat control-plane/headscale/config.yaml
   # Check for YAML syntax errors
   ```

4. **Reset database (WARNING: deletes all data)**

   ```bash
   task down
   docker volume rm mesh_headscale-data
   task controlPlane
   ```

### Nodes Not Registering

**Symptoms:**

- Nodes don't appear in `headscale nodes list`
- Authentication succeeds but node not visible
- "Node not found" errors

**Solutions:**

1. **Check node list**

   ```bash
   docker compose exec headscale headscale nodes list
   ```

2. **Check network exists**

   ```bash
   docker compose exec headscale headscale users list
   ```

3. **Create network if missing**

   If the network does not exist, it can be created by following the [Create a new network](../setup/control-plane.md#step-5-create-a-new-network) instructions.

4. **Check ACLs**

   ```bash
   docker compose exec headscale headscale policy get
   # Verify ACLs allow the connection
   ```

### DERP Relay Not Working

**Symptoms:**

- Peers can't connect even via relay
- "No DERP home" errors
- All connections fail

**Solutions:**

1. **Check STUN port**

   ```bash
   sudo netstat -ulnp | grep 3478
   docker compose logs headscale | grep STUN
   ```

2. **Test STUN**

   ```bash
   sudo apt install stun-client
   stun mesh.yourdomain.com 3478
   ```

3. **Check DERP config**

   ```bash
   cat control-plane/headscale/config.yaml | grep -A 20 "derp:"
   ```

4. **Verify DERP is enabled**

   ```yaml
   derp:
     server:
       enabled: true
       stun_listen_addr: "0.0.0.0:3478"
   ```

5. **Check firewall**

   ```bash
   sudo ufw allow 3478/udp
   ```

## Android App Issues

### App Crashes

**Symptoms:**

- App closes unexpectedly
- "MESH has stopped" error
- Can't open app

**Solutions:**

1. **Clear app data**
   - Settings > Apps > MESH > Storage > Clear Data
   - Reconnect to mesh

2. **Reinstall app**

   ```bash
   adb uninstall com.barghest.mesh
   adb install mesh-android.apk
   ```

3. **Check logs**

   ```bash
   adb logcat | grep MESH
   adb logcat | grep AndroidRuntime
   ```

4. **Check Android version**
   - MESH requires Android 8.0+
   - Update Android if possible

### VPN Permission Denied

**Symptoms:**

- "VPN permission required" error
- Can't connect to mesh
- Permission dialogue doesn't appear

**Solutions:**

1. **Grant permission manually**
   - Settings > Apps > MESH > Permissions
   - Enable VPN permission

2. **Check for conflicting VPN**
   - Disconnect other VPN apps
   - Only one VPN can be active at a time

3. **Reinstall app**
   - Uninstall and reinstall to reset permissions

## Diagnostic commands

### Collect Logs

```bash
# Analyst client logs
docker compose logs --since 1h analyst > mesh-logs.txt

# Control plane logs
docker compose logs --since 1h headscale > headscale-logs.txt

# Android logs
adb logcat -d > android-logs.txt
```

### Generate Bug Report

```bash
# Analyst client
meshcli bugreport --diagnose > mesh-bugreport.txt

# Android
adb bugreport bugreport.zip
```

### Network Diagnostics

```bash
# Check NAT type and DERP latencies
meshcli netcheck

# Trace route to peer
meshcli ping 100.64.X.X --verbose

# Check WireGuard interface
ip addr show tun0
sudo wg show
```

### Connection Diagnostics

```bash
# Detailed status
meshcli status --json

# Check peer connectivity
for peer in $(meshcli status --json | jq -r '.Peer[].TailscaleIPs[0]'); do
    echo "Pinging $peer..."
    meshcli ping $peer --c 3
done

# Monitor connection changes
watch -n 1 'meshcli status --peers'
```

## Getting Help

If you can't resolve your issue:

1. **Check GitHub Issues**: [github.com/BARGHEST-ngo/mesh/issues](https://github.com/BARGHEST-ngo/mesh/issues)
2. **Create Bug Report**: Include logs, config, and steps to reproduce
3. **Contact BARGHEST**: [barghest.asia](https://barghest.asia)

## Next steps

- **[User guide](../user-guide/user-guide.md)** - Learn forensic workflows
- **[CLI reference](cli-reference.md)** - Complete command documentation
- **[Architecture](../overview/architecture.md)** - Understand how MESH works
