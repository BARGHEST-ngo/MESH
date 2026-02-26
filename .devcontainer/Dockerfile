ARG TARGET_JAVA_VERSION=17
ARG GO_VERSION=1.26.0
ARG NODE_VERSION=22
ARG BASE_IMAGE_VERSION_CODENAME=trixie
FROM mcr.microsoft.com/devcontainers/base:${BASE_IMAGE_VERSION_CODENAME}

USER root
ARG TARGET_JAVA_VERSION
ARG GO_VERSION
ARG NODE_VERSION

# Set up environment variables
ENV JAVA_HOME=/usr/lib/jvm/msopenjdk-current
ENV ANDROID_SDK_ROOT=/usr/local/android-sdk
ENV ANDROID_HOME=/usr/local/android-sdk
ENV PATH="${JAVA_HOME}/bin:/usr/local/go/bin:${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin:${ANDROID_SDK_ROOT}/platform-tools:${ANDROID_SDK_ROOT}/build-tools/34.0.0:${PATH}"
# Default to UTF-8 file.encoding
ENV LANG=en_US.UTF-8

# Install Microsoft OpenJDK
RUN arch="$(dpkg --print-architecture)" \
	&& case "$arch" in \
		"amd64") \
			jdkUrl="https://aka.ms/download-jdk/microsoft-jdk-${TARGET_JAVA_VERSION}-linux-x64.tar.gz"; \
			;; \
		"arm64") \
			jdkUrl="https://aka.ms/download-jdk/microsoft-jdk-${TARGET_JAVA_VERSION}-linux-aarch64.tar.gz"; \
			;; \
		*) echo >&2 "error: unsupported architecture: '$arch'"; exit 1 ;; \
	esac \
	\
	&& wget --progress=dot:giga -O msopenjdk.tar.gz "${jdkUrl}" \
	&& wget --progress=dot:giga -O sha256sum.txt "${jdkUrl}.sha256sum.txt" \
	\
	&& sha256sumText=$(cat sha256sum.txt) \
	&& sha256=$(expr substr "${sha256sumText}" 1 64) \
	&& echo "${sha256} msopenjdk.tar.gz" | sha256sum --strict --check - \
	&& rm sha256sum.txt* \
	\
	&& mkdir -p "$JAVA_HOME" \
	&& tar --extract \
		--file msopenjdk.tar.gz \
		--directory "$JAVA_HOME" \
		--strip-components 1 \
		--no-same-owner \
	&& rm msopenjdk.tar.gz* \
	\
	&& ln -s ${JAVA_HOME} /docker-java-home \
	&& ln -s ${JAVA_HOME} /usr/local/openjdk-${TARGET_JAVA_VERSION}

# Install Go
RUN GO_DOWNLOAD_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
	&& echo "Installing Go ${GO_VERSION}..." \
	&& wget --progress=dot:giga -O go.tar.gz "${GO_DOWNLOAD_URL}" \
	&& tar -C /usr/local -xzf go.tar.gz \
	&& rm go.tar.gz \
	&& /usr/local/go/bin/go version

# Install Node.js and pnpm
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
	&& apt-get install -y nodejs \
	&& npm install -g pnpm \
	&& node --version \
	&& npm --version \
	&& pnpm --version

# Install Android SDK
RUN ANDROID_TOOLS_URL="https://dl.google.com/android/repository/commandlinetools-linux-14742923_latest.zip" \
	&& ANDROID_TOOLS_SUM="04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7" \
	&& echo "Installing Android SDK..." \
	&& mkdir -p "${ANDROID_SDK_ROOT}/tmp" \
	&& cd "${ANDROID_SDK_ROOT}/tmp" \
	&& wget --progress=dot:giga -O commandlinetools.zip "${ANDROID_TOOLS_URL}" \
	&& echo "${ANDROID_TOOLS_SUM}  commandlinetools.zip" | sha256sum --strict --check - \
	&& unzip -q commandlinetools.zip \
	&& mkdir -p "${ANDROID_SDK_ROOT}/cmdline-tools" \
	&& mv cmdline-tools "${ANDROID_SDK_ROOT}/cmdline-tools/latest" \
	&& cd / \
	&& rm -rf "${ANDROID_SDK_ROOT}/tmp" \
	&& yes | "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" --licenses || true \
	&& "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" --update \
	&& "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" \
		"platforms;android-34" \
		"extras;android;m2repository" \
		"ndk;23.1.7779620" \
		"platform-tools" \
		"build-tools;34.0.0" \
	&& "${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin/sdkmanager" --list_installed

# Install Task
RUN curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | bash \
	&& apt-get install -y task