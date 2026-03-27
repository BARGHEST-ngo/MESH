# Copyright (c) Tailscale Inc & contributors
# SPDX-License-Identifier: BSD-3-Clause
#
# Portions Copyright (c) BARGHEST
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# This file contains code originally from Tailscale (BSD-3-Clause)
# with modifications by BARGHEST. The modified version is licensed
# under AGPL-3.0-or-later. See LICENSE for details.

## For signed release build JKS_PASSWORD must be set to the password for the jks keystore
## and JKS_PATH must be set to the path to the jks keystore.

# These tags disable features in Tailscale that reach out to Tailscale-controlled servers
# ts_omit_clientmetrics excluded: its omit stub is incomplete (missing Type enum)
# and logtail omission already blocks the only metric exfiltration path.
GOTAGS := ts_omit_logtail,ts_omit_netlog,ts_omit_captiveportal,ts_omit_clientupdate,ts_omit_usermetrics,ts_omit_cloud

# Reproducibility: default SOURCE_DATE_EPOCH to the last git commit timestamp.
# CI or F-Droid may override this via the environment.
export SOURCE_DATE_EPOCH ?= $(shell git log -1 --format=%ct 2>/dev/null || echo 0)

# Reproducibility: always pass -trimpath to Go to strip filesystem paths.
export GOFLAGS := $(strip $(GOFLAGS) -trimpath)

# Reproducibility: pin the exact Go toolchain version from go.mod so that
# machines with a newer patch release (e.g. go1.26.1 vs go1.26.0) still
# produce identical binaries.  GOTOOLCHAIN=<ver> forces the Go command to
# download and use exactly that version instead of the locally installed one.
export GOTOOLCHAIN := $(shell grep '^toolchain ' ../go.mod | awk '{print $$2}')

DEBUG_APK := mesh-debug.apk
RELEASE_AAB := mesh-release.aab

# Define output filenames.
LIBTAILSCALE_AAR := android/libs/libtailscale.aar
UNSTRIPPED_AAR := android/libs/libtailscale_unstripped.aar

# All architecture variants of libgojni.so inside the AAR.
AAR_ARCHS := arm64-v8a armeabi-v7a x86 x86_64

# Compute an absolute path for the unstripped AAR.
ABS_UNSTRIPPED_AAR := $(shell pwd)/$(UNSTRIPPED_AAR)

# Android SDK & Tools settings.
ifeq ($(shell uname),Linux)
    ANDROID_TOOLS_URL := "https://dl.google.com/android/repository/commandlinetools-linux-14742923_latest.zip"
    ANDROID_TOOLS_SUM := "04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7  commandlinetools-linux-14742923_latest.zip"
else
    ANDROID_TOOLS_URL := "https://dl.google.com/android/repository/commandlinetools-mac-14742923_latest.zip"
    ANDROID_TOOLS_SUM := "04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7  commandlinetools-mac-14742923_latest.zip"
endif
ANDROID_SDK_PACKAGES := 'platforms;android-34' 'extras;android;m2repository' 'ndk;23.2.8568313' 'platform-tools' 'build-tools;36.0.0'

# Attempt to find an ANDROID_SDK_ROOT / ANDROID_HOME based either from
# preexisting environment or common locations.
export ANDROID_SDK_ROOT ?= $(shell find $$ANDROID_SDK_ROOT $$ANDROID_HOME $$HOME/Library/Android/sdk $$HOME/Android/Sdk $$HOME/AppData/Local/Android/Sdk /usr/lib/android-sdk -maxdepth 1 -type d 2>/dev/null | head -n 1)
ifeq ($(ANDROID_SDK_ROOT),)
    ifeq ($(shell uname),Linux)
        export ANDROID_SDK_ROOT := $(HOME)/Android/Sdk
    else ifeq ($(shell uname),Darwin)
        export ANDROID_SDK_ROOT := $(HOME)/Library/Android/sdk
    else ifneq ($(WINDIR),)
        export ANDROID_SDK_ROOT := $(HOME)/AppData/Local/Android/sdk
    else
        export ANDROID_SDK_ROOT := $(PWD)/android-sdk
endif
endif
export ANDROID_HOME ?= $(ANDROID_SDK_ROOT)

# Use the exact NDK version specified in android/build.gradle.
NDK_VERSION := $(shell grep 'ndkVersion' android/build.gradle | awk -F'"' '{print $$2}')
NDK_ROOT := $(ANDROID_HOME)/ndk/$(NDK_VERSION)
# Only enforce NDK presence for build targets; allow utility/bootstrap
# targets (androidsdk, help, clean, etc.) to run without it.
ifeq ($(filter androidsdk checkandroidsdk help env androidpath clean,$(MAKECMDGOALS)),)
ifneq ($(MAKECMDGOALS),)
ifeq ($(wildcard $(NDK_ROOT)),)
$(error NDK version $(NDK_VERSION) is not installed at $(NDK_ROOT). Run 'make androidsdk' to install it)
endif
endif
endif

# Reproducibility: force gomobile to use the exact NDK version we specify.
# Without this, gomobile scans $ANDROID_HOME/ndk/ and picks the *highest*
# version it finds, which may differ between build environments.
export ANDROID_NDK_HOME := $(NDK_ROOT)

HOST_OS := $(shell uname | tr A-Z a-z)
ifeq ($(HOST_OS),linux)
    STRIP_TOOL := $(NDK_ROOT)/toolchains/llvm/prebuilt/linux-x86_64/bin/llvm-objcopy
else ifeq ($(HOST_OS),darwin)
    STRIP_TOOL := $(NDK_ROOT)/toolchains/llvm/prebuilt/darwin-x86_64/bin/llvm-objcopy
endif


# Attempt to find Android Studio for Linux configuration, which does not have a
# predetermined location.
ANDROID_STUDIO_ROOT ?= $(shell find ~/android-studio /usr/local/android-studio /opt/android-studio /Applications/Android\ Studio.app $(PROGRAMFILES)/Android/Android\ Studio -type d -maxdepth 1 2>/dev/null | head -n 1)

# Set JAVA_HOME to the Android Studio bundled JDK.
export JAVA_HOME ?= $(shell find "$(ANDROID_STUDIO_ROOT)/jbr" "$(ANDROID_STUDIO_ROOT)/jre" "$(ANDROID_STUDIO_ROOT)/Contents/jbr/Contents/Home" "$(ANDROID_STUDIO_ROOT)/Contents/jre/Contents/Home" -maxdepth 1 -type d 2>/dev/null | head -n 1)
ifeq ($(JAVA_HOME),)
    unexport JAVA_HOME
else
    export PATH := $(JAVA_HOME)/bin:$(PATH)
endif

AVD_BASE_IMAGE := "system-images;android-34;google_apis;"
export HOST_ARCH := $(shell uname -m)
ifeq ($(HOST_ARCH),aarch64)
    AVD_IMAGE := "$(AVD_BASE_IMAGE)arm64-v8a"
else ifeq ($(HOST_ARCH),arm64)
    AVD_IMAGE := "$(AVD_BASE_IMAGE)arm64-v8a"
else
    AVD_IMAGE := "$(AVD_BASE_IMAGE)x86_64"
endif
AVD ?= mesh-$(HOST_ARCH)
export AVD_IMAGE
export AVD

GOBIN := $(HOME)/go/bin
export GOBIN

export PATH := $(GOBIN):$(ANDROID_HOME)/cmdline-tools/latest/bin:$(ANDROID_HOME)/platform-tools:$(PATH)
export GOROOT := # Unset

# ------------------------------------------------------------------------------
# Android Build Targets
# ------------------------------------------------------------------------------

.PHONY: apk
apk: $(DEBUG_APK)

.PHONY: mesh-debug
mesh-debug: $(DEBUG_APK)

$(DEBUG_APK): libtailscale debug-symbols gradle-dependencies build-unstripped-aar
	(cd android && ./gradlew test assembleDebug)
	install -C android/build/outputs/apk/debug/android-debug.apk $@

# Builds the release APK (unsigned, for testing)
RELEASE_APK := mesh-release-unsigned.apk
.PHONY: release-apk
release-apk: $(RELEASE_APK)

$(RELEASE_APK): libtailscale gradle-dependencies
	@echo "Building release APK (unsigned)"
	(cd android && ./gradlew assembleRelease)
	install -C android/build/outputs/apk/release/android-release-unsigned.apk $@

# Builds and signs the release APK (for GitHub releases / sideloading)
RELEASE_APK_SIGNED := mesh-release.apk
.PHONY: release-apk-signed
release-apk-signed: jarsign-env $(RELEASE_APK)
	@echo "Signing release APK"
	@cp $(RELEASE_APK) $(RELEASE_APK_SIGNED)
	# --v1-signing-enabled false: v1 is unnecessary (minSdk=26, v2 needs
	#   API 24+) and adding META-INF files would trigger ZIP re-writing.
	# --alignment-preserved: build-tools 35+ apksigner forcibly converts
	#   zero-byte alignment padding to 0xd935 extra-field headers by
	#   default (fdroiddata#3299).  This flag preserves the existing
	#   zero-byte padding produced by AGP, matching what fdroidserver's
	#   build produces so apksigcopier verification passes.
	@apksigner sign --v1-signing-enabled false --alignment-preserved \
		--ks $(JKS_PATH) --ks-pass env:JKS_PASSWORD --ks-key-alias $(JKS_ALIAS) \
		$(RELEASE_APK_SIGNED)
	@echo "Signed APK: $(RELEASE_APK_SIGNED)"

# Builds the release AAB and signs it (phone/tablet/chromeOS variant)
.PHONY: release
release: jarsign-env $(RELEASE_AAB)
	@jarsigner -sigalg SHA256withRSA -digestalg SHA-256 -keystore $(JKS_PATH) -storepass:env JKS_PASSWORD $(RELEASE_AAB) mesh

# gradle-dependencies groups together the android sources and libtailscale needed to assemble tests/debug/release builds.
.PHONY: gradle-dependencies
gradle-dependencies: $(shell find android -type f -not -path "android/build/*" -not -path '*/.*') $(LIBTAILSCALE_AAR)

$(RELEASE_AAB): gradle-dependencies
	@echo "Building release AAB"
	(cd android && ./gradlew test bundleRelease)
	install -C ./android/build/outputs/bundle/release/android-release.aab $@

mesh-test.apk: gradle-dependencies
	(cd android && ./gradlew assembleApplicationTestAndroidTest)
	install -C ./android/build/outputs/apk/androidTest/applicationTest/android-applicationTest-androidTest.apk $@

# ------------------------------------------------------------------------------
# Go Build Targets (Unstripped AAR, Debug Symbols, Stripped SO, Packaging)
# ------------------------------------------------------------------------------

android/libs:
	mkdir -p android/libs

$(GOBIN):
	mkdir -p $(GOBIN)

$(GOBIN)/gomobile: $(GOBIN)/gobind ../go.mod ../go.sum | $(GOBIN)
	go install golang.org/x/mobile/cmd/gomobile

$(GOBIN)/gobind: ../go.mod ../go.sum
	go install golang.org/x/mobile/cmd/gobind

.PHONY: build-unstripped-aar
build-unstripped-aar: $(GOBIN)/gomobile
	@echo "Running gomobile bind to generate unstripped AAR..."
	@echo "Output file: $(ABS_UNSTRIPPED_AAR)"
	mkdir -p $(dir $(ABS_UNSTRIPPED_AAR))
	rm -f $(ABS_UNSTRIPPED_AAR)
	# Reproducibility flags:
	#   CANONICAL_SRC      symlink with a fixed path so gomobile's generated
	#                      replace directives always embed the same module
	#                      path in the Go build info, regardless of the
	#                      actual workspace location on different machines.
	#   GOCACHE=<fresh-dir> ensures gobind's packages.Load always starts from
	#                       a clean cache, avoiding non-deterministic type
	#                       resolution (golang/go#33687).  Without this, the
	#                       build cache state leaks into gobind's output,
	#                       causing invokevirtual ↔ invokeinterface bytecode
	#                       differences in classes.jar across machines.
	#   -buildid=          zeroes the Go build ID
	#   -trimpath          removes host filesystem paths from the binary
	#                      NOTE: passed directly to gomobile bind because
	#                      gomobile strips the parent process environment
	#                      (GOFLAGS, CGO_CFLAGS, GOCACHE, GOTOOLCHAIN, etc.)
	#                      when spawning `go build` subprocesses.  Only
	#                      gomobile's own CLI flags survive to the subprocess.
	#   -buildvcs=false    prevents embedding VCS info (commit hash, dirty flag)
	#   CGO_CFLAGS         strips absolute paths from C debug info
	#   -linkmode=external uses the NDK linker (required for 16kb page alignment)
	#   -s -w              omit symbol table and DWARF at build time so that
	#                      external strip tools (llvm-strip, llvm-objcopy) have
	#                      nothing to remove and cannot re-layout the ELF.
	#                      Go binaries MUST be stripped this way: external tools
	#                      rewrite LOAD segments, shifting code addresses without
	#                      updating Go's .gopclntab PC tables → runtime crash.
	GOBIND_CACHE=$$(mktemp -d) && \
	CANONICAL_SRC="/tmp/_mesh_build_src" && \
	ln -sfn "$$(cd .. && pwd)" "$$CANONICAL_SRC" && \
	trap "rm -rf $$GOBIND_CACHE $$CANONICAL_SRC" EXIT && \
	cd "$$CANONICAL_SRC/android-client" && \
	CGO_CFLAGS="-ffile-prefix-map=$$(pwd)=. -ffile-prefix-map=$${HOME}=." \
	CGO_CXXFLAGS="-ffile-prefix-map=$$(pwd)=. -ffile-prefix-map=$${HOME}=." \
	GOCACHE="$$GOBIND_CACHE" \
	GONOSUMCHECK=* \
	GONOSUMDB=* \
	GOFLAGS="$(GOFLAGS) -buildvcs=false" \
	gomobile bind -trimpath -target android -androidapi 26 \
		-tags $(GOTAGS) \
		-ldflags "-s -w -buildid= -linkmode=external -extldflags=-Wl,-z,max-page-size=16384 $$(./version-ldflags.sh)" \
		-o $(ABS_UNSTRIPPED_AAR) ./libtailscale || { echo "gomobile bind failed"; exit 1; }
	@if [ ! -f $(ABS_UNSTRIPPED_AAR) ]; then \
	    echo "Error: $(ABS_UNSTRIPPED_AAR) was not created"; exit 1; \
	fi
	@echo "Generated unstripped AAR: $(ABS_UNSTRIPPED_AAR)"

$(UNSTRIPPED_AAR): build-unstripped-aar

# Extract a single arch for standalone debug-symbol extraction (best-effort).
libgojni.so.unstripped: $(UNSTRIPPED_AAR)
	@echo "Extracting libgojni.so from unstripped AAR..."
	@if unzip -p $(ABS_UNSTRIPPED_AAR) jni/arm64-v8a/libgojni.so > libgojni.so.unstripped; then \
	    echo "Found arm64-v8a libgojni.so"; \
	elif unzip -p $(ABS_UNSTRIPPED_AAR) jni/armeabi-v7a/libgojni.so > libgojni.so.unstripped; then \
	    echo "Found armeabi-v7a libgojni.so"; \
	else \
	    echo "Neither jni/arm64-v8a/libgojni.so nor jni/armeabi-v7a/libgojni.so was found."; \
	    echo "Listing contents of $(ABS_UNSTRIPPED_AAR):"; \
	    unzip -l $(ABS_UNSTRIPPED_AAR); exit 1; \
	fi

libgojni.so.debug: libgojni.so.unstripped
	@echo "Extracting debug symbols from libgojni.so..."
	$(STRIP_TOOL) --only-keep-debug libgojni.so.unstripped libgojni.so.debug

libgojni.so.stripped: libgojni.so.unstripped
	@echo "Stripping debug symbols from libgojni.so..."
	$(STRIP_TOOL) --strip-debug libgojni.so.unstripped libgojni.so.stripped

# Repackage the AAR with ALL architecture variants stripped.
# This is critical for reproducibility: unstripped .so files contain
# absolute build paths in debug info, which differ between environments.
$(LIBTAILSCALE_AAR): $(UNSTRIPPED_AAR)
	@echo "Repackaging AAR with stripped libgojni.so for all architectures..."
	rm -rf temp_aar
	mkdir temp_aar
	unzip $(ABS_UNSTRIPPED_AAR) -d temp_aar
	# Strip and normalize every architecture's libgojni.so for reproducibility.
	# Per arch:
	#   1. llvm-objcopy --strip-debug removes DWARF and debug sections.
	#      IMPORTANT: do NOT use --strip-all or --remove-section on Go
	#      shared libraries.  Both re-lay out LOAD segments, shifting code
	#      addresses without updating Go's .gopclntab PC tables, which
	#      causes a fatal stack overflow at runtime.  --strip-debug only
	#      removes non-ALLOC debug sections and preserves the loaded
	#      binary layout.
	#   2. normalize-elf zeros remaining non-deterministic content in-place
	#      (preserving binary layout): Go build info strings in .rodata,
	#      .eh_frame/.eh_frame_hdr, .relro_padding, and the ELF entry
	#      point.  This covers the Go metadata sections (.go.buildinfo,
	#      .go.module, etc.) that cannot be --remove-section'd safely.
	@for arch in $(AAR_ARCHS); do \
		so_path="temp_aar/jni/$${arch}/libgojni.so"; \
		if [ -f "$${so_path}" ]; then \
			echo "Stripping $${so_path}"; \
			$(STRIP_TOOL) --strip-debug "$${so_path}"; \
			go run ./cmd/normalize-elf "$${so_path}"; \
		fi; \
	done
	# Normalize timestamps and permissions inside classes.jar (nested zip)
	# for reproducibility.  The umask differs between build environments,
	# causing different Unix permission bits in the ZIP external attributes.
	@if [ -f temp_aar/classes.jar ]; then \
		mkdir -p temp_aar/classes_tmp && \
		(cd temp_aar/classes_tmp && unzip -q ../classes.jar) && \
		find temp_aar/classes_tmp -type d -exec chmod 0755 {} + && \
		find temp_aar/classes_tmp -type f -exec chmod 0644 {} + && \
		find temp_aar/classes_tmp -exec touch -d @$${SOURCE_DATE_EPOCH:-0} {} + && \
		rm temp_aar/classes.jar && \
		(cd temp_aar/classes_tmp && find . -type f | LC_ALL=C sort | zip -X -@ ../classes.jar) && \
		rm -rf temp_aar/classes_tmp; \
	fi
	# Normalize all file timestamps and permissions, then produce a
	# deterministic zip.
	find temp_aar -type d -exec chmod 0755 {} +
	find temp_aar -type f -exec chmod 0644 {} +
	find temp_aar -exec touch -d @$${SOURCE_DATE_EPOCH:-0} {} +
	(cd temp_aar && find . -type f | LC_ALL=C sort | zip -X -@ ../$(LIBTAILSCALE_AAR))
	rm -rf temp_aar

.PHONY: libtailscale
libtailscale: $(LIBTAILSCALE_AAR) ## Build the stripped libtailscale AAR

.PHONY: debug-symbols
debug-symbols: libgojni.so.debug

# ------------------------------------------------------------------------------
# Utility Targets
# ------------------------------------------------------------------------------

.PHONY: env
env:
	@echo "PATH=$(PATH)"
	@echo "ANDROID_SDK_ROOT=$(ANDROID_SDK_ROOT)"
	@echo "ANDROID_HOME=$(ANDROID_HOME)"
	@echo "ANDROID_STUDIO_ROOT=$(ANDROID_STUDIO_ROOT)"
	@echo "JAVA_HOME=$(JAVA_HOME)"
	@echo "AVD_IMAGE=$(AVD_IMAGE)"

# Ensure that JKS_PATH, JKS_PASSWORD, and JKS_ALIAS are set before we attempt a build
# that requires signing.
.PHONY: jarsign-env
jarsign-env:
ifeq ($(JKS_PATH),)
	$(error JKS_PATH is not set.  export JKS_PATH=/path/to/mesh.jks)
endif
ifeq ($(JKS_PASSWORD),)
	$(error JKS_PASSWORD is not set.  export JKS_PASSWORD=yourPassword)
endif
ifeq ($(JKS_ALIAS),)
	$(error JKS_ALIAS is not set.  export JKS_ALIAS=mesh)
endif
ifeq ($(wildcard $(JKS_PATH)),)
	$(error JKS_PATH does not point to a file)
endif
	@echo "keystore path set to $(JKS_PATH)"

.PHONY: androidpath
androidpath:
	@echo "export ANDROID_HOME=$(ANDROID_HOME)"
	@echo "export ANDROID_SDK_ROOT=$(ANDROID_SDK_ROOT)"
	@echo 'export PATH=$(ANDROID_HOME)/cmdline-tools/latest/bin:$(ANDROID_HOME)/platform-tools:$$PATH'

.PHONY: bump_version_code ## Bump the version code in build.gradle
bump_version_code:
	sed -i'.bak' "s/versionCode .*/versionCode $$(expr $$(awk '/versionCode ([0-9]+)/{print $$2}' android/build.gradle) + 1)/" android/build.gradle && rm android/build.gradle.bak

# Get the commandline tools package, this provides (among other things) the sdkmanager binary.
$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager:
	mkdir -p $(ANDROID_HOME)/tmp
	mkdir -p $(ANDROID_HOME)/cmdline-tools
	(cd $(ANDROID_HOME)/tmp && \
		curl --silent -O -L $(ANDROID_TOOLS_URL) && \
		echo $(ANDROID_TOOLS_SUM) | shasum -c - && \
		unzip $(shell basename $(ANDROID_TOOLS_URL)))
	mv $(ANDROID_HOME)/tmp/cmdline-tools $(ANDROID_HOME)/cmdline-tools/latest
	rm -rf $(ANDROID_HOME)/tmp

.PHONY: androidsdk
androidsdk: $(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager ## Install the set of Android SDK packages we need.
	yes | $(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager --licenses > /dev/null
	$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager --update
	$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager $(ANDROID_SDK_PACKAGES)

# Normally in make you would simply take a dependency on the task that provides
# the binaries, however users may have a decision to make as to whether they
# want to install an SDK or use the one from an Android Studio installation.
.PHONY: checkandroidsdk
checkandroidsdk: ## Check that Android SDK is installed
	@$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager --list_installed | grep -q 'ndk' || (\
		echo -e "\n\tERROR: Android SDK not installed.\n\
		\tANDROID_HOME=$(ANDROID_HOME)\n\
		\tANDROID_SDK_ROOT=$(ANDROID_SDK_ROOT)\n\n\
		See README.md for instructions on how to install the prerequisites.\n"; exit 1)

.PHONY: test
test: gradle-dependencies ## Run the Android tests
	(cd android && ./gradlew test)

.PHONY: fmt
fmt: gradle-dependencies ## Format the Android code
	(cd android && ./gradlew ktfmtFormat)

.PHONY: fmt-check
fmt-check: gradle-dependencies ## Check the Android code is formatted
	(cd android && ./gradlew ktfmtCheck)

.PHONY: emulator
emulator: ## Start an android emulator instance
	@echo "Checking installed SDK packages..."
	@if ! $(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager --list_installed | grep -q "$(AVD_IMAGE)"; then \
		echo "$(AVD_IMAGE) not found, installing..."; \
		$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager "$(AVD_IMAGE)"; \
	fi
	@echo "Checking if AVD exists..."
	@if ! $(ANDROID_HOME)/cmdline-tools/latest/bin/avdmanager list avd | grep -q "$(AVD)"; then \
		echo "AVD $(AVD) not found, creating..."; \
		$(ANDROID_HOME)/cmdline-tools/latest/bin/avdmanager create avd -n "$(AVD)" -k "$(AVD_IMAGE)"; \
	fi
	@echo "Starting emulator..."
	@$(ANDROID_HOME)/emulator/emulator -avd "$(AVD)" -logcat-output /dev/stdout -netdelay none -netspeed full

.PHONY: install
install: $(DEBUG_APK) ## Install the debug APK on a connected device
	adb install -r $<

.PHONY: run
run: install ## Run the debug APK on a connected device
	adb shell am start -n com.tailscale.ipn/com.tailscale.ipn.MainActivity

.PHONY: clean
clean: ## Remove build artifacts. Does not purge docker build envs. Use dockerRemoveEnv for that.
	@echo "Cleaning up old build artifacts"
	-rm -rf android/build $(DEBUG_APK) $(RELEASE_AAB) $(LIBTAILSCALE_AAR) android/libs *.apk *.aab
	@echo "Cleaning cached toolchain"
	-rm -rf $(HOME)/.cache/tailscale-go{,.extracted}

.PHONY: help
help: ## Show this help
	@echo "\nSpecify a command. The choices are:\n"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""

.DEFAULT_GOAL := help
