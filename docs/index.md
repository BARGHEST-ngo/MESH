---
hide:
  - navigation
  - toc
---

# [ MESH ]

!!! warning

    Please be aware this project is currently in public alpha. We take security seriously, and we recommend not using this in production till we have conducted a full penetration-test. This is scheduled at the start of 2026. Features and APIs may change. Please report issues on [GitHub](https://github.com/BARGHEST-ngo/mesh).

MESH is a censorship-resisting, peer-to-peer first, end-to-end encrypted overlay network for digital forensics. It's a fork of the [Tailscale](https://github.com/tailscale/tailscale) protocol, but is self-hostable and heavily modified for civil society and forensic operations.

**Key features:**

- **Censorship resistance** - Uses [AmneziaWG](https://github.com/amnezia-vpn/amneziawg-go) obfuscation for hostile or censored networks, and falls back to encrypted HTTPS relays when UDP is blocked. This offers protection against the GFW and detection by Deep Packet Inspection (DPI) systems.
- **Peer-to-peer connections** - Creates direct encrypted tunnels between devices, removing the need for hub-and-spoke topologies
- **Rapid deployment** - Spin up a forensics mesh, collect data, tear it down, and start again in minutes without complex configuration

**What you can do with MESH:**

- **Remote Android forensics** - Connect to Android devices anywhere in the world via ADB-over-WiFi
- **Automated spyware detection** - Integrates with [AndroidQF](https://github.com/mvt-project/androidqf) and [MVT](https://github.com/mvt-project/mvt) for IOC checks
- **Network monitoring** - Capture network traffic (PCAP) from remote devices
- **Secure data transfer** - Collect forensic artifacts like bug reports and system dumps over encrypted connections
- **Team collaboration** - Multiple analysts can work on different devices through the same secure network
- **Quick setup and teardown** - Create investigation networks in minutes, tear them down when done

## Mesh networking, made for remote mobile forensics

<div class="grid cards" markdown>

- :fontawesome-solid-hexagon-nodes:{ .lg .middle} **MESH networking**

    ---

    Discover why using a MESH network is a powerful way to build secure ephemeral remote forensic networks for threat analysis and forensics

    [:octicons-arrow-right-24: View features](overview/features.md)

- :material-network:{ .lg .middle } **Architecture**

    ---

    Learn about the three core components that power MESH

    [:octicons-arrow-right-24: Explore architecture](overview/architecture.md)

- :material-briefcase:{ .lg .middle } **Use cases**

    ---

    See how MESH helps investigators in real-world scenarios

    [:octicons-arrow-right-24: Read use cases](overview/use-cases.md)

- :material-devices:{ .lg .middle } **Platform support**

    ---

    Check which platforms are supported and what's coming next

    [:octicons-arrow-right-24: View platforms](overview/platform-support.md)

</div>

## Quick start

Ready to get started? Follow our step-by-step guide:

<div class="grid cards" markdown>

- :material-rocket-launch:{ .lg .middle } **Getting started guide**

    ---

    Complete setup guide from zero to your first forensic collection

    [:octicons-arrow-right-24: Start here](getting-started/index.md)

</div>

**Or jump directly to a specific component:**

1. **[Set up the control plane](getting-started/control-plane.md)** - Deploy your self-hosted coordination server
2. **[Install the analyst client](getting-started/analyst-client.md)** - Set up your Linux/macOS forensic workstation
3. **[Deploy the endpoint client](getting-started/endpoint-client.md)** - Install the MESH app on target devices
4. **[Verify your setup](getting-started/verification.md)** - Test connectivity and run your first collection

## Learn more

### Documentation

- **[Architecture](overview/architecture.md)** - Deep dive into how MESH works under the hood
- **[User guide](user-guide/user-guide.md)** - Learn forensic workflows and best practices
- **[AmneziaWG configuration](advanced/amneziawg.md)** - Configure censorship resistance
- **[CLI reference](reference/cli-reference.md)** - Complete command-line documentation
- **[Troubleshooting](reference/troubleshooting.md)** - Common issues and solutions

### About the project

- **[About BARGHEST](overview/about.md)** - Learn about the team behind MESH
- **[Platform support](overview/platform-support.md)** - Current and planned platform support
- **[Get help](overview/about.md#get-help)** - Community, support, and contributing

---
