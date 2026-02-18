#!/bin/bash

if [[ -z "$TOOLCHAIN_DIR" ]]; then
    # MESH fork: We use the Tailscale Go toolchain but disable the
    # "tailscale_go" build tag to avoid strict version validation
    # (version_checkformat.go) that requires Tailscale's version format.
    # MESH uses its own independent versioning scheme.
    echo "mesh_build"
else
    # Otherwise, if TOOLCHAIN_DIR is specified, we assume
    # we're F-Droid or something using a stock Go toolchain.
    # That's fine. But we don't set the tailscale_go build tag.
    # Return some no-op build tag that's non-empty for clarity
    # when debugging.
    echo "not_tailscale_go"
fi
