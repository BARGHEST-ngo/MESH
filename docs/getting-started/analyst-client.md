# Analyst client setup

Now let's set up the MESH client on your acquision node. This will allow you to connect to the mesh network and conduct forensic collections.

## Step 1: Clone and build

Clone the MESH repository and build the analyst client if you haven't already:

```bash
# Clone the repository
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh/mesh-linux-macos-analyst

# Build the MESH client
./build_mesh.sh
```

The build process will compile two binaries:

- `meshcli` - The client CLI for managing connections
- `tailscaled-amnezia` - The background daemon that maintains the mesh connection

!!! tip "Build Time"
    The build process may take 5-10 minutes depending on your system. Go is compiling the entire MESH client from source.

## Step 2: Start the MESH daemon

The MESH daemon runs in the background and maintains your connection to the mesh network.

**Create state directory**

```bash
sudo mkdir -p /var/lib/mesh
```

**Start the daemon**

```bash
./start_mesh_daemon.sh
```

Or to run with custom parameters:

```bash
sudo ./tailscaled-amnezia \
  --socket=/var/run/mesh/tailscaled.sock \
  --state=/var/lib/mesh/tailscaled.state \
  --statedir=/var/lib/mesh
```

!!! info "Keep Terminal Open"
    Keep this terminal open, or run the daemon in a tmux session for persistent operation.

**Running in tmux (Recommended):**

```bash
# Install tmux if not already installed
sudo apt install tmux  # Ubuntu/Debian
brew install tmux      # macOS

# Start a tmux session
tmux new -s mesh-daemon

# Run the daemon
sudo ./tailscaled-amnezia \
  --socket=/var/run/mesh/tailscaled.sock \
  --state=/var/lib/mesh/tailscaled.state \
  --statedir=/var/lib/mesh

# Detach from tmux: Press Ctrl+B, then D
# Reattach later: tmux attach -t mesh-daemon
# List sessions: tmux ls
```

## Step 3: Connect the acquisition node to the MESH network

In a new terminal, connect to your control plane using the pre-authentication key you created earlier. This will allow the control plane to distribute WG keys and allow your node to join the mesh.

```bash
sudo ./meshcli up \
  --login-server=https://your-domain.com \
  --authkey=abc123def456ghi789jkl012mno345pqr678stu901vwx234yz \
  --accept-dns=false
```

!!! important "Replace Values"
    - Replace `your-domain.com` with your control plane URL
    - Replace `abc123...` with your pre-auth key from the control plane setup

**What these flags mean:**

- `--login-server` - URL of your control plane
- `--authkey` - Pre-authentication key for automatic enrollment
- `--accept-dns=false` - Don't override system DNS (optional, use `true` for MagicDNS - you can find out more about this in the Advanced section)

## Step 4: Verify Connection

Check that your analyst client is connected to the mesh:

```bash
# Check connection status
sudo ./meshcli status

# List all peers in the mesh
sudo ./meshcli status --peers
```

**Example output:**

```
Health check:
    - in map poll: true
    - in keep alive: true
    - derp: connected

Logged in as: analyst1
Mesh IP: 100.64.1.1
```

!!! success "Connected!"
    If you see "in map poll: true" and a mesh IP address, you're successfully connected to the mesh!

## Step 5: Test basic functionality

**Check your Mesh IP**

```bash
sudo ./meshcli ip
```

This shows your assigned mesh IP address (e.g., `100.64.1.1`).

**View mesh status**

```bash
sudo ./meshcli status --json
```

This provides detailed status information in JSON format.

## Optional: Install as System Service

For production use, you may want to run the MESH daemon as a system service.

**Create systemd service file (Linux):**

```bash
sudo cat > /etc/systemd/system/mesh.service << EOF
[Unit]
Description=MESH Daemon
After=network.target

[Service]
Type=simple
ExecStart=/path/to/tailscaled-amnezia --socket=/var/run/mesh/tailscaled.sock --state=/var/lib/mesh/tailscaled.state --statedir=/var/lib/mesh
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Replace `/path/to/tailscaled-amnezia` with the actual path to your binary.

**Enable and start the service:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable mesh
sudo systemctl start mesh
sudo systemctl status mesh
```

## Troubleshooting

If you encounter issues during setup, see the [Troubleshooting guide](../reference/troubleshooting.md) for common problems and solutions:

- [Daemon won't start](../reference/troubleshooting.md#daemon-wont-start)
- [Can't connect to control plane](../reference/troubleshooting.md#cant-connect-to-control-plane)
- [Authentication fails](../reference/troubleshooting.md#authentication-fails)
- [No mesh IP assigned](../reference/troubleshooting.md#no-mesh-ip-assigned)

## Next steps

Your analyst client is now connected to the mesh. The next step is to install the endpoint client on an Android device.

For detailed analyst client documentation, see the [Analyst client documentation](../installation/analyst-client.md).

---

← [Previous: Control plane Setup](control-plane.md) | [Next: Endpoint client Setup](endpoint-client.md) →
