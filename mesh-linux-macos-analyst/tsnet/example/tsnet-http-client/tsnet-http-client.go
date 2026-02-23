// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

// The tshello server demonstrates how to use Tailscale as a library.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"tailscale.com/tsnet"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <url in tailnet>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}
	tailnetURL := flag.Arg(0)

	s := new(tsnet.Server)
	defer s.Close()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	cli := s.HTTPClient()

	resp, err := cli.Get(tailnetURL)
	if err != nil {
		log.Fatal(err)
	}

	resp.Write(os.Stdout)
}
