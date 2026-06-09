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

FROM debian:bookworm-slim@sha256:0104b334637a5f19aa9c983a91b54c89887c0984081f2068983107a6f6c21eeb

RUN apt-get update && apt-get install -y --no-install-recommends \
        iptables \
        iproute2 \
        iputils-ping \
        libcap2-bin \
        tcpdump \
        util-linux \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1000 mesh \
    && useradd  --system --uid 1000 --gid 1000 \
                --create-home --home-dir /home/mesh \
                --shell /bin/bash mesh \
    && mkdir -p /home/mesh/.tailscale \
    && chown -R mesh:mesh /home/mesh \
    && setcap cap_net_admin,cap_net_raw+eip /usr/bin/tcpdump

RUN apt-get install -y --no-install-recommends adb \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build       /usr/src/app/analyst/tailscaled  /usr/bin/tailscaled
COPY --from=build       /usr/src/app/analyst/tailscale   /usr/bin/tailscale
COPY --from=build       /usr/src/app/analyst/mesh        /usr/bin/mesh
COPY --chmod=755        analyst/entrypoint.sh            /usr/bin/entrypoint.sh
COPY --chown=mesh:mesh  analyst/bashrc                   /home/mesh/.bashrc
COPY                    analyst/amneziawg.conf.example   /etc/mesh/amneziawg.conf.example

# No USER directive. The entrypoint drops to mesh after fixing volume perms.
WORKDIR /home/mesh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
