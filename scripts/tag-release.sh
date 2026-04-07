#!/usr/bin/env bash
# scripts/tag-release.sh — prepare a MESH release commit and tags.
#
# Usage: scripts/tag-release.sh <version>
#   <version> is strict semver, e.g. 0.1.5-alpha or 1.0.0
#
# What it does:
#   1. Validates preconditions (clean tree, on main, tag doesn't exist, ...)
#   2. Runs `task tidy` and `task generate` as reproducibility invariant checks
#   3. Bumps versionCode via `make -C android-client bump_version_code`
#   4. Rewrites android-client/barghest.version (5 lines, no git hash —
#      the hash is now computed at build time by version-ldflags.sh)
#   5. Appends a new Builds entry to android-client/fdroid/com.barghest.mesh.yml
#      (F-Droid convention: newest at bottom) and updates CurrentVersion/
#      CurrentVersionCode
#   6. Creates a single commit, an annotated tag on the parent repo, and an
#      annotated tag on the tailscale submodule (prefixed `mesh-` to avoid
#      collision with upstream Tailscale tags)
#   7. Prints push instructions — does NOT push
#
# Environment:
#   MESH_SKIP_REMOTE_CHECK=1  skip the network-based remote check (offline mode)
#
# Prerequisites:
#   - Go toolchain (yq is declared as a Go tool dependency in go.mod and
#     invoked via `go tool yq`; first use builds and caches the binary)
#   - Standard Unix tools: git, make, awk, grep, printf
#   - The `task` CLI (go-task) — already required to invoke via Taskfile

set -euo pipefail

# ---- Argument parsing --------------------------------------------------------

if [ $# -ne 1 ] || [ -z "${1:-}" ]; then
  echo "Usage: $0 <version>" >&2
  echo "Example: $0 0.1.6-alpha" >&2
  exit 1
fi

BARE_VERSION="${1#v}"   # strip leading "v" if present
TAG="v${BARE_VERSION}"
SUBMODULE_TAG="mesh-v${BARE_VERSION}"

# ---- Resolve repo root (so the script works regardless of cwd) --------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

VERSION_FILE="android-client/barghest.version"
GRADLE_FILE="android-client/android/build.gradle"
FDROID_YML="android-client/fdroid/com.barghest.mesh.yml"

# ---- Helpers -----------------------------------------------------------------

die() {
  echo "ERROR: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "'$1' is required but not found on PATH"
}

# ---- Preconditions -----------------------------------------------------------

# 1. Strict semver (no leading zeros, no build metadata)
if ! printf '%s' "$BARE_VERSION" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$'; then
  die "VERSION must be strict semver (e.g. 0.1.6-alpha, 1.0.0). Got: '$1'"
fi

# 2. Required tools
require_cmd git
require_cmd make
require_cmd awk
require_cmd grep
require_cmd task
require_cmd go
go tool yq --version 2>&1 | grep -qE 'version v?4\.' \
  || die "go tool yq is not available or not v4. yq is declared in go.mod's tool directive — run 'go mod download' if this is a fresh clone."

# 3. Working tree clean (ignore untracked so stray dotfiles don't block)
[ -z "$(git status --porcelain --untracked-files=no)" ] \
  || die "working tree has uncommitted changes. Commit or stash first."

# 4. Current branch is main
[ "$(git rev-parse --abbrev-ref HEAD)" = "main" ] \
  || die "must be on main branch to tag a release (got: $(git rev-parse --abbrev-ref HEAD))"

# 5. Submodule clean and at registered SHA
[ -z "$(git -C tailscale status --porcelain)" ] \
  || die "tailscale/ submodule has uncommitted changes"

REG_SHA="$(git ls-tree HEAD tailscale | awk '{print $3}')"
HEAD_SHA="$(git -C tailscale rev-parse HEAD)"
[ "$REG_SHA" = "$HEAD_SHA" ] \
  || die "tailscale/ submodule HEAD ($HEAD_SHA) does not match parent's registered SHA ($REG_SHA). Run: git submodule update tailscale"

# 6. Neither tag exists locally
! git rev-parse --verify "refs/tags/$TAG" >/dev/null 2>&1 \
  || die "tag $TAG already exists locally. Delete with: git tag -d $TAG"
! git -C tailscale rev-parse --verify "refs/tags/$SUBMODULE_TAG" >/dev/null 2>&1 \
  || die "submodule tag $SUBMODULE_TAG already exists locally"

# 7. Optional remote check (opt out with MESH_SKIP_REMOTE_CHECK=1)
if [ "${MESH_SKIP_REMOTE_CHECK:-0}" != "1" ]; then
  git fetch origin main --quiet || die "git fetch origin main failed"
  [ "$(git rev-parse HEAD)" = "$(git rev-parse origin/main)" ] \
    || die "local main differs from origin/main. Sync first, or set MESH_SKIP_REMOTE_CHECK=1"
  ! git ls-remote --tags origin "refs/tags/$TAG" 2>/dev/null | grep -q . \
    || die "tag $TAG already exists on origin"
fi

# ---- Reproducibility invariant checks ---------------------------------------

echo "==> Running tidy + generate as reproducibility invariant checks..."
task tidy
task generate
if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
  echo "ERROR: working tree dirty after 'task tidy' + 'task generate'." >&2
  echo "  Either go.mod/go.sum drifted, or the tailscale submodule no longer" >&2
  echo "  matches upstream + patches. Resolve before tagging:" >&2
  git status --short >&2
  exit 1
fi

# ---- Perform the bump -------------------------------------------------------

echo "==> Bumping version to $BARE_VERSION"

# Parse version components
MAJOR="$(printf '%s' "$BARE_VERSION" | cut -d. -f1)"
MINOR="$(printf '%s' "$BARE_VERSION" | cut -d. -f2)"
PATCH="$(printf '%s' "$BARE_VERSION" | cut -d. -f3 | cut -d- -f1)"
SHORT="${MAJOR}.${MINOR}.${PATCH}"
LONG="$BARE_VERSION"

# -- Bump versionCode via the existing Makefile target
make -C android-client bump_version_code >/dev/null
NEW_CODE="$(awk '/versionCode [0-9]+/{print $2; exit}' "$GRADLE_FILE")"
[ -n "$NEW_CODE" ] || die "could not read new versionCode from $GRADLE_FILE"
echo "    versionCode bumped to $NEW_CODE"

# -- Rewrite barghest.version (5 lines; no VERSION_GIT_HASH — computed at build
#    time by version-ldflags.sh)
{
  printf 'VERSION_MAJOR=%s\n' "$MAJOR"
  printf 'VERSION_MINOR=%s\n' "$MINOR"
  printf 'VERSION_PATCH=%s\n' "$PATCH"
  printf 'VERSION_SHORT="%s"\n' "$SHORT"
  printf 'VERSION_LONG="%s"\n' "$LONG"
} > "$VERSION_FILE"
echo "    barghest.version rewritten (VERSION_LONG=$LONG)"

# -- Update fdroid yml via yq: APPEND a new Builds entry (F-Droid convention:
#    oldest at top, newest at bottom). Copy from the current last entry and
#    override versionName/versionCode/commit. Uses strenv() to avoid shell
#    quoting issues and `| tonumber` so versionCode emits as a YAML number.
#    The commit field uses the tag name — F-Droid resolves it at build time.
export FDROID_VERSION_NAME="$LONG"
export FDROID_VERSION_CODE="$NEW_CODE"
export FDROID_COMMIT_REF="$TAG"
go tool yq eval -i '
  .Builds = .Builds + [(.Builds[-1] * {
    "versionName": strenv(FDROID_VERSION_NAME),
    "versionCode": (strenv(FDROID_VERSION_CODE) | tonumber),
    "commit": strenv(FDROID_COMMIT_REF)
  })] |
  .CurrentVersion = strenv(FDROID_VERSION_NAME) |
  .CurrentVersionCode = (strenv(FDROID_VERSION_CODE) | tonumber)
' "$FDROID_YML"
echo "    fdroid yml: appended Builds entry (commit: $TAG), updated CurrentVersion"

# -- Sanity: only the three expected files should be dirty
UNEXPECTED="$(git status --porcelain --untracked-files=no \
  | grep -vE ' (android-client/barghest\.version|android-client/android/build\.gradle|android-client/fdroid/com\.barghest\.mesh\.yml)$' \
  || true)"
if [ -n "$UNEXPECTED" ]; then
  echo "ERROR: unexpected modified files after bump:" >&2
  echo "$UNEXPECTED" >&2
  exit 1
fi

# ---- Commit and tag ----------------------------------------------------------

echo "==> Creating commit and tags"

git add "$VERSION_FILE" "$GRADLE_FILE" "$FDROID_YML"
git commit -m "chore(android): release $TAG" -m "- Bump versionCode to $NEW_CODE in build.gradle
- Update barghest.version to $BARE_VERSION
- Append $TAG entry to fdroid Builds list
- Update CurrentVersion and CurrentVersionCode"

git tag -a "$TAG" -m "Version $BARE_VERSION"
git -C tailscale tag -a "$SUBMODULE_TAG" -m "MESH $BARE_VERSION"

# ---- Success output ----------------------------------------------------------

cat <<EOF

Release $TAG prepared locally.

To publish, push in this order:
  1. git -C tailscale push git@github.com:BARGHEST-ngo/mesh-tailscale.git $SUBMODULE_TAG
  2. git push origin main
  3. git push origin $TAG

Step 3 triggers .github/workflows/android-release.yml: it builds the APK,
verifies reproducibility, signs it, publishes to GitHub Releases, and generates
SLSA L3 provenance. F-Droid auto-picks up the new tag via fdroiddata's own
update mechanism.

To abort before pushing:
  git reset --hard HEAD~1
  git tag -d $TAG
  git -C tailscale tag -d $SUBMODULE_TAG

EOF
