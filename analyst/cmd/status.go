// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var statusArgs struct {
	json    bool
	web     bool
	active  bool
	self    bool
	peers   bool
	listen  string
	browser bool
}

// sanitizeForTerminal strips ASCII control characters, including ESC, from
// peer-supplied strings before they are written to a terminal.
func sanitizeForTerminal(s string) string {
	return fmt.Sprintf("%q", s)
}

func StatusCmd() *ffcli.Command {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.BoolVar(&statusArgs.json, "json", false, "output in JSON format")
	fs.BoolVar(&statusArgs.web, "web", false, "run webserver with HTML showing status")
	fs.BoolVar(&statusArgs.active, "active", false, "filter output to only peers with active sessions")
	fs.BoolVar(&statusArgs.self, "self", true, "show status of local machine")
	fs.BoolVar(&statusArgs.peers, "peers", true, "show status of peers")
	fs.StringVar(&statusArgs.listen, "listen", "127.0.0.1:8384", "listen address for web mode")
	fs.BoolVar(&statusArgs.browser, "browser", true, "open a browser in web mode")

	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "meshcli status [--active] [--web] [--json]",
		ShortHelp:  "Show state of MESH network and its connections",
		LongHelp: strings.TrimSpace(`
Shows the current status of the MESH daemon and its connections.

By default, shows a human-readable summary of the current state.
Use --json for machine-readable output.
Use --web to start a local web server showing the status.
Use --active to show only peers with active sessions.
`),
		FlagSet: fs,
		Exec:    runStatus,
	}
}

func runStatus(ctx context.Context, args []string) error {
	// Tailscale removed, will review flags when adding MESHnet support
	return nil
}
