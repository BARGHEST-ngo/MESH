// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !ios && !js && !android && !ts_omit_useproxy

package netns

import "golang.org/x/net/proxy"

func init() {
	wrapDialer = wrapSocks
}

func wrapSocks(d Dialer) Dialer {
	if cd, ok := proxy.FromEnvironmentUsing(d).(Dialer); ok {
		return cd
	}
	return d
}
