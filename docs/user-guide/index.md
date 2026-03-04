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
sudo ./meshcli adbpair --host 100.63.x.x --hostport 1234 --pairport 1234 --code 1234

#Connect and pair BUT initate AndroidQF instantly on connection
sudo ./meshcli adbpair --host 100.63.x.x --hostport 1234 --pairport 1234 --code 1234 --qf

#Initate ADB acquision using AndroidQF and WARD libraries
sudo ./meshcli abdcollect 

# Check mesh status
sudo ./meshcli status

# List connected nodes
sudo ./meshcli status --peers

# Test connection and establish UDP hole punching
sudo ./meshcli ping 100.67.x.x

# Connect to Android device
adb connect 100.64.x.x:5555

# Run AndroidQF scan
androidqf --adb 100.64.x.x:5555 --output ./artifacts/

# Capture network traffic
sudo tcpdump -i mesh0 -w capture.pcap
```

## Getting Help

- **Troubleshooting:** [Troubleshooting guide](../reference/troubleshooting.md)
- **CLI reference:** [CLI reference](../reference/cli-reference.md)
- **Community:** [GitHub Discussions](https://github.com/BARGHEST-ngo/mesh/discussions)

---

**Next:** Read the [User guide](user-guide.md) for detailed workflows →
