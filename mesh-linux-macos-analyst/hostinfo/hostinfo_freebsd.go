// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build freebsd

package hostinfo

import (
	"bytes"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
	"tailscale.com/types/ptr"
	"tailscale.com/version/distro"
)

func init() {
	osVersion = lazyOSVersion.Get
	distroName = distroNameFreeBSD
	distroVersion = distroVersionFreeBSD
}

var (
	lazyVersionMeta = &lazyAtomicValue[versionMeta]{f: ptr.To(freebsdVersionMeta)}
	lazyOSVersion   = &lazyAtomicValue[string]{f: ptr.To(osVersionFreeBSD)}
)

func distroNameFreeBSD() string {
	return lazyVersionMeta.Get().DistroName
}

func distroVersionFreeBSD() string {
	return lazyVersionMeta.Get().DistroVersion
}

type versionMeta struct {
	DistroName     string
	DistroVersion  string
	DistroCodeName string
}

func osVersionFreeBSD() string {
	var un unix.Utsname
	unix.Uname(&un)
	return unix.ByteSliceToString(un.Release[:])
}

func freebsdVersionMeta() (meta versionMeta) {
	d := distro.Get()
	meta.DistroName = string(d)
	switch d {
	case distro.Pfsense:
		b, _ := os.ReadFile("/etc/version")
		meta.DistroVersion = string(bytes.TrimSpace(b))
	case distro.OPNsense:
		b, _ := exec.Command("opnsense-version").Output()
		meta.DistroVersion = string(bytes.TrimSpace(b))
	case distro.TrueNAS:
		b, _ := os.ReadFile("/etc/version")
		meta.DistroVersion = string(bytes.TrimSpace(b))
	}
	return
}
