<h1><div align="center">
  <img width="150" height="27" alt="Unti150" src="https://github.com/user-attachments/assets/32b90d10-b7c8-4808-b2fd-c1fe1c6bcbcf" />
</div></h1>

<div align="center">
  <p>
    <a href="https://docs.meshforensics.org/">
      <img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" />
    </a>
    <a href="https://discord.com/invite/">
      <img src="https://img.shields.io/discord/1161119546170687619?logo=discord&style=flat-square" alt="Chat" />
    </a>
  </p>
</div>

<h3><div align="center">
  <h3>
    <a href="https://meshforensics.org">
      MESH Docs
    </a>
    <span>  |  </span>
    <a href="https://Barghest.asia">
      BARGHEST
    </a>
  </h3>
</div>
<br/></h3>

> [!IMPORTANT]
> Please be aware this project is currently in public alpha. We take security seriously, and thus we recommend not using this in production till we have conducted a full penetration-test. This is currently on going. Features and APIs may change. Please report issues here on GitHub.

MESH is a censorship-resisting, peer-to-peer first, end-to-end encrypted overlay network for digital forensics. It's a fork of the [Tailscale](https://github.com/tailscale/tailscale) protocol, but is self-hostable and heavily modified for civil society and forensic operations.

MESH adds hardened transport and obfuscation options such as [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) for hostile or censored networks, and falls back to encrypted HTTPS relays when UDP is blocked. This offers protection against the GFW and detection by Deep Packet Inspection (DPI) systems.

It can create and tear down end-to-end mesh links between a locked-down analyst host and forensic clients (for example, Android devices) in seconds by removing hub-and-spoke topologies. MESH lets you spin up a forensics mesh, collect data, tear it down, and start again without complex configuration. This is remote forensic and network capture without a hub-and-spoke design.

Remote investigations that work.

Key functions

- Builds peer-to-peer subnets for device analysis.
- Provides end-to-end encryption via WireGuard/AmneziaWG and distributes keys automatically.
- Runs a self-hostable control plane to create and manage meshes between analysts, enforce ACLs, and block compromised nodes on connection.
- Advertises LAN routes or exit nodes for monitoring and extended forensics.
- Supports ephemeral deployment and teardown on disconnect.
- Creates a virtual TUN interface and assigns CGNAT-range addresses for use with ADB-over-Wifi & [Libimobile](https://github.com/libimobiledevice/libimobiledevice).
- Transfers forensic artifacts such as ADB bug reports and `dumpsys` output.
- Intergrates directly with [AndroidQF](https://github.com/mvt-project/androidqf) for spyware IOC checks with [MVT](https://github.com/mvt-project/mvt)
- Supports kill-switch capabilities to block a device’s other network traffic while the forensics link remains active.
- Enables rapid creation, isolation, and teardown of remote investigation nodes.

# Remote mobile forensics and auditing

MESH creates an encrypted overlay network and assigns addresses from CGNAT ranges to the connected devices. This makes mobile devices reachable over the mesh without exposing them to the wider network.

With a CGNAT-assigned address, Android devices can be accessed over ADB-over-WiFi for collection of artifacts such as bug reports and `dumpsys` output, using tools like [AndroidQF](https://github.com/mvt-project/androidqf). iOS devices can be reached over the same mesh for use with tools like [libimobiledevice](https://github.com/libimobiledevice/libimobiledevice), enabling remote acquisition and analysis workflows even when the devices are not on the same physical LAN.

Because the overlay is ephemeral, analysts can bring devices into the mesh, perform live-state collection, and tear the network down immediately after. 

ACL configuration on the control plane lets analysts lock down services for suspected or compromised devices, restrict who can reach them, and enforce strict access controls against lateral movement for the duration of the investigation.

# Encryption, P2P & censorship resistance.

MESH uses a self-hostable control plane to coordinate the sharing of WireGuard keys between nodes letting you establish direct peer-to-peer connections whenever possible. Each node attempts NAT traversal using UDP hole punching and STUN-like techniques so two peers can exchange encrypted WireGuard packets directly. When hole punching succeeds, traffic is fully peer-to-peer and end-to-end encrypted between the two endpoints.

If UDP hole punching is unavailable or UDP is blocked, MESH falls back to the **DERP (Detoured Encrypted Relay for Packets)** protocol, which provides a relay network for failed or asymmetric NAT traversal. It relays already-encrypted WireGuard packets through DERP servers, so the relay sees metadata but never plaintext. DERP relays can be self-hosted by anyone and is fully open-source. See [DERP](https://github.com/tailscale/tailscale/tree/main/cmd/derper#derp).

Because WireGuard is actively censored in regions such as Russia and China, MESH can be configured to use [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) by simply adding an AmneziaWG config to your client (but also offers backward compatability). This obfuscates WireGuard packet fingerprints.

When UDP is blocked and the connection fails over to DERP, all encrypted traffic is sent over HTTPS via relays, which still support end-to-end encryption, providing a censorship-resilient mesh network.

# Getting started

The easiest way to get started is to build a control plane and then join nodes. We recommend following the official <a href="https://docs.meshforensics.org/">documentation</a>, but if you want to get started from here, do the following:

```
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh/mesh-control-plane
```
Create a configuration file:
> [!NOTE]
> Replace your-domain.com with your actual domain or server IP address. If you want to use your own DERP server, you should change the urls to your own. You can learn more about DERP here. We use Tailscale's DERP relay's by default.


```
mkdir -p config
cat > config/config.yaml << EOF
server_url: https://your-domain.com
listen_addr: 0.0.0.0:8080
metrics_listen_addr: 0.0.0.0:9090

grpc_listen_addr: 0.0.0.0:50443
grpc_allow_insecure: false

private_key_path: /var/lib/headscale/private.key
noise_private_key_path: /var/lib/headscale/noise_private.key
base_domain: your-domain.com

ip_prefixes:
  - 100.64.0.0/10
  - fd7a:115c:a1e0::/48

derp:
  server:
    enabled: true
    region_id: 999
    region_code: "custom"
    region_name: "Custom DERP"
    stun_listen_addr: "0.0.0.0:3478"

  urls:
    - https://controlplane.tailscale.com/derpmap/default

  auto_update_enabled: true
  update_frequency: 24h

database:
  type: sqlite3
  sqlite:
    path: /var/lib/headscale/db.sqlite

log:
  level: info
  format: text

dns_config:
  override_local_dns: false
  nameservers:
    - 1.1.1.1
    - 8.8.8.8
  magic_dns: true
  base_domain: mesh.local

acl_policy_path: /etc/headscale/acl.yaml
EOF
```

Create a ACL policy:
Access Control Lists (ACLs) define which nodes can communicate with each other.
>[!IMPORTANT]
>The default configuration allows all nodes to communicate with each other. For production deployments, see the ACL documentation for more restrictive policies.

```
cat > config/acl.yaml << EOF
acls:
  - action: accept
    src:
      - "*"
    dst:
      - "*:*"
EOF
```

Start the control plane

```
docker-compose up -d
```

Access the Web UI

The control plane includes a web-based management interface. Access it at:
```
    Local access: https://localhost:3000/login
    Remote access: https://your-domain:8443/login
```

>[!IMPORTANT]
>Self-Signed Certificate
>The web UI uses a self-signed certificate by default. Your browser will show a security warning - this is expected. Click "Advanced" and proceed to the site. You should use a reverse proxy to serve the web UI over HTTPS. It is required for the control plane to be accessible on HTTPS publicly

Before using the web UI, create an API key:

```
docker exec headscale headscale apikeys create --expiration 90d
```

Example output:
```
abc123def456ghi789jkl012mno345pqr678stu901vwx234yz
```

Connect to the Web UI

    Open the web UI in your browser
    Enter your Ccntrol plane URL: http://localhost:3000
    Paste your API Key from Step 6
    Click ACCESS SYSTEM

<img width="1035" height="868" alt="image" src="https://github.com/user-attachments/assets/52bb4020-18fe-4b49-9b33-516250055278" />

Your MESH network is now build and ready to accept nodes, you should see the <a href="https://docs.meshforensics.org/">documentation</a> further for joining nodes and doing forensics analysis.

## Repository structure

This repository is a monorepo and contains all components required for MESH:
- `mesh-android-client`: The core library the Android APK that should be installed on endpoints for analysis
- `mesh-control-plane`: The control plane server that manages your nodes in the mesh network and distributes keys.
- `mesh-linux-macos-analyst`: The CLI client and daemon for your analysis/collection node.


## Legal

WireGuard is a registered trademark of Jason A. Donenfeld.



