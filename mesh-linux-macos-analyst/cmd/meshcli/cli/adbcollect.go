// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.
//
// Portions of this code are derived from androidqf (Android Quick Forensics)
// Copyright (c) 2021–2022 Claudio Guarnieri.
// Use of this software is governed by the MVT License 1.1, available at:
// https://license.mvt.re/1.1/

package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mvt-project/androidqf_ward/acquisition"
	"github.com/mvt-project/androidqf_ward/adb"
	"github.com/mvt-project/androidqf_ward/log"
	"github.com/mvt-project/androidqf_ward/modules"
	"github.com/mvt-project/androidqf_ward/utils"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var adbcollectArgs struct {
	verbose bool
	fast    string
	list    string
	module  string
	output  string
	serial  string
	version bool
}

var adbcollectCmd = &ffcli.Command{
	Name:       "adbcollect",
	ShortUsage: "meshcli collect [flags]",
	ShortHelp:  "Collects ADB data, using AndroidQF and WARD",
	LongHelp: `The collect command initiates ADB acquision using AndroidQF and WARD libaries. The output is in the same format ready for WARD or MVT usage. 

Examples:
  meshcli adbpair HOST[:PORT] [PAIRING CODE]

`,
	FlagSet: (func() *flag.FlagSet {
		fs := newFlagSet("adbcollect")
		fs.BoolVar(&adbcollectArgs.verbose, "verbose", false, "Enable verbose output")
		fs.StringVar(&adbcollectArgs.fast, "fast", "", "Fast mode (skip some checks)")
		fs.StringVar(&adbcollectArgs.list, "list", "", "List available modules")
		fs.StringVar(&adbcollectArgs.module, "module", "", "Specific module to run")
		fs.StringVar(&adbcollectArgs.output, "output", "", "Output directory for collected data")
		fs.StringVar(&adbcollectArgs.serial, "serial", "", "Device serial number")
		fs.BoolVar(&adbcollectArgs.version, "version", false, "Show version information")
		return fs
	})(),
	Exec: runcollectCmd,
}

func validateAdbCollectArgs() error {
	return nil
}

func runcollectCmd(ctx context.Context, args []string) error {
	//(ov) my experience is generally it shouldn't take more than 60 minutes.
	//we should verify this with user feedback
	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()

	if len(args) > 0 {
		return fmt.Errorf("unexpected arguments: %v", args)
	}
	if err := validateAdbCollectArgs(); err != nil {
		return err
	}

	verbose := adbcollectArgs.verbose
	version_flag := adbcollectArgs.version
	list_modules := adbcollectArgs.list != ""
	fast := adbcollectArgs.fast != ""
	module := adbcollectArgs.module
	output_folder := adbcollectArgs.output
	serial := adbcollectArgs.serial

	if verbose {
		log.SetLogLevel(log.DEBUG)
	}
	if version_flag {
		log.Infof("AndroidQF version: %s", utils.Version)
		// include WARD version
		os.Exit(0)
	}

	if list_modules {
		mods := modules.List()
		log.Info("List of modules:")
		// include WARD modules if not combined with AndroidQF
		for _, mod := range mods {
			log.Infof("- %s", mod.Name())
		}
		os.Exit(0)
	}

	log.Debug("Starting androidqf")
	adbClient, err := adb.New()
	if err != nil {
		log.Fatal("Impossible to initialize ADB: ", err)
	}

	adb.Client = adbClient

	for {
		serial, err = adbClient.SetSerial(serial)
		if err != nil {
			log.Error(fmt.Sprintf("Error trying to connect over ADB: %s", err))

		} else {

			_, err = adbClient.GetState()
			if err == nil {
				break
			}
			log.Debug(err)
			log.Error("Unable to get device state. Please make sure it is connected and authorized. Trying again in 5 seconds...")
		}
		select {
		case <-time.After(5 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	acq, err := acquisition.New(output_folder)
	if err != nil {
		log.Debug(err)
		log.FatalExc("Impossible to initialise the acquisition", err)
	}

	// Start acquisitions
	log.Info(fmt.Sprintf("Started new acquisition in %s", acq.StoragePath))

	mods := modules.List()
	for _, mod := range mods {
		if (module != "") && (module != mod.Name()) {
			continue
		}
		err = mod.InitStorage(acq.StoragePath)
		if err != nil {
			log.Infof(
				"ERROR: failed to initialize storage for module %s: %v",
				mod.Name(),
				err,
			)
			continue
		}

		err = mod.Run(acq, fast)
		if err != nil {
			log.Infof("ERROR: failed to run module %s: %v", mod.Name(), err)
		}
	}

	if acq.StreamingMode {
		// In streaming mode, all data is already encrypted in the zip stream
		log.Info("Finalizing encrypted acquisition...")
	} else {
		// Traditional mode: hash files, then encrypt if key exists
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
	}

	acq.Complete()
	log.Info("Acquisition completed.")
	return nil
}
