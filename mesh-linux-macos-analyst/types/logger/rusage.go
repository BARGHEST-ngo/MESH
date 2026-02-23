// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package logger

import (
	"fmt"
	"runtime"
)

// RusagePrefixLog returns a Logf func wrapping the provided logf func that adds
// a prefixed log message to each line with the current binary memory usage
// and max RSS.
func RusagePrefixLog(logf Logf) Logf {
	return func(f string, argv ...any) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		goMem := float64(m.HeapInuse+m.StackInuse) / (1 << 20)
		maxRSS := rusageMaxRSS()
		pf := fmt.Sprintf("%.1fM/%.1fM %s", goMem, maxRSS, f)
		logf(pf, argv...)
	}
}
