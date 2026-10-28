// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/BARGHEST-ngo/MESH/analyst/cmd"
	"github.com/BARGHEST-ngo/MESH/analyst/session"
	"github.com/peterbourgon/ff/v3/ffcli"
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

	fmt.Printf("MESH address: %s\n\nRun 'pin' to pair a device\n", addr)

	lines := make(chan string)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()

	for {
		fmt.Print("meshcli> ")
		select {
		case <-ctx.Done():
			fmt.Println()
			return
		case line, ok := <-lines:
			if !ok {
				return // stdin closed (Ctrl-D)
			}
			args := strings.Fields(line)
			switch {
			case len(args) == 0:
				continue
			case args[0] == "quit" || args[0] == "exit":
				return
			}
			if err := commands(s).ParseAndRun(ctx, args); err != nil && err != flag.ErrHelp {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}
	}
}

func commands(s *session.Session) *ffcli.Command {
	root := &ffcli.Command{
		Name:       "meshcli",
		ShortUsage: "<command> [flags]",
		FlagSet:    flag.NewFlagSet("meshcli", flag.ContinueOnError),
		Subcommands: []*ffcli.Command{
			cmd.PinCmd(s),
			cmd.AdbPairCmd(),
			cmd.AdbcollectCmd(),
			cmd.AdbdisableCmd(),
			cmd.AdbcleanCmd(),
			cmd.StatusCmd(),
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
