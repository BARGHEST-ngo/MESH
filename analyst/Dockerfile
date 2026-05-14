FROM golang:1.26-alpine AS build

RUN apk add --no-cache git patch

WORKDIR /usr/src/app

# Disable workspace so each module resolves independently
ENV GOWORK=off

# Cache root module deps (used by tailscale submodule build)
COPY go.mod go.sum ./
RUN go mod download

# Cache analyst module deps separately
COPY analyst/go.mod analyst/go.sum ./analyst/
RUN cd analyst && go mod download

COPY analyst ./analyst
COPY tailscale ./tailscale
RUN ./analyst/src/build.sh

FROM eclipse-temurin:25-jre AS android-sdk

ARG ANDROID_TOOLS_VERSION=14742923
ARG ANDROID_TOOLS_SUM=04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7
ENV ANDROID_SDK_ROOT=/usr/local/android-sdk

RUN apt-get update && apt-get install -y --no-install-recommends \
        wget \
		unzip \
    && rm -rf /var/lib/apt/lists/* \
    && ANDROID_TOOLS_URL="https://dl.google.com/android/repository/commandlinetools-linux-${ANDROID_TOOLS_VERSION}_latest.zip" \
    && mkdir -p "${ANDROID_SDK_ROOT}/tmp" \
    && wget --progress=dot:giga -O "${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" "${ANDROID_TOOLS_URL}" \
    && echo "${ANDROID_TOOLS_SUM}  ${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" | sha256sum --strict --check - \
    && unzip -q "${ANDROID_SDK_ROOT}/tmp/commandlinetools.zip" -d "${ANDROID_SDK_ROOT}/tmp" \
    && mkdir -p "${ANDROID_SDK_ROOT}/cmdline-tools" \
    && mv "${ANDROID_SDK_ROOT}/tmp/cmdline-tools" "${ANDROID_SDK_ROOT}/cmdline-tools/latest" \
    && rm -rf "${ANDROID_SDK_ROOT}/tmp" \
    && yes | "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" --licenses || true \
    && "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" "platform-tools"

FROM debian:bookworm-slim

ENV JAVA_HOME=/opt/java/openjdk
ENV ANDROID_SDK_ROOT=/usr/local/android-sdk
ENV PATH="${JAVA_HOME}/bin:${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin:${ANDROID_SDK_ROOT}/platform-tools:${PATH}"

RUN apt-get update && apt-get install -y --no-install-recommends \
        iptables \
        iproute2 \
        iputils-ping \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /root/.tailscale

COPY --from=android-sdk /opt/java/openjdk              $JAVA_HOME
COPY --from=android-sdk /usr/local/android-sdk         $ANDROID_SDK_ROOT
COPY --from=build       /usr/src/app/analyst/tailscaled  /usr/bin/tailscaled
COPY --from=build       /usr/src/app/analyst/tailscale   /usr/bin/tailscale
COPY --from=build       /usr/src/app/analyst/mesh-analyst /usr/bin/mesh-analyst
COPY --chmod=755        analyst/entrypoint.sh          /usr/bin/entrypoint.sh
COPY                    analyst/amneziawg.conf.example /etc/mesh/amneziawg.conf.example

USER root

# Configure .bashrc: mesh-analyst alias and interactive welcome message
RUN printf '%s\n' \
    '' \
    'alias mesh="/usr/bin/mesh-analyst"' \
    '' \
    '# Display welcome message for interactive shells only' \
    'if [[ $- == *i* ]]; then' \
    '  printf "\n\033[1;36m╔══════════════════════════════════════════════╗\033[0m\n"' \
    '  printf "\033[1;36m║\033[0m\033[1;37m          MESH Analyst Container              \033[0m\033[1;36m║\033[0m\n"' \
    '  printf "\033[1;36m╚══════════════════════════════════════════════╝\033[0m\n\n"' \
    '  printf "  \033[1;33mForensic Tools:\033[0m\n\n"' \
    '  printf "    \033[1;32mmesh\033[0m           MESH CLI — manage MESH network connections\n"' \
    '  printf "               Usage: mesh <subcommand> [command flags]\n\n"' \
    '  printf "               SUBCOMMANDS\n"' \
    '  printf "                 status       Show state of tailscaled and its connections\n"' \
    '  printf "                 adbpair      Pair & connect to a device on the MESH network via ADB\n"' \
    '  printf "                 adbcollect   Collects ADB data, using AndroidQF and WARD\n"' \
    '  printf "                 adbdisable   The ADB disable utility disables the developer mode on the Android node\n\n"' \
    '  printf "               For help on subcommands, add --help after: \"mesh status --help\".\n\n"' \
    '  printf "  \033[1;33mLogs:\033[0m\n"' \
    '  printf "    ~/.tailscale/mesh.log     MESH daemon log\n"' \
    '  printf "    ~/.tailscale/meshcli.log  MESH CLI log\n\n"' \
    'fi' \
    >> /root/.bashrc

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
