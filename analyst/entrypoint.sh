#!/usr/bin/env bash

set -eo pipefail

# Fix volume perms (old deployments are root-owned) then drop to mesh,
# keeping NET_ADMIN/NET_RAW ambient so iptables etc. still work.
if [ "$(id -u)" = "0" ]; then
    if ! su -s /bin/sh mesh -c 'test -w /home/mesh/acquisitions'; then
        chown mesh:mesh /home/mesh/acquisitions ||
            echo "warning: cannot chown /home/mesh/acquisitions; acquisitions may not be writable" >&2
    fi
    export HOME=/home/mesh USER=mesh
    exec setpriv \
        --reuid=mesh --regid=mesh --init-groups \
        --inh-caps=+net_admin,+net_raw \
        --ambient-caps=+net_admin,+net_raw \
        -- "$0" "$@"
fi

# Placeholder until new MESHnet-based commands are in
if [ "$#" -gt 0 ]; then
    exec "$@"
fi
exec sleep infinity