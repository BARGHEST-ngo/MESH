# Getting started with MESH

Welcome to the MESH getting started guide.

This guide will walk you through setting up your first MESH network for remote mobile forensics.

By the end, you'll have a working mesh network connecting an analyst workstation to an Android device for collecting Android forensics data.

## Setup Overview

The setup process involves four main steps:

1. **[Deploy the Control plane](control-plane.md)** - Set up the MESH coordination server
2. **[Install the Analyst client](analyst-client.md)** - Build and configure the MESH CLI on your acquision Linux/MacOS node
3. **[Install the Endpoint client](endpoint-client.md)** - Deploy the MESH app on a client Android device
4. **[Connect and Verify](verification.md)** - Establish the mesh and test connectivity

## Before you begin

Make sure you have the necessary [prerequisites](prerequisites.md) for each component you plan to deploy.

## Quick Navigation

<div class="grid cards" markdown>

- :material-server:{ .lg .middle } **Control plane**

    ---

    Set up the coordination server that manages your mesh network

    [:octicons-arrow-right-24: Deploy Control plane](control-plane.md)

- :material-laptop:{ .lg .middle } **Analyst client**

    ---

    Install the MESH client on your Linux/macOS workstation

    [:octicons-arrow-right-24: Install Analyst client](analyst-client.md)

- :material-cellphone:{ .lg .middle } **Endpoint client**

    ---

    Deploy the MESH app on Android devices for forensic collection

    [:octicons-arrow-right-24: Install Endpoint client](endpoint-client.md)

- :material-check-circle:{ .lg .middle } **Verification**

    ---

    Test connectivity and run your first forensic collection

    [:octicons-arrow-right-24: Verify Setup](verification.md)

</div>

## Need Help?

If you encounter issues during setup:

- Check the [Troubleshooting guide](../reference/troubleshooting.md)
- Review the [Architecture documentation](../overview/architecture.md) to understand how components work together
- Report issues on [GitHub](https://github.com/BARGHEST-ngo/mesh/issues)

---

**Ready to get started?** Begin with the [Prerequisites](prerequisites.md) →
