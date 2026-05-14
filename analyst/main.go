// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/BARGHEST-ngo/MESH/analyst/cmd"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	root := &ffcli.Command{
		Name:       "mesh-analyst",
		ShortUsage: "mesh-analyst <command> [flags]",
		ShortHelp:  "MESH analyst CLI",
		Subcommands: []*ffcli.Command{
			cmd.HelloCmd(),
			cmd.StatusCmd(),
		},
		FlagSet: flag.NewFlagSet("mesh", flag.ContinueOnError),
		Exec: func(ctx context.Context, args []string) error {
			fmt.Fprintf(os.Stderr, "Run 'mesh-analyst help' for usage.\n")
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}
