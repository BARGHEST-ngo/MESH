// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func SetCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "set",
		ShortUsage: "meshcli set [flags]",
		ShortHelp:  "Connect to the MESH network",
		FlagSet:    flag.NewFlagSet("set", flag.ContinueOnError),
		Exec:       runSetCommand,
	}
}

func runSetCommand(ctx context.Context, args []string) error {
	return execTailscale(append([]string{"set"}, args...))
}
