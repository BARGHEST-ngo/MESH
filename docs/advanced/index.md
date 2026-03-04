# Advanced configuration

Advanced topics for power users, including censorship resistance, custom configurations, and production deployments.

## Overview

This section covers advanced MESH configuration topics:

- **AmneziaWG** - Obfuscated WireGuard for censorship resistance
- **Network packet capture** - Capture PCAP's quickly and easily
- **Custom DERP servers** - Deploy your own relay infrastructure
- **Production hardening** - Security best practices for production
- **Advanced ACL policies** - Complex access control scenarios
- **Performance tuning** - Optimize for your use case

## Advanced topics

<div class="grid cards" markdown>

- :material-shield-lock-outline:{ .lg .middle } **AmneziaWG configuration**

    ---

    Configure AmneziaWG obfuscation to bypass deep packet inspection and censorship in restrictive networks.

    **Use case:** Censored environments

    **Difficulty:** Advanced

    [:octicons-arrow-right-24: Configure AmneziaWG](amneziawg.md)

- :material-server-network:{ .lg .middle } **Control plane advanced**

    ---

    Advanced control plane configuration including custom DERP servers and PostgreSQL database setup.

    **Use case:** Production deployments

    **Difficulty:** Advanced

    [:octicons-arrow-right-24: Advanced control plane](control-plane-advanced.md)

- :material-network-outline:{ .lg .middle } **Exit node and packet capture**

    ---

    Configure exit node routing and capture network traffic from endpoints for forensic analysis.

    **Use case:** Network forensics

    **Difficulty:** Advanced

    [:octicons-arrow-right-24: Exit node and PCAP](exit-node-pcap.md)

</div>

## Coming soon

Additional advanced topics will be added to this section:

## Prerequisites

Before diving into advanced topics, ensure you have:

1. **Working MESH deployment** - Complete the [Getting started guide](../getting-started/index.md)
2. **Basic understanding** - Familiar with MESH architecture and components
3. **Production experience** - Tested MESH in development environment
4. **Technical skills** - Comfortable with Linux, networking, and Docker

## Getting help

Advanced configuration can be complex. Resources:

- **Documentation:** Read the specific advanced topic guide
- **Community:** [GitHub Discussions](https://github.com/BARGHEST-ngo/mesh/discussions)
- **Issues:** [GitHub Issues](https://github.com/BARGHEST-ngo/mesh/issues)
- **Professional support:** Contact BARGHEST for consulting

## Safety first

!!! danger
     Advanced configurations can impact security and stability.

     **Best practices:**
     
     **Test in development** - Never test in production
     **Backup configurations** - Save working configs before changes
    - **Document changes** - Keep detailed records
    - **Have rollback plan** - Know how to revert changes

---

**Explore advanced topics:**

- [AmneziaWG configuration](amneziawg.md) - Censorship resistance
- [Control plane advanced](control-plane-advanced.md) - Custom DERP and PostgreSQL
- [Exit node and PCAP](exit-node-pcap.md) - Network forensics
