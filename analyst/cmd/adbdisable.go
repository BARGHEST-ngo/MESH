// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func AdbdisableCmd(s *session.Session, in *Input) *ffcli.Command {
	return &ffcli.Command{
		Name:       "adbdisable",
		ShortUsage: "adbdisable",
		ShortHelp:  "The ADB disable utility disables the developer mode on the Android node",
		LongHelp:   `Since leaving ADB developer mode open is dangerous in light of non-consensual forensics, the ADB disable utility allows an analyst to disable developer mode entirely on the android node they are analyzing.`,
		Exec: func(ctx context.Context, args []string) error {
			return runadbdisable(ctx, args, s, in)
		},
	}
}

func runadbdisable(ctx context.Context, args []string, s *session.Session, in *Input) error {
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

	fmt.Printf("Turning off Developer options on %s...\n", dev.Name)
	if _, err := adb.Client.Shell("settings", "put", "global", "development_settings_enabled", "0"); err != nil {
		return fmt.Errorf("failed to disable developer options: %v", err)
	}

	if _, err := adb.Client.Shell("settings", "put", "global", "adb_wifi_enabled", "0"); err != nil {
		return fmt.Errorf("failed to turn off wireless debugging: %v", err)
	}

	if err := validateDisable(ctx, addr); err != nil {
		return err
	}
	if err := disconnect(addr); err != nil {
		return fmt.Errorf("failed to disconnect devices after ADB disable: %v", err)
	}
	s.SetADBAddr(dev.ID, "")
	fmt.Printf("Done. %s is still connected to MESH. Run 'disconnect' to end the device's MESH session\n", dev.Name)
	return nil
}

func triggerDeveloperModeCheck() {
	checkADBClient()
	adb.Client.Shell("am", "start", "-a", "android.settings.APPLICATION_DEVELOPMENT_SETTINGS")
	time.Sleep(1 * time.Second)
}

func validateDisable(ctx context.Context, addr string) error {
	checkADBClient()
	fmt.Printf("Waiting for device to drop ADB...\n")

	const maxRetries = 10
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelay):
		}

		devices, err := adb.Client.Devices()
		if err != nil {
			return fmt.Errorf("unable to get devices: %w", err)
		}

		if !slices.Contains(devices, addr) {
			fmt.Println("ADB is disabled.")
			return nil
		}

		state, err := adb.Client.GetState()

		if err != nil || state == "offline" {
			fmt.Println("ADB is disabled.")
			return nil
		}

		if attempt == 1 {
			triggerDeveloperModeCheck()
		}

		fmt.Printf("Still reachable, retrying (attempt %d/%d)...\n", attempt, maxRetries)
	}

	return fmt.Errorf("validation timeout: device still connected after %d attempts", maxRetries)
}
