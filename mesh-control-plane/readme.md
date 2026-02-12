# MESH headscale UI

MESH UI wraps the Headscale API with a workflow designed for field deployments:

- **Network management** — Create, rename, and delete isolated networks. Each network maps to a Headscale user namespace.
- **Automatic network isolation** — When a network is created, ACL rules are auto-generated using `tag:net-<network>` patterns so devices in one network cannot reach devices in another. Custom ACL rules are preserved during sync.
- **Pre-auth key generation** — Generate keys with configurable expiration (1–365 days), reusable/ephemeral toggles, and device role tags (`ANALYST` or `MOBILE_NODE`).
- **Node monitoring** — View online/offline status, IP addresses, tags, and last-seen timestamps per network.
- **ACL policy editor** — Edit the raw HuJSON policy directly. Auto-generated isolation rules are merged with your custom rules on every network change.
- **Analyst email tracking** — Attach analyst emails to networks and filter the dashboard by analyst.

<img width="2476" height="1224" alt="image" src="https://github.com/user-attachments/assets/c3d93219-c3fc-428e-8b81-294b1fd16de8" />

<img width="2476" height="1224" alt="image" src="https://github.com/user-attachments/assets/e39dd96c-19e0-41ad-95d8-d5dcb78a708b" />

## Start with Docker compose

From the `mesh-control-plane/` directory:

```bash
docker compose up -d
```

This starts three services:
- **headscale** on port 8081
- **nginx proxy** on port 8080
- **mesh_ui** on port 3000

Then open `http://localhost:3000` and log in with your Headscale API key.

### Generating an API Key

```bash
docker exec headscale headscale apikeys create
```
