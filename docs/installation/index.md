# Configuration

Detailed installation and configuration guides for all MESH components.

## Overview

This section provides comprehensive installation and configuration documentation for:

- **Control plane** - Self-hosted coordination server
- **Analyst client** - Forensic workstation client (Linux/macOS)
- **Endpoint Clients** - Target device clients (Android/iOS)

## Configuration guides

<div class="grid cards" markdown>

- :material-server:{ .lg .middle } **Control plane**

    ---

    After deployment follow this control plane configuration guide.

    **Platform:** Linux with Docker  

    [:octicons-arrow-right-24: Configure](control-plane.md)

- :material-laptop:{ .lg .middle } **Analyst client**

    ---

    Install the analyst client on your forensic node to connect to remote devices and perform investigations.This guide gives you details about configure the client to suit your needs.

    **Platforms:** Linux, macOS  

    [:octicons-arrow-right-24: Install Analyst client](analyst-client.md)

- :material-cellphone:{ .lg .middle } **Endpoint Clients**

    ---

    Configure and build the endpoint client for your needs.

    **Platforms:** Android (iOS coming Q4 2026)  

    [:octicons-arrow-right-24: Install Endpoint client](endpoint-clients.md)

</div>

## After configuration

Once you've installed all components:

1. **Verify connectivity** - Test connections between nodes
2. **Configure ACLs** - Set up access control policies
3. **Run forensic tools** - Start your investigation

See the [User guide](../user-guide/index.md) for forensic workflows and usage instructions.

## Need help?

- **Quick setup:** [Getting started guide](../getting-started/index.md)
- **Troubleshooting:** [Troubleshooting guide](../reference/troubleshooting.md)
- **Community support:** [GitHub Discussions](https://github.com/BARGHEST-ngo/mesh/discussions)

---

**Next:** Install the [Control plane](control-plane.md) →
