// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BARGHEST-ngo/MESH/analyst/cmd"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var nativeCommands = map[string]bool{
	"adbpair":    true,
	"adbcollect": true,
	"adbdisable": true,
	"adbclean":   true,
	"status":     true,
	"help":       true,
}

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if command != "" && !strings.HasPrefix(command, "-") && !nativeCommands[command] {
		if err := cmd.ExecTailscale(os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	root := &ffcli.Command{
		Name:       "meshcli",
		ShortUsage: "meshcli <command> [flags]",
		ShortHelp:  "MESH analyst CLI",
		LongHelp:   "Native MESH commands are listed below. Standard tailscale commands are also supported.",
		Subcommands: []*ffcli.Command{
			cmd.AdbpairCmd(),
			cmd.AdbcollectCmd(),
			cmd.AdbdisableCmd(),
			cmd.AdbcleanCmd(),
			cmd.StatusCmd(),
		},
		FlagSet: flag.NewFlagSet("meshcli", flag.ContinueOnError),
		Exec: func(ctx context.Context, args []string) error {
			fmt.Fprintf(os.Stderr, "Run 'meshcli help' for usage.\n")
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
