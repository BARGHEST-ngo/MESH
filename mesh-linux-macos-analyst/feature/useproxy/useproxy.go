// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// Package useproxy registers support for using system proxies.
package useproxy

import (
	"tailscale.com/feature"
	"tailscale.com/net/tshttpproxy"
)

func init() {
	feature.HookProxyFromEnvironment.Set(tshttpproxy.ProxyFromEnvironment)
	feature.HookProxyInvalidateCache.Set(tshttpproxy.InvalidateCache)
	feature.HookProxyGetAuthHeader.Set(tshttpproxy.GetAuthHeader)
	feature.HookProxySetSelfProxy.Set(tshttpproxy.SetSelfProxy)
	feature.HookProxySetTransportGetProxyConnectHeader.Set(tshttpproxy.SetTransportGetProxyConnectHeader)
}
