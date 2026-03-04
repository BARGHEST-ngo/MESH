# Endpoint Clients

MESH endpoint clients run on the target devices being analysed. Currently, Android is supported with iOS support planned for Q4 2026.

## Android Client

### Overview

The MESH Android client is a native Android app that:

- Connects to the MESH network via WireGuard/AmneziaWG
- Automatically enables ADB-over-WiFi on the mesh interface
- Supports ephemeral connections (disconnect on app close)
- Provides censorship resistance via AmneziaWG obfuscation
- Runs as a VPN service for persistent connectivity

### System Requirements

- **Android Version**: 8.0 (Oreo) or later
- **Permissions**: VPN permission (granted on first connection)
- **Storage**: 50MB for app and state
- **Network**: WiFi or mobile data connectivity

### Installation

#### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh/mesh-android-client

# Build the APK
./gradlew assembleRelease

# The APK will be at:
# app/build/outputs/apk/release/app-release.apk
```

#### Option 2: Download Pre-built APK

```bash
# Download the latest release
wget https://github.com/BARGHEST-ngo/mesh/releases/latest/download/mesh-android.apk
```

#### Install on Device

**Via USB:**

```bash
# Enable USB debugging on the Android device
# Settings > Developer Options > USB Debugging

# Connect device via USB
adb devices

# Install the APK
adb install app-release.apk
```

**Via File Transfer:**

1. Copy the APK to the device (email, cloud storage, etc.)
2. Open the APK file on the device
3. Tap "Install" (you may need to enable "Install from Unknown Sources")

### Configuration

#### First-Time Setup

1. **Open the MESH app**
2. **Grant VPN permission** when prompted
3. **Enter control plane URL**: `https://mesh.yourdomain.com`
4. **Enter pre-authentication key** (obtained from control plane)
5. **Tap "Connect"**

The device will now join the mesh network and be assigned a CGNAT IP address (e.g., `100.64.X.X`).

#### Connection Settings

- **Auto-connect**: Automatically connect when app opens
- **Ephemeral mode**: Disconnect when app closes
- **ADB-over-WiFi**: Automatically enabled on mesh interface
- **Kill-switch**: Block non-MESH traffic (optional)

### Using the Android Client

#### Connecting to the Mesh

1. Open the MESH app
2. Tap "Connect" (if not auto-connected)
3. Wait for "Connected" status
4. Note your mesh IP address displayed in the app

#### Viewing Connection Status

The app displays:

- **Connection status**: Connected/Disconnected/Connecting
- **Mesh IP address**: Your CGNAT IP (e.g., `100.64.2.1`)
- **Peers**: Number of connected peers
- **Transfer stats**: Bytes sent/received
- **Connection type**: Direct P2P or DERP relay

#### Disconnecting

1. Open the MESH app
2. Tap "Disconnect"

Or simply close the app if ephemeral mode is enabled.

### ADB-over-WiFi

The MESH Android client automatically enables ADB-over-WiFi on port 5555 of the mesh interface.

#### Connecting from Analyst

```bash
# Find the device's mesh IP
sudo meshcli status --peers

# Connect via ADB
adb connect 100.64.X.X:5555

# Verify connection
adb devices
```

#### Troubleshooting ADB

**ADB not responding:**

```bash
# On the Android device (via USB first)
adb shell setprop service.adb.tcp.port 5555
adb shell stop adbd
adb shell start adbd
```

**Connection refused:**

- Verify the device is connected to the mesh
- Check firewall settings on the device
- Ensure ADB-over-WiFi is enabled in the MESH app settings

### Forensic Collection

Once connected via ADB, you can collect forensic artifacts:

```bash
# Bug report
adb bugreport bugreport.zip

# System dump
adb shell dumpsys > dumpsys.txt

# Logcat
adb logcat -d > logcat.txt

# Package list
adb shell pm list packages -f > packages.txt

# Pull files
adb pull /sdcard/Download/ ./downloads/

# Run AndroidQF
androidqf --adb 100.64.X.X:5555 --output ./artifacts/

# Run MVT
mvt-android check-adb --serial 100.64.X.X:5555 --output ./mvt-results/
```

### Advanced Features

#### Subnet Routing

Advertise the device's local network to the mesh:

1. Open MESH app settings
2. Enable "Advertise LAN routes"
3. Select network interface (WiFi/Mobile)

Analysts can then access devices on the same network as the endpoint.

#### Exit Node

Allow analysts to route traffic through this device:

1. Open MESH app settings
2. Enable "Act as exit node"

Analysts can then use this device as an internet gateway.

#### AmneziaWG Obfuscation

Enable censorship resistance:

1. Open MESH app settings
2. Tap "AmneziaWG Configuration"
3. Select a preset:
    - **Disabled**: Standard WireGuard (default)
    - **Light**: Minimal obfuscation
    - **Balanced**: Recommended for most cases
    - **Heavy**: Maximum DPI evasion
4. Tap "Apply"
5. Reconnect to the mesh

See [AmneziaWG Configuration](../advanced/amneziawg.md) for details.

### Troubleshooting

#### Can't Connect to Control plane

- Verify control plane URL is correct
- Check internet connectivity
- Try disabling WiFi and using mobile data (or vice versa)
- Check control plane logs for errors

#### No Peer-to-Peer Connection

- Verify both devices show as "online" in the mesh
- Check if UDP is blocked (app will fall back to DERP relay)
- Enable AmneziaWG obfuscation if in a censored network

#### VPN Permission Denied

- Go to Settings > Apps > MESH > Permissions
- Grant VPN permission
- Restart the app

#### High Battery Usage

- Disable "Keep alive" in settings
- Use ephemeral mode (disconnect when not in use)
- Reduce connection check frequency

#### App Crashes

- Clear app data: Settings > Apps > MESH > Storage > Clear Data
- Reinstall the app
- Check logcat: `adb logcat | grep MESH`

### Security Considerations

1. **VPN Permission**: MESH requires VPN permission to create the tunnel interface
2. **ADB Access**: ADB-over-WiFi provides full device access - only connect to trusted analysts
3. **Kill-Switch**: Enable to prevent data leaks outside the mesh
4. **Ephemeral Mode**: Use for temporary investigations
5. **Pre-auth Keys**: Use short-lived keys to limit unauthorized access

### Uninstallation

1. Open MESH app
2. Disconnect from the mesh
3. Go to Settings > Apps > MESH
4. Tap "Uninstall"

Or via ADB:

```bash
adb uninstall com.barghest.mesh
```

## iOS Client (Planned)

iOS support is planned for Q4 2026. The iOS client will provide:

- WireGuard/AmneziaWG connectivity
- libimobiledevice support for forensic access
- Similar features to the Android client
- App Store distribution

Stay tuned for updates!

## Next steps

- **[Getting started](../getting-started/index.md)** - Set up your first mesh network
- **[Analyst client](analyst-client.md)** - Connect from your workstation
- **[User guide](../user-guide/user-guide.md)** - Learn forensic workflows
- **[AmneziaWG Configuration](../advanced/amneziawg.md)** - Enable censorship resistance
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
