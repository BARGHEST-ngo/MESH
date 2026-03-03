FROM golang:1.26-alpine AS build

RUN apk add --no-cache git patch

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY analyst/src ./analyst
RUN go generate ./...
COPY .git ./.git
RUN ./analyst/build.sh

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
    && "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" "platform-tools;37.0.0"

FROM debian:bookworm-slim

ENV JAVA_HOME=/opt/java/openjdk
ENV ANDROID_SDK_ROOT=/usr/local/android-sdk
ENV PATH="${JAVA_HOME}/bin:${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin:${ANDROID_SDK_ROOT}/platform-tools:${PATH}"

RUN groupadd --gid 10001 mesh \
    && useradd  --uid 10001 --gid mesh --shell /sbin/nologin --create-home mesh \
    && mkdir -p /var/run/tailscale \
    && chown mesh:mesh /var/run/tailscale

COPY --from=android-sdk /opt/java/openjdk              $JAVA_HOME
COPY --from=android-sdk /usr/local/android-sdk         $ANDROID_SDK_ROOT
COPY --from=build       /usr/src/app/analyst/mesh      /usr/bin/mesh
COPY --chmod=755        analyst/meshcli                /usr/bin/meshcli
COPY --chmod=755        analyst/entrypoint.sh          /usr/bin/entrypoint.sh

USER mesh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["/usr/bin/meshcli"]