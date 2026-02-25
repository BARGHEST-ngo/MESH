#!/bin/sh

if [ -z "${LOGIN_URL}" ] || [ -z "${AUTH_KEY}" ]; then
    echo "Set the LOGIN_URL and AUTH_KEY environment variables in the Docker container via .env to connect to the MESH on startup."
    exit 1
fi

/usr/bin/mesh --verbose=1 --socket=/var/run/tailscale/tailscaled.sock &

sleep 2

exec /usr/bin/meshcli up --login-server=${LOGIN_URL} --auth-key=${AUTH_KEY} --accept-dns=false
