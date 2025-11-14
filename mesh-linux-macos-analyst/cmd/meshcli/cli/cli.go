// Copyright (c) Your Name & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/net/idna"
	"tailscale.com/client/local"
	"tailscale.com/feature"
	"tailscale.com/ipn"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/paths"
	"tailscale.com/util/dnsname"
)

var Stderr io.Writer = os.Stderr
var Stdout io.Writer = os.Stdout

func errf(format string, a ...any) {
	fmt.Fprintf(Stderr, format, a...)
}

func printf(format string, a ...any) {
	fmt.Fprintf(Stdout, format, a...)
}

func outln(a ...any) {
	fmt.Fprintln(Stdout, a...)
}

func newFlagSet(name string) *flag.FlagSet {
	onError := flag.ExitOnError
	if runtime.GOOS == "js" {
		onError = flag.ContinueOnError
	}
	fs := flag.NewFlagSet(name, onError)
	fs.SetOutput(Stderr)
	return fs
}

// LocalClient is the client used to communicate with tailscaled
var localClient = local.Client{
	Socket: paths.DefaultTailscaledSocket(),
}

// Command function variables that may be set by init functions in other files
var (
	maybeFunnelCmd    func() *ffcli.Command
	maybeServeCmd     func() *ffcli.Command
	maybeWebCmd       func() *ffcli.Command
	fileCmd           func() *ffcli.Command
	maybeCertCmd      func() *ffcli.Command
	maybeNetlockCmd   func() *ffcli.Command
	maybeDriveCmd     func() *ffcli.Command
	maybeSysExtCmd    func() *ffcli.Command
	maybeVPNConfigCmd func() *ffcli.Command
	maybeSystrayCmd   func() *ffcli.Command
)

// Run runs the CLI. The args do not include the binary name.
func Run(args []string) (err error) {
	// Handle version shortcuts
	if len(args) == 1 {
		switch args[0] {
		case "-V", "--version":
			args = []string{"version"}
		case "help":
			args = []string{"--help"}
		}
	}

	var warnOnce sync.Once
	local.SetVersionMismatchHandler(func(clientVer, serverVer string) {
		warnOnce.Do(func() {
			fmt.Fprintf(Stderr, "Warning: client version %q != tailscaled server version %q\n", clientVer, serverVer)
		})
	})

	cfmt.Print(`
	{{[ 𝗠𝗘𝗦𝗛 ]}}::blue
	{{[ 𝐀𝐍𝐀𝐋𝐘𝐒𝐓 𝐂𝐋𝐈𝐄𝐍𝐓 ]}}::blue
	{{by Barghest.asia. No rights reserved.}}::blue
	`)

	rootCmd := newRootCmd()
	if err := rootCmd.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		if noexec := (ffcli.NoExecError{}); errors.As(err, &noexec) {
			cmd := noexec.Command
			args := cmd.FlagSet.Args()
			if len(cmd.Subcommands) > 0 {
				if len(args) > 0 {
					return fmt.Errorf("%s: unknown subcommand: %s", fullCmd(rootCmd, cmd), args[0])
				}
				subs := make([]string, 0, len(cmd.Subcommands))
				for _, sub := range cmd.Subcommands {
					subs = append(subs, sub.Name)
				}
				return fmt.Errorf("%s: missing subcommand: %s", fullCmd(rootCmd, cmd), strings.Join(subs, ", "))
			}
		}
		return err
	}

	err = rootCmd.Run(context.Background())
	if local.IsAccessDeniedError(err) && os.Getuid() != 0 && runtime.GOOS != "windows" {
		return fmt.Errorf("%v\n\nUse 'sudo meshcli %s'.\nTo not require root, use 'sudo tailscale set --operator=$USER' once", err, strings.Join(args, " "))
	}
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	return err
}

func fullCmd(root, cmd *ffcli.Command) string {
	return root.Name + " " + cmd.Name
}

func nilOrCall(f func() *ffcli.Command) *ffcli.Command {
	if f == nil {
		return nil
	}
	return f()
}

func nonNilCmds(cmds ...*ffcli.Command) []*ffcli.Command {
	var result []*ffcli.Command
	for _, cmd := range cmds {
		if cmd != nil {
			result = append(result, cmd)
		}
	}
	return result
}

// hidden is the prefix that hides subcommands and flags from --help output when
// found at the start of the subcommand's LongHelp or flag's Usage.
const hidden = "HIDDEN: "

// usageFuncNoDefaultValues is like usageFunc but doesn't print default values.
func usageFuncNoDefaultValues(c *ffcli.Command) string {
	return usageFuncOpt(c, false)
}

func usageFunc(c *ffcli.Command) string {
	return usageFuncOpt(c, true)
}

func usageFuncOpt(c *ffcli.Command, withDefaults bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "USAGE\n")
	if c.ShortUsage != "" {
		fmt.Fprintf(&b, "  %s\n", c.ShortUsage)
	} else {
		fmt.Fprintf(&b, "  %s\n", c.Name)
	}
	fmt.Fprintf(&b, "\n")

	if help := strings.TrimPrefix(c.LongHelp, hidden); help != "" {
		fmt.Fprintf(&b, "%s\n\n", help)
	}

	if c.FlagSet != nil && c.FlagSet.NFlag() > 0 {
		fmt.Fprintf(&b, "FLAGS\n")
		c.FlagSet.VisitAll(func(f *flag.Flag) {
			if strings.HasPrefix(f.Usage, hidden) {
				return
			}
			fmt.Fprintf(&b, "  -%s", f.Name)
			if withDefaults && f.DefValue != "" {
				fmt.Fprintf(&b, " (default %s)", f.DefValue)
			}
			fmt.Fprintf(&b, "\n        %s\n", f.Usage)
		})
		fmt.Fprintf(&b, "\n")
	}

	if len(c.Subcommands) > 0 {
		fmt.Fprintf(&b, "SUBCOMMANDS\n")
		for _, sub := range c.Subcommands {
			if strings.HasPrefix(sub.LongHelp, hidden) {
				continue
			}
			fmt.Fprintf(&b, "  %s\n", sub.Name)
			if sub.ShortHelp != "" {
				fmt.Fprintf(&b, "        %s\n", sub.ShortHelp)
			}
		}
		fmt.Fprintf(&b, "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// Helper functions used by various commands
func fatalf(format string, args ...any) {
	fmt.Fprintf(Stderr, format, args...)
	os.Exit(1)
}

func fixTailscaledConnectError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to connect to tailscaled: %w", err)
}

func isRunningOrStarting(st *ipnstate.Status) (description string, ok bool) {
	switch st.BackendState {
	default:
		return fmt.Sprintf("unexpected state: %s", st.BackendState), false
	case ipn.Stopped.String():
		return "MESH is stopped.", false
	case ipn.NeedsLogin.String():
		s := "Logged out."
		if st.AuthURL != "" {
			s += fmt.Sprintf("\nLog in at: %s", st.AuthURL)
		}
		return s, false
	case ipn.NeedsMachineAuth.String():
		return "Machine is not yet approved by tailnet admin.", false
	case ipn.Running.String(), ipn.Starting.String():
		return st.BackendState, true
	}
}

// Helper functions and variables used by various commands
var hookPrintFunnelStatus feature.Hook[func(context.Context)]

func lastSeenFmt(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format("2006-01-02 15:04:05")
}

func colorableOutput() (io.Writer, bool) {
	return Stdout, false
}

func dnsOrQuoteHostname(st *ipnstate.Status, ps *ipnstate.PeerStatus) string {
	baseName := dnsname.TrimSuffix(ps.DNSName, st.MagicDNSSuffix)
	if baseName != "" {
		if strings.HasPrefix(baseName, "xn-") {
			if u, err := idna.ToUnicode(baseName); err == nil {
				return fmt.Sprintf("%s (%s)", baseName, u)
			}
		}
		return baseName
	}
	return fmt.Sprintf("(%q)", dnsname.SanitizeHostname(ps.HostName))
}

func newRootCmd() *ffcli.Command {
	rootfs := newFlagSet("meshcli")
	rootfs.Func("socket", "path to tailscaled socket", func(s string) error {
		if len(s) > 0 {
			localClient.Socket = s
		} else {
			localClient.Socket = paths.DefaultTailscaledSocket()
		}
		localClient.UseSocketOnly = true
		return nil
	})
	cfmt.Printf(`{{Using socket: %s}}::yellow|bold`, localClient.Socket)
	cfmt.Println()
	cfmt.Println()
	rootfs.Lookup("socket").DefValue = localClient.Socket

	var rootCmd *ffcli.Command
	rootCmd = &ffcli.Command{
		Name:       "meshcli",
		ShortUsage: "meshcli [flags] <subcommand> [command flags]",
		Subcommands: nonNilCmds(
			upCmd,
			downCmd,
			setCmd,
			logoutCmd,
			adbpairCmd,
			adbcollectCmd,
			configureCmd(),
			netcheckCmd,
			ipCmd,
			dnsCmd,
			statusCmd(),
			metricsCmd,
			pingCmd,
			ncCmd,
			sshCmd,
			nilOrCall(maybeFunnelCmd),
			nilOrCall(maybeServeCmd),
			versionCmd,
			nilOrCall(maybeWebCmd),
			nilOrCall(fileCmd),
			nilOrCall(maybeCertCmd),
			nilOrCall(maybeNetlockCmd),
			exitNodeCmd(),
			updateCmd,
			whoisCmd,
			debugCmd(),
			nilOrCall(maybeDriveCmd),
			idTokenCmd,
			systrayCmd,
			appcRoutesCmd,
		),
		FlagSet: rootfs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("meshcli: unknown subcommand: %s", args[0])
			}
			return flag.ErrHelp
		},
	}

	return rootCmd
}
