// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func HelloCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "hello",
		ShortUsage: "mesh hello",
		ShortHelp:  "hello world!",
		FlagSet:    flag.NewFlagSet("hello", flag.ContinueOnError),
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println("hello world, mesh analyst")
			return nil
		},
	}
}
