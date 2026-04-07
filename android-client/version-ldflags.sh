#!/usr/bin/env bash

source barghest.version || { echo >&2 "no barghest.version file found"; exit 1; }
if [[ -z "${VERSION_LONG}" ]]; then
    exit 1
fi
GIT_HASH="$(git rev-parse HEAD 2>/dev/null || echo unknown)"
echo "-X tailscale.com/version.longStamp=${VERSION_LONG}"
echo "-X tailscale.com/version.shortStamp=${VERSION_SHORT}"
echo "-X tailscale.com/version.gitCommitStamp=${GIT_HASH}"
