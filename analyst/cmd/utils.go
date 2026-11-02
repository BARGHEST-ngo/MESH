package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/mattn/go-isatty"
)

func checkADBClient() {
	if adb.Client == nil {
		panic("ADB client not initialized")
	}
}

// ReadString prompts the user for input with the given message and returns the trimmed response.
// If there is no TTY on both Stdin and Stdout, returns an empty string.
func ReadString(msg string) string {
	if !(isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())) {
		return ""
	}
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(resp)
}

// ValidationFunc is a function that validates input and returns an error if invalid.
type ValidationFunc func(string) error

// ReadStringWithValidation prompts the user for input with the given message,
// validates it using the provided validation function, and re-prompts on validation errors.
// If there is no TTY on both Stdin and Stdout, returns an empty string.
func ReadStringWithValidation(msg string, validate ValidationFunc) string {
	if !(isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())) {
		return ""
	}
	for {
		input := ReadString(msg)
		if input == "" {
			continue
		}
		if err := validate(input); err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			continue
		}
		return input
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
