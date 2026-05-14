// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func StatusCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "mesh status [flags]",
		ShortHelp:  "Show MESH network status",
		FlagSet:    flag.NewFlagSet("status", flag.ContinueOnError),
		Exec: func(ctx context.Context, args []string) error {
			return execTailscale(append([]string{"status"}, args...))
		},
	}
}
