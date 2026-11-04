// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func AdbcleanCmd(s *session.Session, in *Input) *ffcli.Command {

	return &ffcli.Command{
		Name:       "adbclean",
		ShortUsage: "adbclean",
		ShortHelp:  "Clears as much logging as possible from the session",
		LongHelp: `In situations where you want to clean up after a session (for removing traces of forensics in case of seizure), adbclean will clear the logcat ring buffers and reset shell-accessible telemetry counters.

[IMPORTANT] This should only be used AFTER a forensics backup. If used before, you will tamper with evidence.

This is best-effort scrub of shell-accessible state only. Paired ADB key, wireless debugging history, and any root-only logs (tombstones, dropbox, pstore) are NOT removed and remain recoverable by a forensic examiner — only factory reset clears the paired-keys store.
`,
		Exec: func(ctx context.Context, args []string) error {
			return runcleanCmd(ctx, args, s, in)
		},
	}
}

func runcleanCmd(ctx context.Context, args []string, s *session.Session, in *Input) error {
	dev, err := chooseDevice(ctx, s, in)
	if err != nil {
		return err
	}

	addr := s.ADBAddr(dev.ID)
	if addr == "" {
		return errors.New("device not connected over ADB, run adbpair")
	}
	if err := initADB(); err != nil {
		return err
	}
	adb.Client.Serial = addr

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
