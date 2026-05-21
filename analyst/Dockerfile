FROM golang:1.26-alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

RUN apk add --no-cache git patch

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY analyst/src ./analyst/src
COPY tailscale ./tailscale
RUN ./analyst/src/build.sh

FROM eclipse-temurin:25-jre@sha256:04262e8782d6b034ee5d7c1c5d4e8938fcf2063a76b4bfcd84e5d994d09c27bc AS android-sdk

ARG ANDROID_TOOLS_VERSION=14742923
ARG ANDROID_TOOLS_SUM=04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7
ENV ANDROID_SDK_ROOT=/usr/local/android-sdk

RUN apt-get update && apt-get install -y --no-install-recommends wget unzip \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p "${ANDROID_SDK_ROOT}/tmp" \
    && wget --progress=dot:giga -O "${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" \
        "https://dl.google.com/android/repository/commandlinetools-linux-${ANDROID_TOOLS_VERSION}_latest.zip" \
    && echo "${ANDROID_TOOLS_SUM}  ${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" | sha256sum --strict --check - \
    && unzip -q "${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" -d "${ANDROID_SDK_ROOT}/tmp" \
    && mkdir -p "${ANDROID_SDK_ROOT}/cmdline-tools" \
    && mv "${ANDROID_SDK_ROOT}/tmp/cmdline-tools" "${ANDROID_SDK_ROOT}/cmdline-tools/latest" \
    && rm -rf "${ANDROID_SDK_ROOT}/tmp" \
    && yes | "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" --licenses || true \
    && "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" "platform-tools"

FROM debian:bookworm-slim@sha256:0104b334637a5f19aa9c983a91b54c89887c0984081f2068983107a6f6c21eeb

ENV ANDROID_SDK_ROOT=/usr/local/android-sdk
ENV PATH="${ANDROID_SDK_ROOT}/platform-tools:${PATH}"

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

COPY --from=android-sdk /usr/local/android-sdk/platform-tools ${ANDROID_SDK_ROOT}/platform-tools
COPY --from=build       /usr/src/app/analyst/mesh             /usr/bin/mesh
COPY --chmod=755        analyst/entrypoint.sh                 /usr/bin/entrypoint.sh
COPY --chown=mesh:mesh  analyst/bashrc                        /home/mesh/.bashrc
COPY                    analyst/amneziawg.conf.example        /etc/mesh/amneziawg.conf.example

# No USER directive. The entrypoint drops to mesh after fixing volume perms.
WORKDIR /home/mesh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
