#!/usr/bin/env bash

set -eo pipefail

MESH_STATE_DIR="/home/mesh/.tailscale"

if [ -z "${LOGIN_URL}" ] || [ -z "${AUTH_KEY}" ]; then
    echo "Set the LOGIN_URL and AUTH_KEY environment variables in .env to run the MESH analyst client." >&2
    exit 1
fi

/usr/bin/mesh --tun=userspace-networking --verbose=1 \
    --socket="${MESH_STATE_DIR}/tailscaled.sock" \
    --statedir="${MESH_STATE_DIR}" 2>&1 | tee "${MESH_STATE_DIR}/mesh.log" &

/usr/bin/meshcli up \
    --login-server="${LOGIN_URL}" \
    --auth-key="${AUTH_KEY}" \
    --accept-dns=false \
    --timeout=10s
     2>&1 | tee "${MESH_STATE_DIR}/meshcli.log"

exec tail -f "${MESH_STATE_DIR}/mesh.log" "${MESH_STATE_DIR}/meshcli.log"