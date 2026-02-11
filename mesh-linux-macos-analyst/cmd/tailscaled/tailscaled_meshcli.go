// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ts_include_cli

package main

import (
	"fmt"
	"os"

	"tailscale.com/cmd/meshcli/cli"
)

func init() {
	beCLI = func() {
		if err := cli.Run(os.Args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

