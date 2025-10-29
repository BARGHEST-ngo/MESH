// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package version

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var stringLazy = sync.OnceValue(func() string {
	var ret strings.Builder
	ret.WriteString(Short())
	ret.WriteByte('\n')
	if IsUnstableBuild() {
		fmt.Fprintf(&ret, "  track: unstable (dev); frequent updates and bugs are likely\n")
	}
	if gitCommit() != "" {
		fmt.Fprintf(&ret, "  tailscale commit: %s%s\n", gitCommit(), dirtyString())
	}
	fmt.Fprintf(&ret, "  long version: %s\n", Long())
	if extraGitCommitStamp != "" {
		fmt.Fprintf(&ret, "  other commit: %s\n", extraGitCommitStamp)
	}
	fmt.Fprintf(&ret, "  go version: %s\n", runtime.Version())
	return strings.TrimSpace(ret.String())
})

func String() string {
	return stringLazy()
}
