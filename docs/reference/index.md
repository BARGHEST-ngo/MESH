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
meshcli status

# Enable/disable routes
meshcli up --advertise-routes=10.0.0.0/24
meshcli down

# Check version
./meshcli version
```

#### Control plane (headscale)

```bash
# List nodes
docker compose exec headscale headscale nodes list

# List networks
docker compose exec headscale headscale users list

# Check ACL policy
docker compose exec headscale headscale policy get
```

#### Android Forensics

```bash
# Connect to device
meshcli adbpair --host 100.64.x.x --hostport 1234 --pairport 4321 --code 123456

# Run AndroidQF
meshcli adbcollect

# Run MVT
mvt-android check-adb --output ./mvt-output/
```

### Common issues

| Issue | Solution |
|-------|----------|
| Can't connect to control plane | Check firewall, verify URL, check logs |
| Nodes can't see each other | Check ACL policy, verify routes |
| ADB connection fails | Enable ADB over network on device |
| Web UI not accessible | Check port 80 on the host, verify reverse-proxy/tunnel for remote access |
| High latency | Check if using DERP relay, prefer P2P |

[:octicons-arrow-right-24: Full Troubleshooting guide](troubleshooting.md)

### Port Reference

| Port | Service | Protocol | Purpose |
|------|---------|----------|---------|
| 80 | MESH web UI | HTTP | Web management interface |
| 3478 | DERP | UDP/TCP | STUN for NAT traversal |

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
