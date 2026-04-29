# Build troubleshooting

Build and container-specific issues. For runtime network, auth, and peer connectivity problems, see [Reference Troubleshooting](../troubleshooting.md).

## Compose stack won't start

### `docker compose build` fails pulling a base image

**Symptoms:**

- `ERROR: failed to pull image ...: manifest unknown` or similar during `task build`.

**Causes:**

- A pinned base image tag was removed from Docker Hub.
- Docker Hub is rate-limiting.

**Solutions:**

1. **Confirm connectivity**

    ```bash
    docker pull <image-from-the-error>
    ```

2. **If rate-limited**: authenticate with `docker login`. The free tier raises the anonymous rate limit.
3. **If the pinned image reference is actually gone**: update the relevant Dockerfile in a PR and review any related CI or devcontainer pins at the same time.

### `analyst` container exits immediately

**Symptoms:**

- `docker compose up analyst` starts then stops.
- `docker compose logs analyst` shows `Set the LOGIN_URL and AUTH_KEY environment variables in .env to run the MESH analyst client.`

**Solutions:**

1. Run `task analyst` (not `docker compose up analyst`). The Task wrapper prompts for missing values.
2. Or edit `.env` directly. Set both `LOGIN_URL` and `AUTH_KEY`.

### `analyst` container starts but tunnel doesn't come up

**Symptoms:**

- Logs show `tun device creation failed` or `operation not permitted`.

**Causes:**

- `/dev/net/tun` is not available on the host.
- The container isn't getting `NET_ADMIN` or `NET_RAW`.
- User namespace remapping (such as rootless Docker) is blocking the capability.

**Solutions:**

1. **Check `/dev/net/tun`**

    ```bash
    ls -l /dev/net/tun
    lsmod | grep tun
    ```

2. **Check for rootless Docker**

    ```bash
    docker info | grep -i 'userns\|rootless'
    ```

3. Rootless Docker requires [additional capabilities](https://docs.docker.com/engine/security/rootless/#networking-errors) to pass `/dev/net/tun` through. The simplest fix is to run the analyst under normal Docker. The broader fix is documented upstream.

### `ui` container: `envsubst: error while loading shared libraries`

**Symptoms:**

- `ui` container fails to start with a dynamic linker error about `envsubst`.

**Cause:**

- You've replaced the `nginx:alpine` base with a non-Alpine variant that doesn't include `gettext` by default.

**Solutions:**

1. Revert the base image.
2. Or add `gettext` to the Dockerfile's runtime stage. `envsubst` is in the `gettext` package.

### `ui` container: `502 Bad Gateway` on `/api/*`

**Symptoms:**

- The UI loads, but API calls return 502.

**Causes:**

- `headscale` service isn't running.
- `ui` was started with `--no-deps`, bypassing `depends_on`.
- The Compose project name doesn't match because the stack was started from two different directories.

**Solutions:**

1. **Check service state**

    ```bash
    docker compose ps
    docker compose logs headscale
    ```

2. **Check DNS inside `ui`**

    ```bash
    docker compose exec ui sh -c 'getent hosts headscale'
    ```

3. **Restart from the repo root**

    ```bash
    docker compose down
    docker compose up -d headscale ui
    ```

---

## Task (go-task) issues

### `task: command not found`

**Solution:**

Install go-task. See [taskfile.dev/installation](https://taskfile.dev/installation). On Ubuntu:

```bash
sudo snap install task --classic
```

### Prompts stuck in a loop

**Symptoms:**

- `task controlPlane` asks for the same question twice or loops.

**Cause:**

- The prompt writes to `.env` via `echo >> .env`. If the task fails between prompt and append, inspect `.env` for duplicates or stale entries.

**Solutions:**

1. **Inspect `.env`**

    ```bash
    cat .env
    ```

    Remove duplicate or empty lines.

2. **Or start over**

    ```bash
    rm .env
    task controlPlane
    ```

### `task tidy` modifies files unexpectedly

**Symptoms:**

- Fresh checkout followed by `task tidy` modifies `go.sum`.

**Causes:**

- Genuine dependency drift.
- Your Go version differs from the one in `go.mod`'s `toolchain:` directive.

**Solutions:**

1. **Check your Go version**

    ```bash
    go env GOVERSION GOTOOLCHAIN
    cat go.mod | grep -E '^(go|toolchain)'
    ```

2. If the Go versions don't match, switch to the repo's declared version. Otherwise CI and your local tidy will continually diverge.

### `Taskfile.dev.yml` task not found

**Symptoms:**

- `task buildAnalyst` says "task not found".

**Cause:**

- Developer tasks live in `Taskfile.dev.yml`, not the default `Taskfile.yml`.

**Solution:**

Use the `-t` flag:

```bash
task -t Taskfile.dev.yml buildAnalyst
```

---

## Dockerfile build failures

### `go mod download` fails in analyst build stage

**Symptoms:**

- First-stage build fails with `cannot find module` or a hash mismatch.

**Cause:**

- The build context is wrong. `analyst/Dockerfile` requires `context: .` at the repo root so it can see both `go.mod` and the `analyst/src` directory. `docker build -f analyst/Dockerfile analyst/` would only see `analyst/`.

**Solution:**

Always build via `task build`.

### Android SDK stage: `sha256sum --strict --check` fails

**Symptoms:**

- Build fails with `sha256sum: WARNING: 1 computed checksum did NOT match`.

**Cause:**

- Google changed the command-line tools archive behind the pinned filename without changing the version identifier.

**Solutions:**

1. **Read the current pinned archive name from the source of truth**

    Check `ANDROID_TOOLS_VERSION` in `analyst/Dockerfile` or `.devcontainer/Dockerfile`, depending on which build failed.

2. **Manually download and verify that archive**

    ```bash
    wget https://dl.google.com/android/repository/commandlinetools-linux-<ANDROID_TOOLS_VERSION>_latest.zip
    sha256sum commandlinetools-linux-<ANDROID_TOOLS_VERSION>_latest.zip
    ```

3. If it actually changed, open a PR updating the checksum in every place that installs that archive. In this repo that usually means both `analyst/Dockerfile` and `.devcontainer/Dockerfile`.

### `control-plane/mesh-ui` build: `pnpm install --frozen-lockfile` fails

**Symptoms:**

- `ERR_PNPM_OUTDATED_LOCKFILE`.

**Cause:**

- Someone changed `package.json` without running `pnpm install` to update `pnpm-lock.yaml`.

**Solutions:**

```bash
cd control-plane/mesh-ui
pnpm install
git add pnpm-lock.yaml
```

---

## Devcontainer issues

### Devcontainer build fails on Java install

**Symptoms:**

- `error: unsupported architecture: 'riscv64'` or similar during the JDK stage.

**Cause:**

- The Dockerfile only knows about `amd64` and `arm64` (the two JDK binaries Microsoft publishes). Other host architectures fail.

**Solutions:**

1. Use a `linux/amd64` or `linux/arm64` host. On Apple Silicon, VS Code's devcontainer extension builds an `arm64` variant automatically.
2. On a less common host, you need a different JDK source. Modify the Dockerfile, or don't use the devcontainer.

---

## Release script failures

### `working tree dirty after 'task tidy' + 'task generate'`

**Symptoms:**

- The script aborts at the reproducibility check.

**Causes:**

- `go.mod` or `go.sum` has changed.
- The tailscale submodule doesn't match upstream.

**Solutions:**

1. The script prints `git status --short`. If `go.mod` or `go.sum` is listed, commit the tidy changes in a separate PR first.
2. If files under `tailscale/` are listed, the submodule HEAD is wrong:

    ```bash
    git ls-tree HEAD tailscale
    git -C tailscale rev-parse HEAD
    ```

    The two hashes must match.

### `go tool yq is not available or not v4`

**Cause:**

- yq is a Go tool declared in `go.mod` (see the `tool` directive). On a fresh clone it hasn't been built yet.

**Solution:**

```bash
go mod download
go tool yq --version
```

---

## CI workflow failures

### `pin-github-actions.yml` fails with "Some github actions versions need pinning"

**Cause:**

- A PR added or modified a workflow and included a non-SHA action reference.

**Solution:**

Run `task pinGithubActions` locally and commit the result:

```bash
task -t Taskfile.dev.yml pinGithubActions
git add .github/workflows
git commit -m "ci: pin github actions"
```

### `govulncheck` fails on a submodule path

**Cause:**

- MESH is using a tailscale package that contains (or transitively imports) a known-vulnerable Go stdlib function.

**Solutions:**

1. **Upgrade Go**: bump `go` in `go.mod` (and the `toolchain:` line if pinned).
2. **Upgrade the vulnerable dependency**: if it's a direct dependency, run `go get <pkg>@latest && go mod tidy`.
3. **Patch the tailscale submodule**: if the vulnerable code lives there, the fix goes in `github.com/BARGHEST-ngo/mesh-tailscale`.

---

## Getting help

If you can't resolve your issue:

1. **Runtime/network issues**: see [Reference Troubleshooting](../troubleshooting.md).
2. **File an issue** with the log output from `docker compose logs` and the failing command: [github.com/BARGHEST-ngo/MESH/issues](https://github.com/BARGHEST-ngo/MESH/issues).
3. **Supply-chain concerns** (vulnerability, key compromise, SLSA failure): see [`SECURITY.md`](https://github.com/BARGHEST-ngo/MESH/blob/main/SECURITY.md) for private disclosure.

---

← [CI and release](ci-and-release.md)
