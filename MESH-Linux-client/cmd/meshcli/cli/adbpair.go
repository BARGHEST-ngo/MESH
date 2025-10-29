// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	rt "github.com/botherder/go-savetime/runtime"
	"github.com/google/uuid"
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
	LongHelp: `The adbpair command initiates a ADB pairing session with a device on the MESH network so you can issue ADB commands. After pairing, it will connect you to the device ready for ADB sessions.
	You can use the --qf flag to perform quick forensics immediately on connecting. 

Examples:
  meshcli adbpair HOST[:PORT] [PAIRING CODE]

`,
	FlagSet: (func() *flag.FlagSet {
		fs := newFlagSet("adbpair")
		fs.StringVar(&adbpairArgs.host, "host", "", "IP of the Android device")
		fs.StringVar(&adbpairArgs.hostport, "hostport", "", "port number for the device")
		fs.StringVar(&adbpairArgs.pairport, "pairport", "", "port number assigned during 'Pair device with pairing code'")
		fs.StringVar(&adbpairArgs.code, "code", "", "pairing code from Android device")
		fs.BoolVar(&adbpairArgs.qf, "qf", false, "perform adbcollect (AndroidQF/WARD) on connection")
		return fs
	})(),
	Exec: runadbpair,
}

//TODO: adbpairArgs.host should really be defined by the tailscaled network later down the line

func runadbpair(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unexpected arguments: %v", args)
	}
	if err := validateAdbArgs(); err != nil {
		return err
	}
	if err := adbPair(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := adbConnect(); err != nil {
		return err
	}
	if err := validateConnect(); err != nil {
		return err
	}
	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("failed to initialize ADB: %v", err)
	}
	adb.Client = adbClient
	if adbpairArgs.qf {
		if err := qf(); err != nil {
			return err
		}
	}
	return nil
}

func validateAdbArgs() error {
	if adbpairArgs.host == "" && adbpairArgs.hostport == "" &&
		adbpairArgs.pairport == "" && adbpairArgs.code == "" {
		return errors.New("ADBPair requires --host --hostport --pairport --code. Use --help for more. For help on where these values are for the Android device, refer to docs.meshforensics.org")
	}
	if adbpairArgs.host == "" {
		return errors.New("--host (Host IP) is required")
	}
	if adbpairArgs.hostport == "" {
		return errors.New("--hostport (Wireless Debugging port) is required")
	}
	if adbpairArgs.pairport == "" {
		return errors.New("--pairport (Pairing port) is required")
	}
	if adbpairArgs.code == "" {
		return errors.New("--code (Pairing code) is required")
	}
	return nil
}

func adbPair() error {
	fmt.Println("Pairing to device...")
	cmd := exec.Command("adb", "pair",
		fmt.Sprintf("%s:%s", adbpairArgs.host, adbpairArgs.pairport),
		adbpairArgs.code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ADB pair failed: %v\nOutput: %s", err, string(output))
	}
	fmt.Printf("ADB pair successful:\n%s", string(output))
	if err := saveHost(adbpairArgs.host); err != nil {
		return fmt.Errorf("failed to save host config: %v", err)
	}
	return nil
}

func adbConnect() error {
	fmt.Println("ADB pair successful, now connecting...")
	cmd := exec.Command("adb", "connect",
		fmt.Sprintf("%s:%s", adbpairArgs.host, adbpairArgs.hostport))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ADB connect failed: %v\nOutput: %s", err, string(output))
	}
	fmt.Printf("ADB connect successful:\n%s", string(output))
	if err := saveHostport(adbpairArgs.hostport); err != nil {
		return fmt.Errorf("failed to save hostport config: %v", err)
	}
	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("failed to initalize ADB: $v", err)
	}
	adb.Client = adbClient

	if err := validateConnect(); err != nil {
		return err
	}
	return nil
}

func validateConnect() error {
	fmt.Println("Validating ADB session...")
	if adb.Client == nil {
		return fmt.Errorf("ADB client not initialized")
	}
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	if len(devices) == 0 {
		fmt.Println("no devices connected")
	}
	if len(devices) == 1 {
		fmt.Println("Success! Device connected")
		fmt.Printf("Found %d device(s): %v\n", len(devices), devices)
		for _, d := range devices {
			if strings.HasPrefix(d, "100.") {
				fmt.Println("Success! Valid MESH network device")
				fmt.Println("You may proceed with forensics acquision")
				fmt.Println("Use ./meshcli adbcollect")
			} else {
				fmt.Println("Wrong device connected/paried")
			}

		}
	}
	return nil
}

func qf() error {
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	for _, d := range devices {
		if strings.HasPrefix(d, "100.") {
			fmt.Printf("Performing forensics acquisition\n")
			log.Debug("Starting androidqf")

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
