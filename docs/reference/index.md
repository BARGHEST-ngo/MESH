# Reference

Technical reference documentation, CLI commands, and troubleshooting guides.

## Overview

This section provides reference documentation for:

- **CLI reference** - Complete command-line interface documentation
- **Troubleshooting** - Common issues and solutions
- **API Reference** - Control plane API documentation (coming soon)
- **Configuration Reference** - All configuration options (coming soon)

## Reference Documentation

<div class="grid cards" markdown>

- :material-console:{ .lg .middle } **CLI reference**

    ---

    Complete reference for all MESH CLI commands, including meshcli, headscale, and related tools.

    [:octicons-arrow-right-24: View CLI reference](cli-reference.md)

- :material-help-circle:{ .lg .middle } **Troubleshooting**

    ---

    Common issues, error messages, and solutions for MESH deployment and operation.

    [:octicons-arrow-right-24: View Troubleshooting](troubleshooting.md)

</div>

## Quick Reference

### Common commands

#### Mesh Client (meshcli)

```bash
# Check status
sudo ./meshcli status

# List peers
sudo ./meshcli status --peers

# Enable/disable routes
sudo ./meshcli up --advertise-routes=10.0.0.0/24
sudo ./meshcli down

# Check version
./meshcli version
```

#### Control plane (headscale)

```bash
# List nodes
docker exec headscale headscale nodes list

# Create pre-auth key
docker exec headscale headscale preauthkeys create --user default --expiration 24h

# List users
docker exec headscale headscale users list

# Check ACL policy
docker exec headscale headscale policy check
```

#### Android Forensics

```bash
# Connect to device
adb connect 100.64.x.x:5555

# Run AndroidQF
androidqf --adb 100.64.x.x:5555 --output ./artifacts/

# Run MVT
mvt-android check-adb --output ./mvt-output/
```

### Common issues

| Issue | Solution |
|-------|----------|
| Can't connect to control plane | Check firewall, verify URL, check logs |
| Nodes can't see each other | Check ACL policy, verify routes |
| ADB connection fails | Enable ADB over network on device |
| Web UI not accessible | Check port 8443, verify TLS certificate |
| High latency | Check if using DERP relay, prefer P2P |

[:octicons-arrow-right-24: Full Troubleshooting guide](troubleshooting.md)

### Port Reference

| Port | Service | Protocol | Purpose |
|------|---------|----------|---------|
| 8080 | Headscale | HTTP | Control plane API |
| 8443 | Headscale-UI | HTTPS | Web management interface |
| 3478 | DERP | UDP/TCP | STUN for NAT traversal |
| 41641 | WireGuard | UDP | Encrypted mesh traffic |
| 5555 | ADB | TCP | Android Debug Bridge |

### Network Ranges

| Range | Purpose |
|-------|---------|
| 100.64.0.0/10 | MESH node addresses (CGNAT range) |
| 10.0.0.0/8 | Private subnet routing (example) |
| 172.16.0.0/12 | Private subnet routing (example) |
| 192.168.0.0/16 | Private subnet routing (example) |

## Documentation Sections

### CLI reference

Complete documentation for all command-line tools:

- **meshcli** - MESH client commands
- **headscale** - Control plane management
- **Docker commands** - Container management
- **Forensic tools** - ADB, AndroidQF, MVT integration

[:octicons-arrow-right-24: CLI reference](cli-reference.md)

### Troubleshooting

Solutions for common issues:

- **Connection problems** - Can't connect to control plane or peers
- **Performance issues** - High latency, low throughput
- **Configuration errors** - ACL policy, routing, DNS
- **Platform-specific** - Android, Linux, macOS issues

[:octicons-arrow-right-24: Troubleshooting guide](troubleshooting.md)

## Coming Soon

Additional reference documentation:

### API Reference

Complete REST API documentation for the control plane:

- **Authentication** - API key management
- **Nodes** - List, create, delete nodes
- **Users** - User management
- **Routes** - Route management
- **Pre-auth keys** - Key generation and management

### Configuration Reference

Complete reference for all configuration files:

- **meshcli.conf** - Client configuration
- **headscale.yaml** - Control plane configuration
- **acl.json** - Access control policies
- **docker-compose.yml** - Container orchestration

## Getting Help

Can't find what you're looking for?

- **Search the docs** - Use the search bar at the top
- **Check troubleshooting** - [Troubleshooting guide](troubleshooting.md)
- **Community support** - [GitHub Discussions](https://github.com/BARGHEST-ngo/mesh/discussions)
- **Report issues** - [GitHub Issues](https://github.com/BARGHEST-ngo/mesh/issues)

---

**Next:** Browse the [CLI reference](cli-reference.md) or [Troubleshooting guide](troubleshooting.md) →
