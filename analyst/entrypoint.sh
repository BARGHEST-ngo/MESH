#!/bin/sh

if [ -z "${LOGIN_URL}" ] || [ -z "${AUTH_KEY}" ]; then
    echo "Set the LOGIN_URL and AUTH_KEY environment variables in .env to run the MESH analyst client."
    exit 1
fi

/usr/bin/mesh --tun=userspace-networking --verbose=1 --socket=/var/run/tailscale/tailscaled.sock --statedir=/var/run/tailscale &> /var/run/tailscale/mesh.log &

sleep 2

/usr/bin/meshcli up --login-server=${LOGIN_URL} --auth-key=${AUTH_KEY} --accept-dns=false &> /var/run/tailscale/meshcli.log

exec "$@"