#!/usr/bin/env bash
set -euo pipefail

# Parse barghest.version as data so it's not executable
# each value to match allow list
VERSION_FILE="$(dirname "$0")/barghest.version"
[[ -f "$VERSION_FILE" ]] || { echo >&2 "no barghest.version file found"; exit 1; }

parse_version_field() {
    local key="$1"
    local value
    value="$(grep -E "^${key}=" "$VERSION_FILE" | head -n1 | cut -d= -f2- | tr -d '"')"
    if ! [[ "$value" =~ ^[0-9A-Za-z.+-]+$ ]]; then
        echo >&2 "invalid ${key} value in barghest.version"
        exit 1
    fi
    printf '%s' "$value"
}

VERSION_SHORT="$(parse_version_field VERSION_SHORT)"
VERSION_LONG="$(parse_version_field VERSION_LONG)"
GIT_HASH="$(git rev-parse HEAD 2>/dev/null || echo unknown)"
echo "-X tailscale.com/version.longStamp=${VERSION_LONG}"
echo "-X tailscale.com/version.shortStamp=${VERSION_SHORT}"
echo "-X tailscale.com/version.gitCommitStamp=${GIT_HASH}"
