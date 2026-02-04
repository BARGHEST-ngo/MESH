// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/mvt-project/androidqf_ward/adb"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var adbdisableArgs struct {
	serial string
}

var adbdisableCmd = &ffcli.Command{
	Name:       "adbdisable",
	ShortUsage: "meshcli adbdisable [flags]",
	ShortHelp:  "The ADB disable utility disables the developer mode on the Android node",
	LongHelp: `Since leaving ADB developer mode open is dangerous in light on non-consensual forensics, ADB disable utility allows for an analyst to disable developer mode entirely on the android node they are analyzing.

Examples:
  meshcli adbdisable -s devicename

`,
	FlagSet: (func() *flag.FlagSet {
		fs := newFlagSet("adbpair")
		fs.StringVar(&adbdisableArgs.serial, "serial", "", "serial of the device")
		return fs
	})(),
	Exec: runadbdisable,
}

func validateDisable() error {
	printf("Validating ADB disablement...\n")
	if adb.Client == nil {
		return fmt.Errorf("ADB client not initialized")
	}
	devices, err := adb.Client.Devices()

	if err != nil {
		return fmt.Errorf("Unable to get devices: %w", err)
	}
	if len(devices) == 0 {
		printf("No devices connected - ADB disablement successful\n")
		return nil
	}
	for _, d := range devices {
		if strings.Contains(d, adbdisableArgs.serial) {
			if strings.Contains(d, "offline") {
				printf("ADB disablement successful - device is offline\n")
				return nil
			} else {
				return fmt.Errorf("ADB disablement failed: device still online")
			}
		}
	}
	printf("Device not found in ADB devices list - ADB disablement successful\n")
	return nil
}

func runadbdisable(ctx context.Context, args []string) error {
	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("impossible to initialize ADB: %v", err)
	}
	adb.Client = adbClient
	printf("Disabling ADB on Android device\n")
	out, err := adb.Client.Shell("settings", "put", "global", "development_settings_enabled", "0")
	if err != nil {
		return fmt.Errorf("failed to run `adb shell settings put global development_settings_enabled 0`: %v", err)
	}
	printf("ADB disabled: %s\n", out)
	if err := validateDisable(); err != nil {
		return err
	}
	return nil
}
