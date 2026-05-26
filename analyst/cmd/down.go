// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func DownCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "down",
		ShortUsage: "meshcli down [flags]",
		ShortHelp:  "Disconnect from the MESH network",
		FlagSet:    flag.NewFlagSet("down", flag.ContinueOnError),
		Exec:       runDownCommand,
	}
}

func runDownCommand(ctx context.Context, args []string) error {
	return execTailscale(append([]string{"down"}, args...))
}
