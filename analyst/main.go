// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/BARGHEST-ngo/MESH/analyst/cmd"
	"github.com/BARGHEST-ngo/MESH/analyst/session"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/skip2/go-qrcode"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	s := session.New(os.Stdout)
	addr, err := s.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "start: %v\n", err)
		os.Exit(1)
	}
	defer s.Stop()

	qr, err := qrcode.New(addr, qrcode.Low)
	if err != nil {
		fmt.Fprintf(os.Stderr, "qr code: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(qr.ToSmallString(false))
	fmt.Printf("MESH address: %s\n\nRun 'pin' to pair a device\n", addr)

	in := cmd.NewInput(os.Stdin)
	for {
		line, err := in.ReadLine(ctx, "meshcli> ")
		if err != nil {
			return // Ctrl-C or end of input
		}
		args := strings.Fields(line)
		switch {
		case len(args) == 0:
			continue
		case args[0] == "quit" || args[0] == "exit":
			return
		}
		if err := commands(s, in).ParseAndRun(ctx, args); err != nil && err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func commands(s *session.Session, in *cmd.Input) *ffcli.Command {
	root := &ffcli.Command{
		Name:       "meshcli",
		ShortUsage: "<command> [flags]",
		FlagSet:    flag.NewFlagSet("meshcli", flag.ContinueOnError),
		Subcommands: []*ffcli.Command{
			cmd.PinCmd(s),
			cmd.AdbPairCmd(s, in),
			cmd.AdbcollectCmd(),
			cmd.AdbdisableCmd(),
			cmd.AdbcleanCmd(),
			cmd.StatusCmd(),
			cmd.DisconnectCmd(s, in),
		},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Fprintf(os.Stderr, "Run 'meshcli help' for usage.\n")
			return flag.ErrHelp
		},
	}

	root.Exec = func(ctx context.Context, args []string) error {
		fmt.Println(ffcli.DefaultUsageFunc(root))
		return nil
	}
	return root
}
