# Endpoint client Setup (Android)

The endpoint client runs on mobile devices and enables remote forensic collection over the mesh network.

!!! info "Current support is limited to Android"
    We currently only support Android forensics acquisions

## Step 1: Build the APK or download pre-compiled

**Download**

Download the APK from [GitHub Releases](https://github.com/BARGHEST-ngo/mesh/releases) page.

**Building from source**

Navigate to the Android client directory and build the APK:

```bash
cd mesh/mesh-android-client

# Build the release APK
./gradlew assembleRelease
```

The APK will be created at:

```
app/build/outputs/apk/release/app-release.apk
```

!!! tip "Build Time"
    The first build may take 10-15 minutes as Gradle downloads dependencies. Subsequent builds will be faster.

## Step 2: Install on Android device

!!! info "App is still in Alpha"
    We're still in public Alpha, so we haven't launched the app on Google Playstore or F-Droid yet.

Since the MESH app is not available on Google Play Store or F-Droid, you'll need to install it manually using one of the following methods:

### Method 1: Transfer via USB (Recommended)

This is the most reliable method for installing the APK.

**1. Enable USB debugging on the Android device:**

1. Open **Settings** on your Android device
2. Navigate to **About phone**
3. Tap **Build number** 7 times to enable Developer options
4. Go back to **Settings** → **System** → **Developer options**
5. Enable **USB debugging**

**2. Connect the device to your computer:**

```bash
# Connect your Android device via USB cable
# Verify the device is detected
adb devices
```

You should see your device listed. If prompted on the device, tap **"Allow"** to authorise USB debugging.

**3. Install the APK:**

```bash
# Install the APK using ADB
adb install app-release.apk

# Or if you downloaded from GitHub
adb install mesh-android-client.apk
```

**4. Verify installation:**

```bash
# Check if the app is installed
adb shell pm list packages | grep mesh
```

### Method 2: Transfer via file sharing (WiFi/Bluetooth)

If you don't have a USB cable or prefer wireless transfer:

**1. Enable installation from unknown sources:**

1. Open **Settings** on your Android device
2. Navigate to **Security** or **Privacy**
3. Find **Install unknown apps** or **Unknown sources**
4. Select your file manager app (e.g., Files, Downloads)
5. Enable **Allow from this source**

!!! warning "Security consideration"
    Only enable installation from unknown sources temporarily, and disable it after installing MESH.

**2. Transfer the APK to your device:**

Choose one of these methods:

**Option A: Email**

- Email the APK to yourself
- Open the email on your Android device
- Download the attachment

**Option B: Cloud storage**

- Upload the APK to Google Drive, Dropbox, or similar
- Download it on your Android device

**Option C: Direct WiFi transfer**

- Use a file transfer app like [LocalSend](https://localsend.org/) or [KDE Connect](https://kdeconnect.kde.org/)
- Transfer the APK wirelessly

**Option D: Bluetooth**

- Send the APK via Bluetooth from your computer
- Accept the file on your Android device

**3. Install the APK:**

1. Open your **Files** or **Downloads** app on the Android device
2. Navigate to where you saved the APK
3. Tap on the **APK file**
4. Tap **"Install"**
5. Wait for the installation to complete
6. Tap **"Open"** or find the MESH app in your app drawer

### Method 3: Install via ADB over WiFi (Advanced)

If you want to install without a USB cable and have already enabled USB debugging:

**1. Connect device and computer to the same WiFi network**

**2. Enable ADB over WiFi on the device:**

```bash
# First connect via USB
adb devices

# Enable TCP/IP mode on port 5555
adb tcpip 5555

# Disconnect USB cable
```

**3. Find your device's IP address:**

On the Android device:

1. Go to **Settings** → **About phone** → **Status** → **IP address**
2. Note the IP address (e.g., `192.168.1.100`)

**4. Connect via WiFi:**

```bash
# Connect to the device over WiFi
adb connect 192.168.1.100:5555

# Verify connection
adb devices
```

**5. Install the APK:**

```bash
# Install the APK wirelessly
adb install app-release.apk
```

### Troubleshooting installation

**"App not installed" error:**

- Ensure you have enough storage space (at least 100MB free)
- Try uninstalling any previous version of MESH first
- Reboot your device and try again

**"Installation blocked" error:**

- Make sure you've enabled installation from unknown sources
- Check that the APK file downloaded completely (not corrupted)
- Verify the APK is compatible with your Android version (8.0+)

**ADB device not found:**

- Ensure USB debugging is enabled
- Try a different USB cable or port
- Install the appropriate USB drivers for your device
- Revoke USB debugging authorisations and try again

## Step 3: Configure the app

Now configure the MESH app on your Android device.

### Open the app

1. On your Android device, open the app drawer
2. Find and tap the **MESH** app icon
3. You'll see the connection screen

### Connect to the mesh

1. Tap **"Connect to Network"**
2. Enter your **Control plane URL**: `https://your-domain.com`
3. Enter the **Pre-auth key** you created during control plane setup
4. Tap **"Connect"**

!!! important "Replace values"
    - Replace `your-domain.com` with your actual control plane URL
    - Use the same pre-auth key you created in the control plane setup (if it's reusable)

### Grant permissions

The app will request VPN permissions:

1. Android will show a connection request dialogue
2. Tap **"OK"** to allow MESH to create a VPN connection
3. The app will establish a connection to the mesh

!!! info "VPN permission"
    MESH uses Android's VPN API to create the mesh network interface. This is required for the app to function.

## Step 4: Verify connection

### Check connection status

In the MESH app, you should see:

- **Status**: Connected
- **Mesh IP**: Your assigned IP address (e.g., `100.64.2.1`)
- **Peers**: Number of connected peers

### Enable ADB over WiFi

The MESH app automatically enables ADB-over-WiFi on the mesh interface. This allows remote forensic collection.

**Verify ADB is enabled:**

1. In the MESH app, check the **ADB Status** section
2. It should show: **"ADB enabled on 100.64.2.1:5555"**

!!! success "Ready for forensics"
    Your Android device is now connected to the mesh and ready for remote forensic collection!

## Step 5: Test from analyst workstation

From your analyst workstation, verify you can see the Android device:

```bash
# List all peers
sudo ./meshcli status --peers
```

You should see your Android device listed with its mesh IP address.

## Troubleshooting

### App won't install

**Error: "App not installed"**

- Ensure USB debugging is enabled on the device
- Check that the device is authorised (check phone screen for prompt)
- Try uninstalling any existing version first

**Error: "INSTALL_PARSE_FAILED"**

- The APK may be corrupted. Try rebuilding:

  ```bash
  ./gradlew clean assembleRelease
  ```

### Can't connect to control plane

**"Connection failed" error:**

1. Verify the control plane URL is correct
2. Check that the device has internet connectivity (WiFi or mobile data)
3. Verify the pre-auth key is valid and not expired
4. Check control plane logs: `docker logs headscale`

**"Invalid auth key" error:**

- The pre-auth key may have expired
- Create a new pre-auth key in the control plane web UI
- Make sure you're using a reusable key if connecting multiple devices

### VPN permission denied

If Android won't grant VPN permission:

1. Go to **Settings** → **Apps** → **MESH**
2. Check that all permissions are granted
3. Try uninstalling and reinstalling the app

### ADB not enabled

If ADB-over-WiFi doesn't enable automatically:

1. Check the MESH app settings
2. Verify the device is connected to the mesh (check mesh IP)
3. Try toggling the ADB option in the app settings

### Device not visible from analyst

If the analyst workstation can't see the Android device:

```bash
# Check if both nodes are online
sudo ./meshcli status --peers

# Try pinging the device
ping 100.64.2.1  # Replace with actual mesh IP
```

If ping fails:

- Both devices may be behind restrictive NAT (MESH will use DERP relay)
- Check control plane logs for connection issues
- Verify firewall rules allow UDP traffic

## Next steps

Your endpoint client is now connected to the mesh! The next step is to verify connectivity and run your first forensic collection.

For detailed endpoint client documentation, see the [Endpoint client documentation](../installation/endpoint-clients.md).

---

← [Previous: Analyst client Setup](analyst-client.md) | [Next: Verification](verification.md) →
