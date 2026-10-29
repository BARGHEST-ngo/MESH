// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
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

func StatusCmd(s *session.Session) *ffcli.Command {
	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "status [--active] [--web] [--json]",
		ShortHelp:  "Show state of MESH network and its connections",
		Exec: func(ctx context.Context, args []string) error {
			return runStatus(ctx, args, s)
		},
	}
}

func runStatus(ctx context.Context, args []string, s *session.Session) error {
	pin := s.PairingPIN()
	if pin == "" {
		fmt.Println("Pairing: none, run 'pin' to pair a device.")
	} else {
		fmt.Printf("Pairing: pending (PIN: %s)\n", pin)
	}

	devices := s.Clients()
	fmt.Printf("Devices: %d connected\n", len(devices))
	for _, d := range devices {
		fmt.Printf("\t%q\t%s\n", d.Name, shortID(d.ID))
	}
	return nil
}

func shortID(id string) string {
	id = strings.TrimPrefix(id, "nodekey:")
	if len(id) > 8 {
		id = id[:8]
	}
	return id
}
