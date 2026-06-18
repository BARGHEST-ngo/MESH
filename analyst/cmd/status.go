// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/toqueteos/webbrowser"
	"tailscale.com/client/local"
	"tailscale.com/net/netmon"
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

var localClient local.Client

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
	if len(args) > 0 {
		return errors.New("unexpected non-flag arguments to 'meshcli status'")
	}

	getStatus := localClient.Status
	if !statusArgs.peers {
		getStatus = localClient.StatusWithoutPeers
	}
	st, err := getStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if statusArgs.json {
		if statusArgs.active {
			for peer, ps := range st.Peer {
				if !ps.Active {
					delete(st.Peer, peer)
				}
			}
		}
		j, err := json.MarshalIndent(st, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s", j)
		return nil
	}

	if statusArgs.web {
		ln, err := net.Listen("tcp", statusArgs.listen)
		if err != nil {
			return err
		}
		statusURL := netmon.HTTPOfListener(ln)
		fmt.Fprintf(os.Stdout, "Serving MESH status at %v ...\n", statusURL)
		go func() {
			<-ctx.Done()
			ln.Close()
		}()
		if statusArgs.browser {
			go webbrowser.Open(statusURL)
		}
		err = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI != "/" {
				http.NotFound(w, r)
				return
			}
			st, err := localClient.Status(ctx)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			st.WriteHTML(w)
		}))
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	if len(st.Health) > 0 {
		fmt.Fprintf(os.Stdout, "# Health check:\n")
		for _, m := range st.Health {
			fmt.Fprintf(os.Stdout, "#     - %s\n", m)
		}
		fmt.Fprintf(os.Stdout, "\n")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "IP\tDNS Name\tOS\tRelay\tHostname\n")

	if statusArgs.self && st.Self != nil {
		ip := ""
		if len(st.Self.TailscaleIPs) > 0 {
			ip = st.Self.TailscaleIPs[0].String()
		}
		fmt.Fprintf(w, "*%s\t%s\t%s\t%s\t%s\n",
			ip, st.Self.DNSName, "-", "-", st.Self.HostName)
	}

	if statusArgs.peers {
		for _, peerKey := range st.Peers() {
			peer := st.Peer[peerKey]
			if statusArgs.active && !peer.Active {
				continue
			}

			ip := ""
			if len(peer.TailscaleIPs) > 0 {
				ip = peer.TailscaleIPs[0].String()
			}

			relay := peer.CurAddr
			if relay == "" {
				relay = peer.Relay
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				ip, peer.DNSName, peer.OS, relay, peer.HostName)
		}
		w.Flush()
	}

	return nil
}
