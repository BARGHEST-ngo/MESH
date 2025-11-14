// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package cli

import (
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var dnsCmd = &ffcli.Command{
	Name:      "dns",
	ShortHelp: "Diagnose the internal DNS forwarder",
	LongHelp: strings.TrimSpace(`
The 'tailscale dns' subcommand provides tools for diagnosing the internal DNS
forwarder (100.100.100.100).

For more information about the DNS functionality built into Tailscale, refer to
https://tailscale.com/kb/1054/dns.
`),
	ShortUsage: strings.Join([]string{
		dnsStatusCmd.ShortUsage,
		dnsQueryCmd.ShortUsage,
	}, "\n"),
	UsageFunc: usageFuncNoDefaultValues,
	Subcommands: []*ffcli.Command{
		dnsStatusCmd,
		dnsQueryCmd,

		// TODO: implement `tailscale log` here

		// The above work is tracked in https://github.com/tailscale/tailscale/issues/13326
	},
}
