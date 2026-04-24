# Runtime stack

This page documents `compose.yml` and every `Dockerfile` that is pulled into the runtime stack. The devcontainer (`.devcontainer/Dockerfile`) is covered separately in [Build orchestration](build-orchestration.md#development-container) because it is a build-time environment, not a runtime component.

## Compose file

Path: `compose.yml` (repo root).

The Compose file defines three services that together form the control plane plus an analyst workstation. Services set `restart: unless-stopped` where the service is long-running, and rely on Compose defaults for networking (a single project-scoped bridge network, service-name DNS).

```yaml
services:
  headscale:
    image: docker.io/headscale/headscale:<pinned-tag>
    # config, state volumes, tmpfs for /var/run/headscale

  ui:
    build: ./control-plane/mesh-ui
    ports: ["127.0.0.1:80:80"]
    depends_on: [headscale]

  analyst:
    build: { context: ., dockerfile: analyst/Dockerfile }
    cap_add: [NET_ADMIN, NET_RAW]
    devices: ["/dev/net/tun:/dev/net/tun"]
```

### Service: `headscale`

Pulls the upstream Headscale image unchanged. MESH does not fork Headscale; the exact tag is pinned in `compose.yml`.

Three mounts shape its state:

- **`control-plane/headscale` to `/etc/headscale`** (bind mount). Holds `config.yaml` generated from `config.example.yaml` at first run. The path is expanded via `${LOCAL_WORKSPACE_FOLDER:-.}` so the devcontainer can bind-mount the host repo path rather than the path inside the container.
- **`headscale-data` to `/var/lib/headscale`** (named volume). SQLite database, noise private key, optional DERP private key. Survives `docker compose down` but removable with `docker volume rm`.
- **`/var/run/headscale`** (tmpfs). Unix socket and other runtime state that must not persist.

Entrypoint is `headscale serve`. No host ports are published directly. Traffic reaches Headscale through the `ui` service's nginx proxy.

### Service: `ui`

Built from `control-plane/mesh-ui/` via the two-stage `Dockerfile` described below.

**Environment variables:**

- **`CORS_ORIGIN`** (default `*`). Substituted into the `Access-Control-Allow-Origin` response header.
- **`NGINX_MODE`** (default `prod`). Switches the nginx config between `dev` (reverse-proxy all requests to Vite on `localhost:3000`) and `prod` (serve pre-built static files and fall through to `/index.html` for SPA routing).

The entrypoint does not use the image's default. It replaces it with:

```bash
envsubst '$CORS_ORIGIN $NGINX_MODE' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf && nginx -g 'daemon off;'
```

The bind-mounted `control-plane/nginx.conf` is treated as a template. Only `$CORS_ORIGIN` and `$NGINX_MODE` are substituted. Other nginx variables (such as `$http_upgrade` or `$uri`) survive because `envsubst` is given an explicit allowlist.

nginx proxies the headscale API endpoints (`/api`, `/register`, `/key`, `/machine`, `/derp`, `/verify`, `/health`, `/ts2021`, `/swagger`) to `http://headscale:8080` and serves (or proxies to Vite) everything else. Long-lived connection timeouts are set to 24 hours for WebSocket support.

Only port `80` is exposed, bound to `127.0.0.1`. A reverse proxy or tunnel is responsible for TLS termination and public exposure.

### Service: `analyst`

Built locally from `analyst/Dockerfile` with the repo root as build context, so Go source and the tailscale submodule are both available.

**Capabilities and devices:** `NET_ADMIN` and `NET_RAW` are required for WireGuard/AmneziaWG tunnel creation. `/dev/net/tun` is passed through from the host, so the host kernel must expose it. IP forwarding is enabled per-namespace via `sysctls`, which lets the container act as a gateway between its interfaces without changing host-wide kernel settings.

**Environment:**

- **`LOGIN_URL`**: defaults to `https://${CONTROL_PLANE_DOMAIN}` if unset. Wires the analyst to whichever control plane the user configured.
- **`AUTH_KEY`**: a Headscale pre-auth key. Required. No default.

Both values are read from `.env`.

**State:** the `analyst-data` named volume mounts at `/root/.tailscale`. Holds the daemon state and the two log files (`mesh.log`, `meshcli.log`). Removed by `task setAuthKey` when rotating the auth key.

### Volumes

Only two are declared:

- `headscale-data`: Headscale's persistent state.
- `analyst-data`: analyst daemon state.

Both are project-scoped. Compose prefixes volume names with the project name (derived from the directory, lowercased). Typically `mesh_headscale-data`, but `mesh-fork_headscale-data` if you've cloned into `mesh-fork/`. List them with `docker volume ls`.

---

## `analyst/Dockerfile`

Path: `analyst/Dockerfile`. Three stages. Final image: a slim Debian runtime with a JRE, Android command-line tooling, and the mesh daemon.

### Stage 1: Go build stage

Compiles the mesh daemon from the tailscale submodule.

```dockerfile
FROM golang:<pinned>-alpine AS build
RUN apk add --no-cache git patch
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY analyst/src ./analyst/src
RUN go generate ./...
RUN ./analyst/src/build.sh
```

`go mod download` runs before the source copy so the layer is cacheable across code changes. `analyst/src/build.sh` is the single source of truth for the build command. It is called both here and from `task buildAnalyst`. See [Build orchestration](build-orchestration.md#analyst-build-script) for what it does.

### Stage 2: Android SDK bootstrap stage

Installs the Android SDK command-line tools in a throwaway stage, so the final image can `COPY --from=android-sdk` without having to install `wget`, `unzip`, or the JDK's build-time dependencies into the runtime.

```dockerfile
ARG ANDROID_TOOLS_VERSION=<pinned>
ARG ANDROID_TOOLS_SUM=<vendored-sha256>
# wget -> sha256sum --strict --check -> unzip -> sdkmanager --licenses -> sdkmanager "platform-tools"
```

The SHA256 of the command-line tools zip is checked strictly, and a mismatch fails the build immediately. The SDK is installed to `/usr/local/android-sdk` using the `cmdline-tools/latest` layout. Only `platform-tools` are installed here. That keeps the runtime image useful for device interaction without turning it into a full Android build environment; the heavier Android build toolchain stays in the devcontainer and CI.

### Stage 3: Runtime stage

Final layer. Installs only three packages: `iptables`, `iproute2`, `iputils-ping`.

Copies in:

- JRE from the android-sdk stage.
- Android SDK from the android-sdk stage.
- `/usr/src/app/analyst/mesh` from the build stage to `/usr/bin/mesh`.
- `analyst/entrypoint.sh` to `/usr/bin/entrypoint.sh`.
- `analyst/amneziawg.conf.example` to `/etc/mesh/amneziawg.conf.example`.

### `.bashrc` welcome screen

A multi-line `RUN printf ... >> /root/.bashrc` injects a `meshcli` alias and an interactive-shell-only welcome banner listing subcommands. The `if [[ $- == *i* ]]` guard ensures non-interactive invocations (such as `docker compose exec analyst bash -c '...'`) don't print the banner into piped output.

### Entrypoint

`analyst/entrypoint.sh` (bash, `set -eo pipefail`) is responsible for the ephemeral node behaviour:

1. Fails fast if `LOGIN_URL` or `AUTH_KEY` is unset.
2. Generates a random hostname of 6 to 15 `[a-z0-9]` characters via `/dev/urandom`.
3. Launches the mesh daemon in the background with `--statedir=/root/.tailscale`, tee'd to `mesh.log`.
4. Runs `mesh cli up` with the generated hostname, `--accept-dns=false`, and a 10-second timeout, tee'd to `meshcli.log`.
5. `exec tail -f` on both log files so Compose keeps the container alive and `docker logs` shows daemon output.

Because the hostname is regenerated on every container start, reconnects show up as new node names. Whether old entries disappear automatically depends on the control-plane configuration and the kind of pre-auth key in use. With ephemeral keys and the current Headscale config, stale entries age out after `ephemeral_node_inactivity_timeout`.

---

## `control-plane/mesh-ui/Dockerfile`

Path: `control-plane/mesh-ui/Dockerfile`. Two stages.

### Stage 1: `node:22-alpine` (build)

```dockerfile
FROM node:22-alpine AS build
WORKDIR /app
RUN npm install -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build
```

`--frozen-lockfile` makes the build fail if `pnpm-lock.yaml` and `package.json` disagree, which catches accidental unversioned dep changes. `pnpm` is installed via `npm install -g` rather than via corepack or a feature flag.

### Stage 2: `nginx:alpine`

```dockerfile
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
```

The image ships only the built SPA. The nginx config is bind-mounted by Compose (`control-plane/nginx.conf:/etc/nginx/nginx.conf.template`) and substituted at container start, so config changes apply without rebuilding the image.

---

## `control-plane/nginx.conf`

### Location routing

- `/(api|register|key|machine|derp|verify|health|ts2021|swagger)` proxies to `http://headscale:8080` (service-name DNS).
- `/` in dev mode (if `$should_proxy` is set) proxies to `http://localhost:3000` (Vite).

### Security headers

Set on every response, always:

- A Content-Security-Policy allowing `'self'`, inline styles and scripts, and unpkg.com for SwaggerUI scripts/styles.
- `X-Content-Type-Options: nosniff`.
- `X-Frame-Options: DENY`.
- `X-XSS-Protection: 1; mode=block`.

---

## Dependencies and assumptions

- **Docker Engine with Compose V2 support**. The stack expects the modern `docker compose` plugin rather than legacy `docker-compose`.
- **Host kernel** must expose `/dev/net/tun`. On Linux this is standard. On a stripped-down kernel you'll need `modprobe tun`.
- **Host networking** must allow inbound on the `ui` service's bound port (defaults to `127.0.0.1:80`). A reverse proxy or tunneling service terminates TLS.
- **Workspace folder path**: `${LOCAL_WORKSPACE_FOLDER:-.}` is set by the devcontainer (`remoteEnv` in `devcontainer.json`) so bind mounts resolve to the host path, not the path inside the devcontainer. Outside the devcontainer, it falls back to the repo root directory.

---

← [Overview](index.md) | [Build orchestration](build-orchestration.md) →
