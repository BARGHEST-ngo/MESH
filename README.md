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
    <span> | </span>
    <a href="https://Barghest.asia">
      BARGHEST
    </a>
  </h3>
</div>
<br/>

## What is MESH

MESH is a peer-to-peer encrypted overlay network built on top of the [Tailscale](https://github.com/tailscale/tailscale) protocol, but heavily adapted for civil society and digital forensics use cases. It adds hardened transport and obfuscation options such as [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) to operate in hostile or censored networks.

It creates and tears down end-to-end mesh networks between a locked-down analyst host and untrusted Android or iOS devices in seconds.

Build a forensics mesh, acquire the data, tear it down, and restart without complex configuration. This is remote forensic and network capture without a hub-and-spoke topology.

Remote investigations that work.

Key functions:

- Builds P2P subnets for device analysis.
- End-to-end encrypts using WireGuard.
- Enabled a control plane to manage mesh networks between analysts, control mesh ACLs, and restrict access for compromised nodes on connection.
- Advertise LAN routes or an exit nodes for network monitoring or further forensics use cases. 
- Allows for ephemeral deployment and teardown on disconnection.
- Creates a virtual TUN interface (assigning an address from CGNAT ranges)
- Transfers forensic artifacts such as ADB bug reports and dumpsys data.
- Enables rapid creation, isolation, and teardown of remote investigation nodes.

## Encryption, P2P & Censorship resistance.

MESH uses a self-hostable control plane to coordinate the sharing of WireGuard keys between nodes. Each node then attempts NAT traversal using UDP hole punching and STUN-like techniques so two peers can exchange encrypted WireGuard packets directly. When hole punching succeeds, traffic is fully end-to-end encrypted WireGuard between the two endpoints.

If UDP hole punching is unavailable or UDP is blocked, MESH falls back to the **DERP (Detoured Encrypted Relay for Packets)** protocol, which provides a relay network for failed or asymmetric NAT traversal. It relays already-encrypted WireGuard packets through DERP servers, so the relay sees metadata but never plaintext. Relays can be self-hosted by anyone. See [DERP servers](https://tailscale.com/kb/1232/derp-servers).

Because WireGuard is actively censored in regions such as Russia and China, MESH uses WireGuard by default but can be configured to use [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) by adding an AmneziaWG config to your client. This obfuscates WireGuard packet fingerprints.

When UDP is blocked and the connection fails over to DERP, all encrypted traffic is sent over HTTPS via relays, providing a censorship-resilient mesh network.

