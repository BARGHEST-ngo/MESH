<h1 align="center">MESH</h1>

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

## MESH

> [!IMPORTANT]
> Please be aware this project is currently in public alpha.

MESH is a censorship-resisting, peer-to-peer first, end-to-end encrypted overlay network for digital forensics.

It's a fork of the [Tailscale](https://github.com/tailscale/tailscale) protocol, but is self-hostable and heavily modified for civil society and forensic operations.

MESH adds hardened transport and obfuscation options such as [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) for hostile or censored networks, and falls back to encrypted HTTPS relays when UDP is blocked. This offers protection against GFW and detection by Deep Packet Inspection (DPI) systems.

It can create and tear down end-to-end mesh links between a locked-down analyst host and forensic clients (for example, Android devices) in seconds.

By removing hub-and-spoke topologies, MESH lets you spin up a forensics mesh, collect data, tear it down, and start again without complex configuration. This is remote forensic and network capture without a hub-and-spoke design.

Remote investigations that work.

Key functions

- Builds peer-to-peer subnets for device analysis.
- Provides end-to-end encryption via WireGuard.
- Runs a self-hostable control plane to create and manage meshes between analysts, enforce ACLs, and block compromised nodes on connection.
- Advertises LAN routes or exit nodes for monitoring and extended forensics.
- Supports ephemeral deployment and teardown on disconnect.
- Creates a virtual TUN interface and assigns CGNAT-range addresses for use with ADB-over-Wifi & [Libimobile](https://github.com/libimobiledevice/libimobiledevice).
- Transfers forensic artifacts such as ADB bug reports and `dumpsys` output.
- Intergrates directly with [AndroidQF](https://github.com/mvt-project/androidqf) for spyware IOC checks with [MVT](https://github.com/mvt-project/mvt)
- Supports kill-switch capabilities to block a device’s other network traffic while the forensics link remains active.
- Enables rapid creation, isolation, and teardown of remote investigation nodes.

## Encryption, P2P & Censorship resistance.

MESH uses a self-hostable control plane to coordinate the sharing of WireGuard keys between nodes letting you establish direct peer-to-peer connections whenever possible. Each node attempts NAT traversal using UDP hole punching and STUN-like techniques so two peers can exchange encrypted WireGuard packets directly. When hole punching succeeds, traffic is fully peer-to-peer and end-to-end encrypted between the two endpoints.

If UDP hole punching is unavailable or UDP is blocked, MESH falls back to the **DERP (Detoured Encrypted Relay for Packets)** protocol, which provides a relay network for failed or asymmetric NAT traversal. It relays already-encrypted WireGuard packets through DERP servers, so the relay sees metadata but never plaintext. DERP relays can be self-hosted by anyone and is fully open-source. See [DERP](https://github.com/tailscale/tailscale/tree/main/cmd/derper#derp).

Because WireGuard is actively censored in regions such as Russia and China, MESH can be configured to use [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) by simply adding an AmneziaWG config to your client (but also offers backward compatability). This obfuscates WireGuard packet fingerprints.

When UDP is blocked and the connection fails over to DERP, all encrypted traffic is sent over HTTPS via relays, which still support end-to-end encryption, providing a censorship-resilient mesh network.


## Get started

Docs coming soon.


