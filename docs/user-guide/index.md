# User guide

Learn how to use MESH for forensic investigations, manage your mesh network, and perform common tasks.

## Overview

This section covers:

- **Forensic workflows** - How to perform remote device analysis
- **Network management** - Managing nodes, users, and access control
- **Common tasks** - Day-to-day operations and best practices
- **Architecture details** - Deep dive into how MESH works

## Documentation

<div class="grid cards" markdown>

- :material-book-open-variant:{ .lg .middle } **User guide**

    ---

    Complete guide to using MESH for forensic investigations, including workflows, tools, and best practices.

    [:octicons-arrow-right-24: View User guide](user-guide.md)

- :material-sitemap:{ .lg .middle } **Architecture Details**

    ---

    Deep dive into MESH architecture, including technical details, protocols, and implementation.

    [:octicons-arrow-right-24: View Architecture](architecture.md)

</div>

## Quick Reference

### Common commands

```bash
#Connect and pair to an Android device over ADB-over-Wifi
meshcli adbpair --host 100.64.x.x --hostport 1234 --pairport 4321 --code 123456

#Connect and pair BUT initiate AndroidQF instantly on connection
meshcli adbpair --host 100.64.x.x --hostport 1234 --pairport 4321 --code 123456 --qf

#Initate ADB acquisition using AndroidQF and WARD libraries
meshcli adbcollect

# Check mesh status
meshcli status

# Test connection and establish UDP hole punching
meshcli ping 100.64.x.x

# Capture network traffic
tcpdump -i tailscale0 -w capture.pcap
```

## Getting Help

- **Troubleshooting:** [Troubleshooting guide](../reference/troubleshooting.md)
- **CLI reference:** [CLI reference](../reference/cli-reference.md)
- **Community:** [GitHub Discussions](https://github.com/BARGHEST-ngo/mesh/discussions)

---

**Next:** Read the [User guide](user-guide.md) for detailed workflows →
