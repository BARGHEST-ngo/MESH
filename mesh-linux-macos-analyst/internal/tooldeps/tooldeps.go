// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build for_go_mod_tidy_only

// Package tooldeps contains dependencies for tools used in the Tailscale repository,
// so they're not removed by "go mod tidy".
package tooldeps

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/tailscale/depaware/depaware"
	_ "golang.org/x/tools/cmd/goimports"
)
