// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package vms

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/tailscale/hujson"
)

type Distro struct {
	Name           string // amazon-linux
	URL            string // URL to a qcow2 image
	SHA256Sum      string // hex-encoded sha256 sum of contents of URL
	MemoryMegs     int    // VM memory in megabytes
	PackageManager string // yum/apt/dnf/zypper
	InitSystem     string // systemd/openrc
	HostGenerated  bool   // generated image rather than downloaded
}

func (d *Distro) InstallPre() string {
	switch d.PackageManager {
	case "yum":
		return ` - [ yum, update, gnupg2 ]
 - [ yum, "-y", install, iptables ]
 - [ sh, "-c", "printf '\n\nUseDNS no\n\n' | tee -a /etc/ssh/sshd_config" ]
 - [ systemctl, restart, "sshd.service" ]`
	case "zypper":
		return ` - [ zypper, in, "-y", iptables ]`

	case "dnf":
		return ` - [ dnf, install, "-y", iptables ]`

	case "apt":
		return ` - [ apt-get, update ]
 - [ apt-get, "-y", install, curl, "apt-transport-https", gnupg2 ]`

	case "apk":
		return ` - [ apk, "-U", add, curl, "ca-certificates", iptables, ip6tables ]
 - [ modprobe, tun ]`
	}

	return ""
}

//go:embed distros.hujson
var distroData string

var Distros []Distro = func() []Distro {
	var result []Distro
	b, err := hujson.Standardize([]byte(distroData))
	if err != nil {
		log.Fatalf("error decoding distros: %v", err)
	}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Fatalf("error decoding distros: %v", err)
	}
	return result
}()
