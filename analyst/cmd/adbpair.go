// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/BARGHEST-ngo/MESH/analyst/session"
	meshnet "github.com/BARGHEST-ngo/MESHnet"
	"github.com/BARGHEST-ngo/androidqf_mesh/acquisition"
	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/BARGHEST-ngo/androidqf_mesh/log"
	"github.com/BARGHEST-ngo/androidqf_mesh/modules"
	rt "github.com/botherder/go-savetime/runtime"
	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type AndroidPeer struct {
	IP       string
	HostName string
	DNSName  string
}

type PairingArgs struct {
	Host        string
	PairPort    int
	DebugPort   int
	PairingCode string
}

const maxAutoPairAttempts = 5

var adbpairliteArgs struct {
	qf bool
}

func AdbPairCmd(s *session.Session, in *Input) *ffcli.Command {
	fs := flag.NewFlagSet("adbpair", flag.ContinueOnError)
	runQF := fs.Bool("qf", false, "perform adbcollect (AndroidQF/WARD) immediately after connection")

	return &ffcli.Command{
		Name:       "adbpair",
		ShortUsage: "adbpair [flags]",
		ShortHelp:  "Pair & connect to a device's wireless debugging via ADB",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			return runAdbPair(ctx, args, s, in, *runQF)
		},
	}
}

func runAdbPair(ctx context.Context, args []string, s *session.Session, in *Input, runQF bool) error {
	dev, err := chooseDevice(ctx, s, in)
	if err != nil {
		return err
	}

	fmt.Printf("Using device: %s\n\n", dev.Name)
	fmt.Println("On the device:")
	fmt.Println("1. Turn on Wireless debugging")
	fmt.Println("2. Tap 'Pair device with pairing code'")

	code, err := in.ReadValid(ctx, "Pairing code: ", validatePairingCode)
	if err != nil {
		return err
	}

	pairPort, err := in.ReadPort(ctx, "Pairing port (the number after ':' under the code): ")
	if err != nil {
		return err
	}

	if err := initADB(); err != nil {
		return err
	}

	pairAddr, err := s.BridgeADB(dev.ID, pairPort)
	if err != nil {
		return err
	}
	if err := pair(pairAddr, code); err != nil {
		return err
	}

	debugPort, err := in.ReadPort(ctx, "Debug port (Wireless debugging screen, under 'IP address & Port'): ")
	if err != nil {
		return err
	}
	debugAddr, err := s.BridgeADB(dev.ID, debugPort)
	if err != nil {
		return err
	}
	if err := connect(debugAddr); err != nil {
		return err
	}
	if err := validateConnect(debugAddr); err != nil {
		return err
	}

	fmt.Println("Device connected. You may proceed with the acquisition: run 'adbcollect'")
	if runQF {
		return qf(debugAddr)
	}
	return nil
}

func chooseDevice(ctx context.Context, s *session.Session, in *Input) (meshnet.ClientInfo, error) {
	devices := s.Clients()
	switch len(devices) {
	case 0:
		return meshnet.ClientInfo{}, errors.New("no device connected")
	case 1:
		return devices[0], nil
	}

	names := make([]string, len(devices))
	for i, d := range devices {
		names[i] = d.Name
	}

	i, err := in.Select(ctx, "Several devices are connected:", names)
	if err != nil {
		return meshnet.ClientInfo{}, err
	}
	return devices[i], nil
}

func initADB() error {
	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("failed to initialise ADB: %w", err)
	}
	adb.Client = adbClient

	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	if len(devices) > 0 {
		fmt.Printf("found existing ADB devices: %v\n", devices)
		return disconnect("")
	}
	return nil
}

func pair(addr, code string) error {
	output, err := adb.Client.Exec("pair", addr, code)
	if err != nil {
		return fmt.Errorf("ADB pair failed: %w\nOutput: %s", err, string(output))
	}
	fmt.Println("ADB pair successful")
	return nil
}

func connect(addr string) error {
	fmt.Println("Connecting to device...")
	output, err := adb.Client.Exec("connect", addr)
	if err != nil {
		return fmt.Errorf("ADB connect failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func validateConnect(addr string) error {
	fmt.Println("Validating ADB session...")
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	if !slices.Contains(devices, addr) {
		return fmt.Errorf("device not connected through MESH (adb sees: %v)", devices)
	}
	return nil
}

func qf(serial string) error {
	fmt.Println("Performing forensics acquisition")
	fmt.Println("Starting androidqf")

	acqUUID := uuid.New().String()

	for {
		var err error
		serial, err = adb.Client.SetSerial(serial)
		if err != nil {
			log.Error(fmt.Sprintf("Error trying to connect over ADB: %s", err))
		} else {
			if _, err = adb.Client.GetState(); err == nil {
				break
			}
			log.Debug(err)
			log.Error("Unable to get device state. Please make sure it is connected and authorized. Trying again in 5 seconds...")
		}
		time.Sleep(5 * time.Second)
	}

	outputFolder := filepath.Join(rt.GetExecutableDirectory(), acqUUID)
	acq, err := acquisition.New(outputFolder)
	if err != nil {
		return fmt.Errorf("impossible to initialise the acquisition: %w", err)
	}

	log.Info(fmt.Sprintf("Started new acquisition in %s", acq.StoragePath))

	for _, mod := range modules.List() {
		if err := mod.InitStorage(acq.StoragePath); err != nil {
			log.Infof("ERROR: failed to initialize storage for module %s: %v", mod.Name(), err)
			continue
		}
		if err := mod.Run(acq, false); err != nil {
			log.Infof("ERROR: failed to run module %s: %v", mod.Name(), err)
		}
	}

	if err := acq.HashFiles(); err != nil {
		log.ErrorExc("Failed to generate list of file hashes", err)
		return err
	}

	acq.StoreInfo()

	if err := acq.StoreSecurely(); err != nil {
		log.ErrorExc("Something failed while encrypting the acquisition", err)
		log.Warning("WARNING: The secure storage of the acquisition folder failed! The data is unencrypted!")
	}

	acq.Complete()
	log.Info("Acquisition completed.")
	return nil
}
