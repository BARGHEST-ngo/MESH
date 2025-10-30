// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build linux && !android && !ts_package_container

package hostinfo

import (
	"testing"
)

func TestQnap(t *testing.T) {
	version_info := `commit 2910d3a594b068024ed01a64a0fe4168cb001a12
Date:   2022-05-30 16:08:45 +0800
================================================
* QTSFW_5.0.0
remotes/origin/QTSFW_5.0.0`

	got := getQnapQtsVersion(version_info)
	want := "5.0.0"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}

	got = getQnapQtsVersion("")
	want = ""
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}

	got = getQnapQtsVersion("just a bunch of junk")
	want = ""
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestPackageTypeNotContainer(t *testing.T) {
	var got string
	if packageType != nil {
		got = packageType()
	}
	if got == "container" {
		t.Fatal("packageType = container; should only happen if build tag ts_package_container is set")
	}
}
