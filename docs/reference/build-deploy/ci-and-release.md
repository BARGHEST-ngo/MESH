# CI and release

Every MESH release job runs on a freshly provisioned GitHub-hosted runner. This page documents the Android release flow, its guarantees, and its inputs. Exact action SHAs, cache keys, and toolchain pins live in the workflow and script files that enforce them.

## Release workflow

The documented release path has two entry points:

- Tag pushes matching `v*` for real releases.
- `workflow_dispatch` for maintainers who need to rebuild a specific version string manually.

Pull requests touching the Android build inputs also exercise the same build logic in unsigned mode so reproducibility problems show up before release.

## `android-release.yml`

Path: `.github/workflows/android-release.yml`. Four jobs cooperate here: `build-release`, `verify-reproducibility`, `create-release`, and `provenance`.

### Triggers

- Tag pushes matching `v*` run the full release path.
- Pull requests touching the Android build inputs run an unsigned build plus reproducibility verification.
- `workflow_dispatch` lets a maintainer rebuild a specific version string without moving tags.

### Job 1: `build-release`

Runs on `ubuntu-latest` with `android-client` as the working directory. It is the primary build job:

1. Checks out the repo and sets up Go from `go.mod`.
2. Installs the Java and Android toolchain pinned in the workflow.
3. On tag/manual release runs, decodes the signing keystore from secrets and builds the signed APK. On PRs, builds the unsigned APK only.
4. Hashes the unsigned APK and exposes that hash as a job output for the reproducibility job.
5. Derives a version string from the tag, manual input, or PR number. Non-PR runs then rename the artifact, compute the provenance subject hash, and upload the APK artifact.

The outputs are the unsigned APK hash for reproducibility verification and the base64-encoded release artifact hash for the provenance job.

### Job 2: `verify-reproducibility`

Runs after `build-release`, on a separate `ubuntu-latest` runner. It installs the same pinned toolchain, rebuilds the unsigned APK, and compares its SHA-256 hash to the one exported by `build-release`.

**What this proves:** the unsigned APK is reproducible across two clean runners. If the hashes differ, the release path introduced nondeterminism somewhere in source generation, filesystem ordering, timestamps, or toolchain behavior.

**Why unsigned?** The unsigned APK is the deterministic artifact. Signing adds a signature block and related metadata, so the signed APK is expected to differ even when the unsigned payload is reproducible.

### Job 3: `create-release`

Runs only on `push` of a tag. Downloads the APK artifact from `build-release`, creates a GitHub Release with `generate_release_notes: true`, and attaches the APK.

### Job 4: `provenance`

Also runs only on tag push. Calls the reusable workflow:

```yaml
uses: slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@<pinned-release-ref>
with:
  base64-subjects: ${{ needs.build-release.outputs.hashes }}
  upload-assets: true
```

This produces SLSA Level 3 provenance and attaches it to the GitHub Release alongside the APK. It runs in a separate reusable workflow, which reduces the chance that the primary build job can tamper with provenance generation. The repo's action-pinning rule intentionally excludes this reusable workflow because SLSA verification expects a release-style ref. See [Build and release safeguards](#build-and-release-safeguards) below.

### Secrets

- `KEYSTORE_BASE64` becomes a temporary `mesh-release.jks` file for the signing step.
- `KEYSTORE_PASSWORD` and `KEYSTORE_ALIAS` are passed as environment variables to Gradle.
- The keystore file is removed unconditionally before job exit.

None of the secret values are written intentionally to logs, and GitHub masks them automatically.

## Release script

Path: `scripts/tag-release.sh`. Invoked via `task -t Taskfile.dev.yml tagRelease VERSION=<semver>`.

### Preconditions

The script fails unless all of the following are true:

1. The requested version is strict semver.
2. Required tooling is installed, including `go tool yq`.
3. The working tree is clean.
4. The current branch is `main`.
5. The `tailscale/` submodule is clean and matches the parent repo's recorded SHA.
6. Neither the main-repo tag nor the submodule tag already exists locally.
7. If remote checks are enabled, local `main` matches `origin/main` and the tag does not already exist on the remote.

### Version bump

The script intentionally modifies only three files:

- `android-client/android/build.gradle` for the Android `versionCode`.
- `android-client/barghest.version` for the human-readable version fields.
- `android-client/fdroid/com.barghest.mesh.yml` for the new F-Droid build entry and current-version pointers.

### Commit and tags

The script creates one release-preparation commit plus two annotated tags: the main repo tag (`vX.Y.Z`) and the tailscale submodule tag (`mesh-vX.Y.Z`).

**Submodule tag prefix.** The submodule tag is `mesh-vX.Y.Z`, not `vX.Y.Z`. The tailscale submodule's remote is shared with upstream Tailscale, so the `mesh-` prefix prevents a tag collision.

---

## Build and release safeguards

- **SLSA L3 provenance** on every Android release, generated on an isolated runner.
- **Reproducibility verification** on every release and on PRs that modify Android build inputs.
- **Pinned Java and Android toolchains** in both CI and the devcontainer, with exact values kept in the workflow/Dockerfile sources.

---

← [Build orchestration](build-orchestration.md) | [Troubleshooting](troubleshooting.md) →
