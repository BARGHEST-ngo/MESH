// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func UpCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "up",
		ShortUsage: "meshcli up [flags]",
		ShortHelp:  "Connect to the MESH network",
		FlagSet:    flag.NewFlagSet("up", flag.ContinueOnError),
		Exec:       runUpCommand,
	}
}

func runUpCommand(ctx context.Context, args []string) error {
	return execTailscale(append([]string{"up"}, args...))
}
