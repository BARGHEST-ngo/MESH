// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build linux && !ts_omit_systray

package cli

import (
	"context"

	"github.com/peterbourgon/ff/v3/ffcli"
	"tailscale.com/client/systray"
)

func init() {
	maybeSystrayCmd = systrayCmd
}

func systrayCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "systray",
		ShortUsage: "tailscale systray",
		ShortHelp:  "Run a systray application to manage Tailscale",
		LongHelp:   "Run a systray application to manage Tailscale.",
		Exec:       runSystray,
	}
}

func runSystray(ctx context.Context, _ []string) error {
	new(systray.Menu).Run(&localClient)
	return nil
}
