// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package portlist

import (
	"path/filepath"
	"strings"
)

// argvSubject takes a command and its flags, and returns the
// short/pretty name for the process. This is usually the basename of
// the binary being executed, but can sometimes vary (e.g. so that we
// don't report all Java programs as "java").
func argvSubject(argv ...string) string {
	if len(argv) == 0 {
		return ""
	}
	ret := filepath.Base(argv[0])

	// Handle special cases.
	switch {
	case ret == "mono" && len(argv) >= 2:
		// .Net programs execute as `mono actualProgram.exe`.
		ret = filepath.Base(argv[1])
	}

	// Handle space separated argv
	ret, _, _ = strings.Cut(ret, " ")

	// Remove common noise.
	ret = strings.TrimSpace(ret)
	ret = strings.TrimSuffix(ret, ".exe")

	return ret
}
