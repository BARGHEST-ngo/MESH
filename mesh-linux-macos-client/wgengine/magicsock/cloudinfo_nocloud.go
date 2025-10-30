// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ios || android || js

package magicsock

import (
	"context"
	"net/netip"

	"tailscale.com/types/logger"
)

type cloudInfo struct{}

func newCloudInfo(_ logger.Logf) *cloudInfo {
	return &cloudInfo{}
}

func (ci *cloudInfo) GetPublicIPs(_ context.Context) ([]netip.Addr, error) {
	return nil, nil
}
