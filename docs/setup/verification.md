# Verification and Testing

Now that you have all components set up, let's verify everything is working correctly and run your first forensic collection.

## Step 1: Check Mesh Status

On your analyst workstation, verify all nodes are connected:

```bash
sudo ./meshcli status --peers
```

**Example output:**

```
# Peers:
analyst-laptop    100.64.1.1    online    user:analyst1
android-device    100.64.2.1    online    user:analyst1
```

You should see:

- Your analyst workstation (e.g., `analyst-laptop`)
- Your Android device (e.g., `android-device`)
- Both showing as **online**
- Both assigned mesh IP addresses

!!! success "All Nodes Online"
    If you see both nodes listed as "online", congratulations! Your mesh network is operational.

## Step 2: Test Connectivity

### Ping the Android Device

Test basic network connectivity to the Android device:

```bash
ping 100.64.2.1
```

Replace `100.64.2.1` with your Android device's actual mesh IP address.

**Example output:**

```
PING 100.64.2.1 (100.64.2.1) 56(84) bytes of data.
64 bytes from 100.64.2.1: icmp_seq=1 ttl=64 time=45.2 ms
64 bytes from 100.64.2.1: icmp_seq=2 ttl=64 time=42.8 ms
```

!!! tip "High Latency?"
    If you see high latency (>100ms), the connection may be using the DERP relay instead of a direct peer-to-peer connection. This is normal if both devices are behind restrictive NAT.

### Check Connection Type

See if you have a direct peer-to-peer connection or are using the DERP relay:

```bash
sudo ./meshcli status --peers --json | grep -A 5 "android-device"
```

Look for the `relay` field:

- If empty or shows a direct IP: You have a P2P connection ✓
- If shows a DERP server: You're using the relay (still works, just higher latency)

## Step 3: Connect via ADB

The Android device should automatically enable ADB-over-WiFi on its mesh IP. Let's connect to it.

### Connect to the Device

```bash
adb connect 100.64.2.1:5555
```

Replace `100.64.2.1` with your Android device's mesh IP.

**Example output:**

```
connected to 100.64.2.1:5555
```

### Verify ADB Connection

```bash
adb devices
```

**Example output:**

```
List of devices attached
100.64.2.1:5555    device
```

!!! success "ADB Connected"
    If you see your device listed, you now have remote ADB access over the mesh!

### Run a Test Command

```bash
# Get device model
adb shell getprop ro.product.model

# Get Android version
adb shell getprop ro.build.version.release

# List installed packages
adb shell pm list packages | head -10
```

## Step 4: Run a Forensic Collection

Now let's run some basic forensic collection commands to verify everything works.

### Collect a Bug Report

```bash
# Collect a full bug report
adb bugreport bugreport.zip
```

This creates a comprehensive bug report containing system logs, diagnostics, and configuration.

!!! info "Bug Report Time"
    Bug reports can take 1-2 minutes to generate. Wait for the command to complete.

### Get System Information

```bash
# Dump all system services
adb shell dumpsys > dumpsys.txt

# Get running processes
adb shell ps > processes.txt

# Get installed packages
adb shell pm list packages -f > packages.txt
```

### Collect Logs

```bash
# Get system logs
adb logcat -d > logcat.txt

# Get kernel logs
adb shell dmesg > dmesg.txt
```

### Pull Files from Device

```bash
# Create output directory
mkdir -p ./artifacts

# Pull system build info
adb pull /system/build.prop ./artifacts/

# Pull package list
adb shell pm list packages -f > ./artifacts/packages.txt
```

## Step 5: Run AndroidQF (Optional)

If you have AndroidQF installed, you can run automated spyware detection:

```bash
# Run AndroidQF against the mesh-connected device
androidqf --adb 100.64.2.1:5555 --output ./artifacts/
```

AndroidQF will:

- Collect installed packages
- Check for known spyware indicators
- Extract system information
- Generate a report

!!! tip "Installing AndroidQF"
    If you don't have AndroidQF installed, see the [AndroidQF documentation](https://github.com/mvt-project/androidqf) for installation instructions.

## Verification Checklist

Confirm all of the following are working:

- [ ] Control plane is running (`docker ps` shows headscale containers)
- [ ] Analyst client shows as "online" in mesh
- [ ] Android device shows as "online" in mesh
- [ ] Can ping Android device over mesh IP
- [ ] Can connect via ADB over mesh IP
- [ ] Can run ADB commands successfully
- [ ] Can collect forensic artifacts

!!! success "Setup Complete!"
    If all items are checked, your MESH network is fully operational and ready for forensic investigations!

## Troubleshooting

### Can't Ping Android Device

**Check both nodes are online:**

```bash
sudo ./meshcli status --peers
```

**Check firewall rules:**

Some networks block ICMP (ping). Try ADB connection instead - it uses TCP which is more likely to work.

**Check DERP relay:**

If direct P2P fails, MESH should fall back to DERP relay. Check control plane logs:

```bash
docker logs headscale | grep DERP
```

### ADB Connection Fails

**Error: "Connection refused"**

1. Verify ADB is enabled in the MESH app on Android
2. Check the mesh IP address is correct
3. Ensure port 5555 is not blocked

**Error: "Connection timed out"**

1. Verify the Android device is online in the mesh
2. Try pinging the device first
3. Check if the device's firewall is blocking port 5555

**Fix: Restart ADB**

```bash
adb kill-server
adb start-server
adb connect 100.64.2.1:5555
```

### No Devices Showing in Mesh

**On analyst workstation:**

```bash
# Check daemon is running
ps aux | grep tailscaled-amnezia

# Check connection status
sudo ./meshcli status

# Try reconnecting
sudo ./meshcli down
sudo ./meshcli up --login-server=https://your-domain.com --authkey=YOUR_KEY
```

**On control plane:**

```bash
# Check registered nodes
docker exec headscale headscale nodes list

# Check logs
docker logs headscale
```

## Next steps

Congratulations! You now have a fully functional MESH network for mobile forensics. Here's what to explore next:

- **[User guide](../user-guide/user-guide.md)** - Learn forensic workflows and best practices
- **[AmneziaWG Configuration](../advanced/amneziawg.md)** - Enable censorship resistance
- **[CLI reference](../reference/cli-reference.md)** - Complete command documentation
- **[Architecture](../overview/architecture.md)** - Understand how MESH works under the hood

---

← [Previous: Endpoint client Setup](endpoint-client.md) | [Next: Next steps](next-steps.md) →
