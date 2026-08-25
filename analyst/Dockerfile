FROM golang:1.26.3-alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

RUN apk add --no-cache git patch

WORKDIR /usr/src/app

# Disable workspace so each module resolves independently
ENV GOWORK=off

# Cache root module deps (used by tailscale submodule build)
COPY go.mod go.sum ./
RUN go mod download

# Cache analyst module deps separately
COPY tailscale/go.mod tailscale/go.sum ./tailscale/
COPY analyst/go.mod analyst/go.sum ./analyst/
RUN cd analyst && go mod download

COPY analyst ./analyst
COPY tailscale ./tailscale
RUN ./analyst/src/build.sh

FROM debian:trixie-slim@sha256:4e401d95de7083948053197a9c3913343cd06b706bf15eb6a0c3ccd26f436a0e

RUN apt-get update && apt-get install -y --no-install-recommends \
        adb \
        iptables \
        iproute2 \
        iputils-ping \
        libcap2-bin \
        suricata \
        tcpdump \
        util-linux \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1000 mesh \
    && useradd  --system --uid 1000 --gid 1000 \
                --create-home --home-dir /home/mesh \
                --shell /bin/bash mesh \
    && mkdir -p /home/mesh/.tailscale \
    && mkdir -p /home/mesh/acquisitions \
    && chown -R mesh:mesh /home/mesh \
    && setcap cap_net_admin,cap_net_raw+eip /usr/bin/tcpdump

# MVT (Mobile Verification Toolkit) for offline analysis of adbcollect output:
# Installed into a venv so the Debian-managed python3 stays untouched,
# python3-venv is only needed to create the venv, so it is purged afterwards
# python3 itself must stay, as the venv symlinks back to it
RUN apt-get update && apt-get install -y --no-install-recommends \
        python3 \
        python3-venv \
    && python3 -m venv /opt/mvt \
    && /opt/mvt/bin/pip install --no-cache-dir mvt==2026.7.29 \
    && apt-get purge -y python3-venv \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /home/mesh/.local/share/mvt/indicators \
    && chown -R mesh:mesh /home/mesh/.local.

ENV PATH="/opt/mvt/bin:$PATH"

COPY --from=build       /usr/src/app/analyst/tailscaled  /usr/bin/tailscaled
COPY --from=build       /usr/src/app/analyst/tailscale   /usr/bin/tailscale
COPY --from=build       /usr/src/app/analyst/mesh        /usr/bin/mesh
COPY --chmod=755        analyst/entrypoint.sh            /usr/bin/entrypoint.sh
COPY --chown=mesh:mesh  analyst/bashrc                   /home/mesh/.bashrc
COPY                    analyst/amneziawg.conf.example   /etc/mesh/amneziawg.conf.example

#entrypoint.sh fixes its ownership before dropping
#privileges, since bind mounts arrive root-owned.
WORKDIR /home/mesh/acquisitions

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
