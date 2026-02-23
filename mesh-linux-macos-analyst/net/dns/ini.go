// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build windows

package dns

import (
	"regexp"
	"strings"
)

// parseIni parses a basic .ini file, used for wsl.conf.
func parseIni(data string) map[string]map[string]string {
	sectionRE := regexp.MustCompile(`^\[([^]]+)\]`)
	kvRE := regexp.MustCompile(`^\s*(\w+)\s*=\s*([^#]*)`)

	ini := map[string]map[string]string{}
	var section string
	for _, line := range strings.Split(data, "\n") {
		if res := sectionRE.FindStringSubmatch(line); len(res) > 1 {
			section = res[1]
			ini[section] = map[string]string{}
		} else if res := kvRE.FindStringSubmatch(line); len(res) > 2 {
			k, v := strings.TrimSpace(res[1]), strings.TrimSpace(res[2])
			ini[section][k] = v
		}
	}
	return ini
}
