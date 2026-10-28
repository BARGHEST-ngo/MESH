package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
)

func checkADBClient() {
	if adb.Client == nil {
		panic("ADB client not initialized")
	}
}

// validatePairingCode validates that the input is a valid 6-digit pairing code.
func validatePairingCode(code string) error {
	if code == "" {
		return errors.New("pairing code cannot be empty")
	}
	// ADB pairing codes are typically 6 digits
	matched, _ := regexp.MatchString(`^\d{6}$`, code)
	if !matched {
		return errors.New("pairing code must be exactly 6 digits")
	}
	return nil
}

func disconnect(serial string) error {
	if serial == "" {
		fmt.Printf("Disconnecting all devices...\n")
		out, err := exec.Command(adb.Client.ExePath, "disconnect").Output()
		if err != nil {
			return fmt.Errorf("failed to disconnect all devices: %v\nOutput: %s", err, string(out))
		}
	} else {
		fmt.Printf("Disconnecting device %s...\n", serial)
		out, err := exec.Command(adb.Client.ExePath, "disconnect", serial).Output()
		if err != nil {
			return fmt.Errorf("failed to disconnect device %s: %v\nOutput: %s", serial, err, string(out))
		}
	}
	return nil
}
