# Prerequisites

Before you begin setting up MESH, ensure you have the necessary requirements for each component you plan to deploy.

## Control plane

The control plane is the coordination server that manages your mesh network.

!!! warning "Security Note"
    MESH is built to be lightweight, ephemeral and quick to deploy.
    Most devices can host the control plane, but remember this device will be exposed to the public internet either via it's public IP/Domain or via a reverse proxy.

**Server requirements:**

- A Linux server (VPS or local machine)
- Docker and Docker Compose installed
- At least 1GB RAM and 10GB disk space

**Software:**

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later

!!! tip "VPS Providers"
    Popular VPS providers for hosting the control plane include:

    - DigitalOcean (starting at $6/month)
    - Linode (starting at $5/month)
    - Vultr (starting at $5/month)
    - AWS Lightsail (starting at $3.50/month)

## Analyst client

The analyst client runs on a secure forensics node and provides the MESH CLI for managing connections and running forensic tools.

!!! warning "Security Note"
    While MESH locks down ports and services between the analyst node and the compromised device, you must treat that network path as untrusted. The analyst node must therefore be a hardened, dedicated system that can be snapshotted after forensic collection, and must not be your everyday personal workstation. Although it would require a very complex attack to achieve lateral movement from a mobile device to a Linux-based analyst node, you should still assume the risk exists and use a secure, controlled node for this purpose.”

**Operating system:**

- Linux (Ubuntu 20.04+, Debian 11+, Fedora 35+, or similar)
- macOS 11 (Big Sur) or later

**Software requirements:**

- Go 1.21 or later
- Git
- Root/sudo access
- ADB (Android Debug Bridge)

**Optional tools:**

- AndroidQF - For automated spyware detection
- MVT (Mobile Verification Toolkit) - For iOS/Android forensics
- tmux - For running the daemon in the background

**Installing prerequisites on Ubuntu/Debian:**

```bash
# Update package list
sudo apt update

# Install Go
sudo apt install golang-go

# Install Git
sudo apt install git

# Install ADB
sudo apt install adb

# Install tmux (optional)
sudo apt install tmux
```

**Installing prerequisites on macOS:**

```bash
# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go
brew install go

# Install Git
brew install git

# Install ADB
brew install android-platform-tools

# Install tmux (optional)
brew install tmux
```

## Endpoint client

The endpoint client runs on Android devices that you want to connect to the mesh for forensic collection. You should distribute the APK to your clients who you're conducting forensics for.

**Device Requirements:**

- Android 8.0 (Oreo) or later
- WiFi connectivity
- At least 50MB free storage space

**Development Setup (for building APK):**

- Android Studio or Android SDK command-line tools
- Java Development Kit (JDK) 11 or later
- Gradle (usually bundled with Android Studio)

**Enabling Developer Options on Android:**

1. Go to **Settings** → **About phone**
2. Tap **Build number** 7 times
3. Go back to **Settings** → **System** → **Developer options**
4. Enable **USB debugging**

!!! warning "Security Note"
    The client should only enable USB debugging during the forensics collection. Disable it after.

## Network considerations

### Firewall configuration

If you're running the control plane behind a firewall, ensure the required ports are open:

```bash
# Example: UFW (Ubuntu/Debian)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 3478/udp

# Example: firewalld (RHEL/CentOS/Fedora)
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=3478/udp
sudo firewall-cmd --reload
```

### DNS configuration

For production deployments, you'll need:

- A domain name (e.g., `mesh.yourdomain.com`) or reverse proxy address
- DNS A record pointing to your control plane server's IP address
- (Optional) SSL/TLS certificate (can use Let's Encrypt)

### NAT and restrictive networks

MESH is designed to work in restrictive network environments:

- **NAT traversal** - Uses STUN/DERP for peer-to-peer connections behind NAT
- **Censorship resistance** - Optional AmneziaWG obfuscation for censored networks
- **Relay fallback** - Automatically uses DERP relay if direct connection fails

---

**Next:** [Control plane Setup](control-plane.md) →
