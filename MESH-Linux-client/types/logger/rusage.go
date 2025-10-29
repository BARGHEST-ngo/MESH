// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
