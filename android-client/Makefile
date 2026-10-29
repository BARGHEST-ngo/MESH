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

# Reproducibility: default SOURCE_DATE_EPOCH to the last git commit timestamp.
# CI or F-Droid may override this via the environment.
export SOURCE_DATE_EPOCH ?= $(shell git log -1 --format=%ct 2>/dev/null || echo 0)

DEBUG_APK := mesh-debug.apk
RELEASE_AAB := mesh-release.aab

# The MESHnet library, from a MESHnet release or built with MESHnet's scripts/build-aar.sh.
MESHNET_AAR := android/libs/meshnet.aar

# Android SDK & Tools settings.
ifeq ($(shell uname),Linux)
    ANDROID_TOOLS_URL := "https://dl.google.com/android/repository/commandlinetools-linux-14742923_latest.zip"
    ANDROID_TOOLS_SUM := "04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7  commandlinetools-linux-14742923_latest.zip"
else
    ANDROID_TOOLS_URL := "https://dl.google.com/android/repository/commandlinetools-mac-14742923_latest.zip"
    ANDROID_TOOLS_SUM := "04453066b540409d975c676d781da1477479dde3761310f1a7eb92a1dfb15af7  commandlinetools-mac-14742923_latest.zip"
endif
ANDROID_SDK_PACKAGES := 'platforms;android-36' 'extras;android;m2repository' 'platform-tools' 'build-tools;36.0.0'

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

AVD_BASE_IMAGE := "system-images;android-36;google_apis;"
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

export PATH := $(ANDROID_HOME)/cmdline-tools/latest/bin:$(ANDROID_HOME)/platform-tools:$(PATH)

# ------------------------------------------------------------------------------
# Android Build Targets
# ------------------------------------------------------------------------------

.PHONY: apk
apk: $(DEBUG_APK)

.PHONY: mesh-debug
mesh-debug: $(DEBUG_APK)

$(DEBUG_APK): gradle-dependencies
	(cd android && ./gradlew test assembleDebug)
	install -C android/build/outputs/apk/debug/android-debug.apk $@

# Builds the release APK (unsigned, for testing)
RELEASE_APK := mesh-release-unsigned.apk
.PHONY: release-apk
release-apk: $(RELEASE_APK)

$(RELEASE_APK): gradle-dependencies
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

# gradle-dependencies groups together the android sources and the MESHnet AAR needed to assemble tests/debug/release builds.
.PHONY: gradle-dependencies
gradle-dependencies: $(shell find android -type f -not -path "android/build/*" -not -path '*/.*') $(MESHNET_AAR)

$(MESHNET_AAR):
	@echo "Missing $@: download meshnet.aar from a MESHnet release, or build it with MESHnet's scripts/build-aar.sh" >&2
	@exit 1

$(RELEASE_AAB): gradle-dependencies
	@echo "Building release AAB"
	(cd android && ./gradlew test bundleRelease)
	install -C ./android/build/outputs/bundle/release/android-release.aab $@

mesh-test.apk: gradle-dependencies
	(cd android && ./gradlew assembleApplicationTestAndroidTest)
	install -C ./android/build/outputs/apk/androidTest/applicationTest/android-applicationTest-androidTest.apk $@

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
	@$(ANDROID_HOME)/cmdline-tools/latest/bin/sdkmanager --list_installed | grep -q 'platform-tools' || (\
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
	adb shell am start -n com.barghest.mesh/.MainActivity

.PHONY: clean
clean: ## Remove build artifacts. Keeps android/libs/meshnet.aar.
	@echo "Cleaning up old build artifacts"
	-rm -rf android/build $(DEBUG_APK) $(RELEASE_AAB) $(RELEASE_APK) $(RELEASE_APK_SIGNED) mesh-test.apk

.PHONY: help
help: ## Show this help
	@echo "\nSpecify a command. The choices are:\n"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""

.DEFAULT_GOAL := help
