// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	rt "github.com/botherder/go-savetime/runtime"
	"github.com/google/uuid"
	"github.com/mattn/go-isatty"
	"github.com/mvt-project/androidqf_ward/acquisition"
	"github.com/mvt-project/androidqf_ward/adb"
	"github.com/mvt-project/androidqf_ward/log"
	"github.com/mvt-project/androidqf_ward/modules"
	"github.com/mvt-project/androidqf_ward/utils"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var adbpairArgs struct {
	host     string
	hostport string
	pairport string
	code     string
	qf       bool
}

type Acquisition struct {
	UUID             string                          `json:"uuid"`
	AndroidQFVersion string                          `json:"androidqf_version"`
	StoragePath      string                          `json:"storage_path"`
	Started          time.Time                       `json:"started"`
	Completed        time.Time                       `json:"completed"`
	Collector        *adb.Collector                  `json:"collector"`
	TmpDir           string                          `json:"tmp_dir"`
	SdCard           string                          `json:"sdcard"`
	Cpu              string                          `json:"cpu"`
	closeLog         func()                          `json:"-"`
	EncryptedWriter  *acquisition.EncryptedZipWriter `json:"-"`
	StreamingMode    bool                            `json:"streaming_mode"`
	StreamingPuller  *acquisition.StreamingPuller    `json:"-"`
}

var adbpairCmd = &ffcli.Command{
	Name:       "adbpair",
	ShortUsage: "meshcli adbpair [flags]",
	ShortHelp:  "Pair & connect to a device on the MESH network via ADB",
	LongHelp: `The adbpair command initiates an ADB pairing session with a device on the MESH network so you can issue ADB commands. After pairing, it will connect you to the device ready for ADB sessions.

If you run the command without flags, you will be prompted interactively for each required value.

You can use the --qf flag to perform quick forensics immediately on connecting.

Examples:
  # Interactive mode (recommended for first-time users)
  meshcli adbpair

  # With all flags specified
  meshcli adbpair --host 100.64.1.5 --hostport 37891 --pairport 45281 --code 123456

  # Partial flags (will prompt for missing values)
  meshcli adbpair --host 100.64.1.5

  # Short form with host:port
  meshcli adbpair 100.64.1.5:37891

  # With quick forensics
  meshcli adbpair --qf

`,
	FlagSet: (func() *flag.FlagSet {
		fs := newFlagSet("adbpair")
		fs.StringVar(&adbpairArgs.host, "host", "", "(optional) IP address of the Android device on the MESH network")
		fs.StringVar(&adbpairArgs.hostport, "hostport", "", "(optional) Wireless Debugging port number")
		fs.StringVar(&adbpairArgs.pairport, "pairport", "", "(optional) Pairing port from 'Pair device with pairing code' dialog")
		fs.StringVar(&adbpairArgs.code, "code", "", "(optional) 6-digit pairing code from Android device")
		fs.BoolVar(&adbpairArgs.qf, "qf", false, "perform adbcollect (AndroidQF/WARD) immediately after connection")
		return fs
	})(),
	Exec: runadbpair,
}

//TODO: adbpairArgs.host should really be defined by the tailscaled network later down the line

func runadbpair(ctx context.Context, args []string) error {
	// For brevity, we also allow `meshcli HOST[:HOSTPORT]`
	if len(args) == 1 {
		if strings.Contains(args[0], ":") {
			var err error
			adbpairArgs.host, adbpairArgs.hostport, err = net.SplitHostPort(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse host:port: %v", err)
			}
		} else {
			adbpairArgs.host = args[0]
		}
	} else if len(args) > 1 {
		return fmt.Errorf("unexpected arguments: %v", args)
	}

	// Check if any arguments are missing and prompt for them interactively
	if adbpairArgs.host == "" || adbpairArgs.hostport == "" ||
		adbpairArgs.pairport == "" || adbpairArgs.code == "" {
		promptForMissingArgs()
	}

	// Validate all arguments
	if err := validateAdbArgs(); err != nil {
		return err
	}

	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("failed to initialize ADB: %v", err)
	}
	adb.Client = adbClient

	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	if len(devices) > 0 {
		printf("Found existing ADB devices: %v\n", devices)
		if err := disconnect(""); err != nil {
			return err
		}
	}

	if err := adbPair(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := adbConnect(); err != nil {
		return err
	}

	if adbpairArgs.qf {
		if err := qf(); err != nil {
			return err
		}
	}
	return nil
}

func validateAdbArgs() error {
	// Validate host if provided
	if adbpairArgs.host != "" {
		if err := validateIPAddress(adbpairArgs.host); err != nil {
			return fmt.Errorf("invalid --host: %v", err)
		}
	} else {
		return errors.New("--host (Host IP) is required")
	}

	// Validate hostport if provided
	if adbpairArgs.hostport != "" {
		if err := validatePort(adbpairArgs.hostport); err != nil {
			return fmt.Errorf("invalid --hostport: %v", err)
		}
	} else {
		return errors.New("--hostport (Wireless Debugging port) is required")
	}

	// Validate pairport if provided
	if adbpairArgs.pairport != "" {
		if err := validatePort(adbpairArgs.pairport); err != nil {
			return fmt.Errorf("invalid --pairport: %v", err)
		}
	} else {
		return errors.New("--pairport (Pairing port) is required")
	}

	// Validate code if provided
	if adbpairArgs.code != "" {
		if err := validatePairingCode(adbpairArgs.code); err != nil {
			return fmt.Errorf("invalid --code: %v", err)
		}
	} else {
		return errors.New("--code (Pairing code) is required")
	}

	return nil
}

func adbPair() error {
	checkADBClient()
	printf("Pairing to device...\n")
	output, err := adb.Client.Exec("pair",
		fmt.Sprintf("%s:%s", adbpairArgs.host, adbpairArgs.pairport),
		adbpairArgs.code)
	if err != nil {
		return fmt.Errorf("ADB pair failed: %v\nOutput: %s", err, string(output))
	}
	printf("ADB pair successful:\n%s", string(output))
	if err := saveHost(adbpairArgs.host); err != nil {
		return fmt.Errorf("failed to save host config: %v", err)
	}
	return nil
}

func adbConnect() error {
	checkADBClient()
	printf("ADB pair successful, now connecting...\n")
	output, err := adb.Client.Exec("connect",
		fmt.Sprintf("%s:%s", adbpairArgs.host, adbpairArgs.hostport))
	if err != nil {
		return fmt.Errorf("ADB connect failed: %v\nOutput: %s", err, string(output))
	}
	printf("ADB connect successful:\n%s", string(output))
	if err := saveHostport(adbpairArgs.hostport); err != nil {
		return fmt.Errorf("failed to save hostport config: %v", err)
	}
	return validateConnect()
}

func validateConnect() error {
	checkADBClient()
	printf("Validating ADB session...\n")
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	switch len(devices) {
	case 0:
		return fmt.Errorf("no devices connected after ADB connect")
	case 1:
		printf("Success! Device connected\n")
		for _, d := range devices {
			if strings.HasPrefix(d, "100.") {
				printf("Success! Valid MESH network device\n")
				printf("You may proceed with forensics acquision\n")
				printf("Use ./meshcli adbcollect\n")
				return nil
			}
		}
		return fmt.Errorf("wrong device connected/paried")
	default:
		return fmt.Errorf("multiple devices connected after ADB connect: %v", devices)
	}
}

func qf() error {
	checkADBClient()
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	for _, d := range devices {
		if strings.HasPrefix(d, "100.") {
			printf("Performing forensics acquisition\n")
			printf("Starting androidqf\n")

			acqt := Acquisition{
				UUID:             uuid.New().String(),
				Started:          time.Now().UTC(),
				AndroidQFVersion: utils.Version,
			}

			for {
				d, err = adb.Client.SetSerial(d)
				if err != nil {
					log.Error(fmt.Sprintf("Error trying to connect over ADB: %s", err))
				} else {
					_, err = adb.Client.GetState()
					if err == nil {
						break
					}
					log.Debug(err)
					log.Error("Unable to get device state. Please make sure it is connected and authorized. Trying again in 5 seconds...")
				}
				time.Sleep(5 * time.Second)
			}

			output_folder := filepath.Join(rt.GetExecutableDirectory(), acqt.UUID)
			acq, err := acquisition.New(output_folder)
			if err != nil {
				log.Debug(err)
				log.FatalExc("Impossible to initialise the acquisition", err)
			}

			log.Info(fmt.Sprintf("Started new acquisition in %s", acq.StoragePath))

			mods := modules.List()
			for _, mod := range mods {
				err = mod.InitStorage(acq.StoragePath)
				if err != nil {
					log.Infof("ERROR: failed to initialize storage for module %s: %v", mod.Name(), err)
					continue
				}

				err = mod.Run(acq, false)
				if err != nil {
					log.Infof("ERROR: failed to run module %s: %v", mod.Name(), err)
				}
			}

			err = acq.HashFiles()
			if err != nil {
				log.ErrorExc("Failed to generate list of file hashes", err)
				return err
			}

			acq.StoreInfo()

			err = acq.StoreSecurely()
			if err != nil {
				log.ErrorExc("Something failed while encrypting the acquisition", err)
				log.Warning("WARNING: The secure storage of the acquisition folder failed! The data is unencrypted!")
			}

			acq.Complete()
			log.Info("Acquisition completed.")
			return nil
		}
	}
	return nil
}

func disconnect(serial string) error {
	checkADBClient()
	if serial == "" {
		printf("Disconnecting all devices...\n")
		out, err := exec.Command(adb.Client.ExePath, "disconnect").Output()
		if err != nil {
			return fmt.Errorf("failed to disconnect all devices: %v\nOutput: %s", err, string(out))
		}
	} else {
		printf("Disconnecting device %s...\n", serial)
		out, err := exec.Command(adb.Client.ExePath, "disconnect", serial).Output()
		if err != nil {
			return fmt.Errorf("failed to disconnect device %s: %v\nOutput: %s", serial, err, string(out))
		}
	}
	return nil
}

func checkADBClient() {
	if adb.Client == nil {
		panic("ADB client not initialized")
	}
}

// validateIPAddress validates that the input is a valid IP address (IPv4 or IPv6).
func validateIPAddress(ip string) error {
	if ip == "" {
		return errors.New("IP address cannot be empty")
	}
	if _, err := netip.ParseAddr(ip); err != nil {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

// validatePort validates that the input is a valid port number (1-65535).
func validatePort(port string) error {
	if port == "" {
		return errors.New("port cannot be empty")
	}
	portNum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return errors.New("port must be a number")
	}
	if portNum < 1 || portNum > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
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

// promptForMissingArgs interactively prompts the user for any missing ADB pairing arguments.
func promptForMissingArgs() {
	if adbpairArgs.host == "" {
		adbpairArgs.host = ReadStringWithValidation("Enter target device IP address: ", validateIPAddress)
	}

	if adbpairArgs.hostport == "" {
		adbpairArgs.hostport = ReadStringWithValidation("Enter target device Wireless Debugging port: ", validatePort)
	}

	if adbpairArgs.pairport == "" {
		adbpairArgs.pairport = ReadStringWithValidation("Enter target device pairing port: ", validatePort)
	}

	if adbpairArgs.code == "" {
		adbpairArgs.code = ReadStringWithValidation("Enter pairing code: ", validatePairingCode)
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
