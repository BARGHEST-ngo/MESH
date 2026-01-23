#!/usr/bin/env bash

source barghest.version || { echo >&2 "no barghest.version file found"; exit 1; }
if [[ -z "${VERSION_LONG}" ]]; then
    exit 1
fi
echo "-X tailscale.com/version.longStamp=${VERSION_LONG}"
echo "-X tailscale.com/version.shortStamp=${VERSION_SHORT}"
echo "-X tailscale.com/version.gitCommitStamp=${VERSION_GIT_HASH}"
echo "-X tailscale.com/version.extraGitCommitStamp=${VERSION_EXTRA_HASH}"
