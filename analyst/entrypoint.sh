#!/usr/bin/env bash

set -eo pipefail

MESH_STATE_DIR="/root/.tailscale"

if [ -z "${LOGIN_URL}" ] || [ -z "${AUTH_KEY}" ]; then
    echo "Set the LOGIN_URL and AUTH_KEY environment variables in .env to run the MESH analyst client." >&2
    exit 1
fi

HOSTNAME_LEN=$(( (RANDOM % 10) + 6 ))
RANDOM_HOSTNAME=$(LC_ALL=C head -c 512 /dev/urandom | tr -dc 'a-z0-9')
RANDOM_HOSTNAME="${RANDOM_HOSTNAME:0:$HOSTNAME_LEN}"

/usr/bin/tailscaled \
    --statedir="${MESH_STATE_DIR}" \
    2>&1 | tee "${MESH_STATE_DIR}/mesh.log" &

/usr/bin/tailscale up \
    --login-server="${LOGIN_URL}" \
    --auth-key="${AUTH_KEY}" \
    --hostname="${RANDOM_HOSTNAME}" \
    --accept-dns=false \
    --timeout=10s \
    2>&1 | tee "${MESH_STATE_DIR}/meshcli.log"

exec tail -f "${MESH_STATE_DIR}/mesh.log" "${MESH_STATE_DIR}/meshcli.log"
