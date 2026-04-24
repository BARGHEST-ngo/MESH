# Build orchestration

MESH uses [go-task](https://taskfile.dev) for build orchestration and task running. Two Taskfiles split responsibilities:

- `Taskfile.yml` is user-facing (deploy-time).
- `Taskfile.dev.yml` is developer-facing (build-time).

Both auto-load `.env` and run silently by default.

## Why Task?

Three constraints shaped the choice:

- **Interactive prompts that work cross-platform.** The setup flow asks for the control plane type, domain, and auth key, then writes them to `.env`. Task's `prompt:` and `requires.vars.enum:` directives make this declarative. `make` would need a custom shell wrapper per platform.
- **Embedded platform-specific tasks.** Several tasks branch on OS (sed vs PowerShell, Linux sed vs macOS sed). Task exposes `platforms: [linux, darwin, windows]` on individual commands, so one Taskfile serves all three.
- **Dotenv loading.** `dotenv: [.env]` means every task sees configuration variables without sourcing scripts.

## `Taskfile.yml`: user-facing tasks

Path: `Taskfile.yml`. Default invocation target. All user-facing tasks live here.

### Task reference

| Task | Purpose | Depends on |
|---|---|---|
| `build` | `docker compose build` | |
| `controlPlane` | `docker compose up -d headscale ui` | `configureControlPlane` |
| `down` | `docker compose down` | |
| `apikey` | Calls `headscale apikeys create --expiration 3h` inside the running container. Prints a help message if the container isn't running locally. | |
| `analyst` | Brings up the analyst service and opens an interactive bash shell via `docker compose exec analyst /bin/bash`. | `configureAnalyst` |
| `setAuthKey` | Stops the analyst, removes its volume, and re-prompts for a new `AUTH_KEY`. Use when rotating keys or starting a fresh session. | |

### The configuration flow

`configureControlPlane` and `configureAnalyst` are high-level tasks that run the interactive setup prompts. Each prompt task is guarded by an `if:` so it only runs when the relevant `.env` variable is empty. Rerunning `task controlPlane` on a configured checkout skips the prompts.

```text
controlPlane
  └── configureControlPlane
        ├── controlPlaneType      (prompts Ephemeral vs Persistent)
        ├── controlPlaneDomain    (prompts domain name)
        ├── headscaleConfig       (sed-substitutes config.example.yaml if config.yaml missing)
        └── confirmDotEnv         (cat .env, prompt "looks correct?")
```

### Headscale config templating

`headscaleConfig` (internal) substitutes `CONTROL_PLANE_DOMAIN` in `config.example.yaml` and writes `config.yaml`, but only if the latter doesn't exist:

```yaml
- task: headscaleConfig
  if: '[ ! -f control-plane/headscale/config.yaml ]'
```

The task declares `sources:` and `generates:` so Task skips it if the source hasn't changed since the generate target's mtime. On Windows, a PowerShell `-replace` command does the same thing.

### Key rotation: `setAuthKey`

```yaml
cmds:
  - docker compose down analyst && docker volume rm $(docker compose volumes -q | grep analyst-data) || true
  - <prompt>
  - task: promptAnalystAuthKey
```

Stopping the analyst and deleting its volume before asking for the new key ensures the container can't come up with stale state between the prompt and the restart.

---

## `Taskfile.dev.yml`: developer-facing tasks

Path: `Taskfile.dev.yml`. Invoked via `task -t Taskfile.dev.yml <task>`. Holds every task that a contributor runs and an end-user should not.

### Task reference

| Task | What it does |
|---|---|
| `tidy` | `go mod tidy`. Kept separate because it mutates `go.sum` and affects reproducibility. |
| `generate` | `go generate ./...`. Depends on `submoduleInit`. |
| `docs` | `mkdocs build`. |
| `serveDocs` | `mkdocs serve` with live reload on `:8000`. |
| `buildAnalyst` | Runs `analyst/src/build.sh`. Produces `analyst/mesh`. Depends on `submoduleInit`. |
| `buildAndroid` | `make -C android-client apk`. |
| `buildAndroidRelease` | `make -C android-client release-apk`. Unsigned. |
| `buildAndroidSigned` | `make -C android-client release-apk-signed`. Requires `JKS_PATH`, `JKS_PASSWORD`, `JKS_ALIAS` env vars. |
| `clean` | `rm -f analyst/mesh` and `make -C android-client clean`. |
| `tagRelease` | Calls `scripts/tag-release.sh $VERSION`. See [CI and release](ci-and-release.md#release-script). |
| `fdroidInit` | Initialise a local F-Droid repo for reproducibility testing. |
| `fdroidBuild` | `cd fdroid && fdroid build -v -l com.barghest.mesh`. |
| `pinGithubActions` | `go tool github.com/stacklok/frizbee actions .github/workflows`. |

---

## Analyst build script

Path: `analyst/src/build.sh`. Called by the Dockerfile (stage 1) and by `task buildAnalyst`. Single source of truth for the mesh binary build.

### What it does

1. Locates the Go module directory via `go env GOMOD`.
2. Assembles build tags, always adding `ts_omit_logtail` (omits Tailscale telemetry) and any extra tags supplied through `$TAGS`.
3. `cd` into `$GO_MOD_DIR/tailscale` (the submodule).
4. Runs:

    ```bash
    go build -tags "$tags" -trimpath -o "$GO_MOD_DIR/analyst/mesh" ./cmd/tailscaled
    ```

### Why a separate script rather than inline in the Dockerfile?

Two reasons:

- `task buildAnalyst` runs it for local dev without Docker. Contributors iterating on CLI changes rebuild in seconds instead of rebuilding the whole container.
- Build flags live in one place. A change to build tags updates both the container image and local dev builds without drift.

---

## Development container

Path: `.devcontainer/Dockerfile`, `.devcontainer/devcontainer.json`. Defines an ephemeral development environment compatible with VS Code's Remote Containers and GitHub Codespaces.

### What the image installs

The image is built from a devcontainers Debian base image. It installs:

- Microsoft OpenJDK, Go, Node.js, and pnpm.
- The Android command-line tools and the SDK components needed for Android builds.
- Task CLI.

Java and the Android command-line tools are checksum-verified in the Dockerfile. Go is version-pinned via the `GO_VERSION` build arg. Android SDK licenses are auto-accepted with `yes | sdkmanager --licenses`.

### `devcontainer.json` features

- Docker-outside-of-Docker so `docker compose` inside the devcontainer drives the host Docker daemon.
- MkDocs tooling so `task docs` and `task serveDocs` work without extra setup.
- Gradle on `PATH` for Android builds.
- Python plus `fdroidserver` for local F-Droid reproducibility work.

### Mounts and environment

- Binds `~/.android` from the host to `/root/.android` to persist the Android debug keystore across rebuilds.
- Sets `LOCAL_WORKSPACE_FOLDER=${localWorkspaceFolder}` so bind mounts in `compose.yml` resolve to the host's repo path, not a path inside the devcontainer.
- Forwards port `80` so the UI is reachable from the host browser.

---

## Dependencies and assumptions

- **A `task` release with `prompt:` and `requires.vars.enum:` support**.
- **`go` on PATH** for `task tidy`, `task generate`, `task buildAnalyst`. The CI workflow sets this up via `actions/setup-go` reading `go-version-file: go.mod`.
- **`git submodule` initialised** before running any build task. `submoduleInit` handles this but requires network access on first run.
- **F-Droid work** requires `fdroidserver`, which requires Python, gradle, JDK, and the Android SDK. The devcontainer has all of this.
- **Signed releases** require `JKS_PATH`, `JKS_PASSWORD`, and `JKS_ALIAS` environment variables and an existing `.jks` keystore. `task buildAndroidSigned` fails if these are missing.

---

← [Runtime stack](runtime-stack.md) | [CI and release](ci-and-release.md) →
