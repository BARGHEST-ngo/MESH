// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

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

func triggerDeveloperModeCheck() {
	adb.Client.Shell("am", "start", "-a", "android.settings.APPLICATION_DEVELOPMENT_SETTINGS")
	time.Sleep(1 * time.Second)
}

func validateDisable() error {
	printf("Validating ADB disablement...\n")
	if adb.Client == nil {
		return fmt.Errorf("ADB client not initialized")
	}

	const maxRetries = 10
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		time.Sleep(retryDelay)

		devices, err := adb.Client.Devices()
		if err != nil {
			return fmt.Errorf("unable to get devices: %w", err)
		}

		if len(devices) == 0 {
			printf("ADB disablement successful - no devices connected\n")
			return nil
		}

		var targetDevice string
		for _, d := range devices {
			if adbdisableArgs.serial == "" || strings.Contains(d, adbdisableArgs.serial) {
				targetDevice = d
				break
			}
		}

		if targetDevice == "" {
			printf("ADB disablement successful - device not found\n")
			return nil
		}

		adb.Client.Serial = targetDevice
		state, err := adb.Client.GetState()

		if err != nil || state == "offline" {
			printf("ADB disablement successful - device is offline\n")
			return nil
		}

		if attempt == 1 {
			triggerDeveloperModeCheck()
		}

		printf("Validating disablement (attempt %d/%d)...\n", attempt, maxRetries)
	}

	return fmt.Errorf("validation timeout: device still connected after %d attempts", maxRetries)
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
