// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
