# Build and deploy

This section documents MESH's containerization strategy and its build and release pipelines. It's written for analysts, contributors, and anyone auditing the supply chain. For end-user setup, see the [Setup guide](../../setup/index.md).

## What's in this section

- **[Runtime stack](runtime-stack.md)**: the `compose.yml` stack and every `Dockerfile` that ships in the repo.
- **[Build orchestration](build-orchestration.md)**: the Task-based build system (`Taskfile.yml`, `Taskfile.dev.yml`), the analyst build script, and the VS Code devcontainer.
- **[CI and release](ci-and-release.md)**: the Android release workflow and `scripts/tag-release.sh`.
- **[Build troubleshooting](troubleshooting.md)**: build and container-layer issues. For network, auth, and runtime problems, see [Reference Troubleshooting](../troubleshooting.md).

## Key design decisions

### 1. Docker Compose for the runtime, go-task for orchestration

`compose.yml` defines three services (`headscale`, `ui`, `analyst`) and is invoked via `docker compose` directly. All orchestration lives in two Taskfiles rather than `make` or shell scripts. This includes interactive setup prompts, `.env` management, templated config generation, and release tagging. [go-task](https://taskfile.dev) is a better fit than `make` for this job: tasks are declarative, dependencies between them are explicit, and the built-in `prompt:` directive handles interactive setup nicely.

### 2. Multi-stage Dockerfiles with pinned base images

Every `Dockerfile` uses multi-stage builds and pins its base images in the source files that actually consume them. Build-time downloads such as the Android command-line tools are checksum-verified where the Dockerfiles fetch them.

### 3. Ephemeral-first analyst container

The analyst container generates a random hostname on every start. It accepts credentials only via `LOGIN_URL` and `AUTH_KEY` environment variables, and stores per-session state in a named Docker volume. Rotating the auth key with `task setAuthKey` stops the analyst and removes that volume before prompting for the new key. No long-lived identity is baked into the image.

### 4. Reproducible Android releases

The Android release pipeline rebuilds the unsigned APK on a second, isolated runner and fails the release if the SHA-256 hashes differ. SLSA Level 3 provenance is generated via the `slsa-github-generator` reusable workflow. Go comes from `go.mod`, and the Java/Android toolchain is pinned in the Dockerfiles and workflow files that install it.

Reproducibility is a security benefit on its own. F-Droid has argued that VPN apps require an unusually high trust baseline that reproducible builds help establish. See F-Droid's [VPN trust requires Free Software](https://f-droid.org/en/2023/03/08/vpn-trust-requires-free-software.html) post.

### 5. Supply-chain hardening

- GitHub Actions are pinned to commit SHAs via [frizbee](https://github.com/stacklok/frizbee), refreshed on demand with `task pinGithubActions`.
- Go vulnerability scanning (`govulncheck`) and linting (`golangci-lint`) run on pushes to `main` and on PRs that modify Go sources.
- The tailscale submodule is tagged separately with a `mesh-v*` prefix to prevent collisions with upstream Tailscale tags in a shared remote.

### 6. Config via `.env`, never committed

All runtime configuration (`CONTROL_PLANE_DOMAIN`, `LOGIN_URL`, `AUTH_KEY`, `CORS_ORIGIN`, `NGINX_MODE`) is read from `.env` at the repo root. The file is created on first `task` run with `chmod 600`. There is no committed `.env.example`. Instead, the Taskfile prompts define the variables.

## Assumptions

- You have Docker Engine with the Compose V2 plugin (`docker compose`, not `docker-compose`).
- You have the [go-task](https://taskfile.dev) CLI installed. Most tasks work without it via plain `docker compose`, but interactive setup does not.
- You're on Linux, macOS, or WSL. Native Windows works for many tasks but is untested in CI.
- For Android work: Go matching `go.mod`, plus the Java/Android toolchain declared in `.devcontainer/Dockerfile` and `.github/workflows/android-release.yml`. Or just use the [devcontainer](build-orchestration.md#development-container).
- For release work: the `tailscale/` submodule is initialised. Clone with `git clone --recurse-submodules https://github.com/BARGHEST-ngo/MESH.git`, or on an existing clone run `git submodule update --init --recursive`.

---

**Next:** [Runtime stack](runtime-stack.md)
