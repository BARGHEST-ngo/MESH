// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var adbcleanArgs struct {
	serial string
}

func AdbcleanCmd() *ffcli.Command {
	fs := flag.NewFlagSet("adbclean", flag.ContinueOnError)
	fs.StringVar(&adbcleanArgs.serial, "serial", "", "Device serial number")

	return &ffcli.Command{
		Name:       "adbclean",
		ShortUsage: "mesh adbclean [flags]",
		ShortHelp:  "Clears as much logging as possible from the session",
		LongHelp: `In situations where you want to clean up after a session (for removing traces of forensics in case of seizure), adbclean will clear the logcat ring buffers and reset shell-accessible telemetry counters.

[IMPORTANT] This should only be used AFTER a forensics backup. If used before, you will tamper with evidence.

This is best-effort scrub of shell-accessible state only. Paired ADB key, wireless debugging history, and any root-only logs (tombstones, dropbox, pstore) are NOT removed and remain recoverable by a forensic examiner — only factory reset clears the paired-keys store.

Examples:
  mesh adbclean
  mesh adbclean --serial devicename
`,
		FlagSet: fs,
		Exec:    runcleanCmd,
	}
}

func runcleanCmd(ctx context.Context, args []string) error {
	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("impossible to initialize ADB: %v", err)
	}
	adb.Client = adbClient
	if adbcleanArgs.serial != "" {
		adb.Client.Serial = adbcleanArgs.serial
	}

	fmt.Printf("Resetting battery stats history\n")
	out, err := adb.Client.Shell("dumpsys", "batterystats", "--reset")
	if err != nil {
		fmt.Printf("failed to reset batterystats (continuing): %v\n", err)
	} else {
		fmt.Printf("batterystats reset: %s\n", out)
	}

	fmt.Printf("Clearing process stats\n")
	out, err = adb.Client.Shell("dumpsys", "procstats", "--clear")
	if err != nil {
		fmt.Printf("failed to clear procstats (continuing): %v\n", err)
	} else {
		fmt.Printf("procstats cleared: %s\n", out)
	}

	fmt.Printf("Clearing statsd puller cache\n")
	out, err = adb.Client.Shell("cmd", "stats", "clear-puller-cache")
	if err != nil {
		fmt.Printf("failed to clear statsd puller cache (continuing): %v\n", err)
	} else {
		fmt.Printf("statsd puller cache cleared: %s\n", out)
	}

	fmt.Printf("Clearing logcat ring buffers\n")
	out, err = adb.Client.Shell("logcat", "-b", "all", "-c")
	if err != nil {
		return fmt.Errorf("failed to run `adb shell logcat -b all -c`: %v", err)
	}
	fmt.Printf("logcat cleared: %s\n", out)

	fmt.Printf("adbclean complete. Root-only artifacts remain — see --help.\n")
	return nil
}
