<div align="center">
  <img width="150" height="27" alt="Unti150" src="https://github.com/user-attachments/assets/32b90d10-b7c8-4808-b2fd-c1fe1c6bcbcf" />
</div>
<br>

<div align="center">
  <p>
    <a href="https://meshforensics.org/docs">
      <img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" />
    </a>
    <a href="https://discord.com/invite/">
      <img src="https://img.shields.io/discord/1161119546170687619?logo=discord&style=flat-square" alt="Chat" />
    </a>
  </p>
</div>

<div align="center">
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
<br/>

> [!IMPORTANT]
> Please be aware this project is currently in public alpha.

MESH makes remote mobile forensics secure, scalable and easy. 

It’s a censorship-resisting, peer-to-peer, end-to-end encrypted overlay network for digital forensics. It’s a fork of the [Tailscale](https://github.com/tailscale/tailscale) protocol, but is completely self-hostable and heavily modified for civil society and forensic operations.

MESH adds hardened transport and obfuscation options such as [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) for hostile or censored networks, and falls back to encrypted HTTPS relays when UDP is blocked. This offers protection against the GFW and detection by Deep Packet Inspection (DPI) systems.

It can create and tear down end-to-end peer-to-peer links between a locked-down analyst host and forensic clients (for example, Android devices) in seconds, without time consuming configurations that rely on hub-and-spoke topologies. It provides automated, segregated networks/namespaces for each investigation or analyst, making it easy to conduct multiple investigations at once for helplines. MESH lets you spin up a forensic mesh, collect data, tear it down, and start again without complex configuration. This is remote forensic and network capture without a hub-and-spoke design.

It’s remote investigations that work and scale to meet the needs of civil society.

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

## Remote mobile forensics and auditing

MESH creates an encrypted overlay network and assigns addresses from CGNAT ranges to the connected devices. This makes mobile devices reachable over the mesh without exposing them to the wider network.

With a CGNAT-assigned address, Android devices can be accessed over ADB-over-WiFi for collection of artifacts such as bug reports and `dumpsys` output, using tools like [AndroidQF](https://github.com/mvt-project/androidqf). iOS devices can be reached over the same mesh for use with tools like [libimobiledevice](https://github.com/libimobiledevice/libimobiledevice), enabling remote acquisition and analysis workflows even when the devices are not on the same physical LAN.

Because the overlay is ephemeral, analysts can bring devices into the mesh, perform live-state collection, and tear the network down immediately after. 

ACL configuration on the control plane lets analysts lock down services for suspected or compromised devices, restrict who can reach them, and enforce strict access controls against lateral movement for the duration of the investigation.

## Encryption, P2P & censorship resistance.

MESH uses a self-hostable control plane to coordinate the sharing of WireGuard keys between nodes letting you establish direct peer-to-peer connections whenever possible. Each node attempts NAT traversal using UDP hole punching and STUN-like techniques so two peers can exchange encrypted WireGuard packets directly. When hole punching succeeds, traffic is fully peer-to-peer and end-to-end encrypted between the two endpoints.

If UDP hole punching is unavailable or UDP is blocked, MESH falls back to the **DERP (Detoured Encrypted Relay for Packets)** protocol, which provides a relay network for failed or asymmetric NAT traversal. It relays already-encrypted WireGuard packets through DERP servers, so the relay sees metadata but never plaintext. DERP relays can be self-hosted by anyone and is fully open-source. See [DERP](https://github.com/tailscale/tailscale/tree/main/cmd/derper#derp).

Because WireGuard is actively censored in regions such as Russia and China, MESH can be configured to use [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) by simply adding an AmneziaWG config to your client (but also offers backward compatability). This obfuscates WireGuard packet fingerprints.

When UDP is blocked and the connection fails over to DERP, all encrypted traffic is sent over HTTPS via relays, which still support end-to-end encryption, providing a censorship-resilient mesh network.


## Get started

Docs coming soon.


