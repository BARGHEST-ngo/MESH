// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !tailscale_go || !(darwin || ios || android || ts_enable_sockstats)

package sockstats

import (
	"context"

	"tailscale.com/net/netmon"
	"tailscale.com/types/logger"
)

const IsAvailable = false

func withSockStats(ctx context.Context, label Label, logf logger.Logf) context.Context {
	return ctx
}

func get() *SockStats {
	return nil
}

func getInterfaces() *InterfaceSockStats {
	return nil
}

func getValidation() *ValidationSockStats {
	return nil
}

func setNetMon(netMon *netmon.Monitor) {
}

func debugInfo() string {
	return ""
}
