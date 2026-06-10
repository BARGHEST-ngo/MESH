// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BARGHEST-ngo/androidqf_mesh/acquisition"
	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/BARGHEST-ngo/androidqf_mesh/log"
	"github.com/BARGHEST-ngo/androidqf_mesh/modules"
	rt "github.com/botherder/go-savetime/runtime"
	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v3/ffcli"
	"tailscale.com/tailcfg"
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

func AdbPairCmd() *ffcli.Command {
	fs := flag.NewFlagSet("adbpair", flag.ContinueOnError)
	fs.BoolVar(&adbpairliteArgs.qf, "qf", false, "perform adbcollect (AndroidQF/WARD) immediately after connection")

	return &ffcli.Command{
		Name:       "adbpair",
		ShortUsage: "mesh adbpair [flags]",
		ShortHelp:  "Pair & connect to a device on the MESH network via ADB",
		FlagSet:    fs,
		Exec:       runAdbPair,
	}
}

func runAdbPair(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unexpected arguments: %v", args)
	}

	pairingArgs := PairingArgs{}

	fmt.Println("Starting automatic pairing...")

	chosenPeer, err := selectAndroidPeer(ctx)
	if err != nil {
		return fmt.Errorf("unable to select an Android client: %w", err)
	}
	pairingArgs.Host = chosenPeer.IP

	fmt.Printf("Using Android device: %s (%s)\n\n", chosenPeer.HostName, chosenPeer.IP)
	fmt.Println("On the Android device:")
	fmt.Println("1. Enable Wireless Debugging")
	fmt.Println("2. Tap 'Pair device with pairing code'")
	ReadString("Press Enter when the pairing dialog is open...")

	if err := checkDataPath(ctx, chosenPeer.IP); err != nil {
		return err
	}

	fmt.Println("Scanning for open ports...")
	openPorts, err := scanOpenPorts(chosenPeer.IP)
	if err != nil {
		return fmt.Errorf("port scan failed: %w", err)
	}
	if len(openPorts) == 0 {
		return fmt.Errorf("no open ports found on %s - is the pairing dialog open?", chosenPeer.IP)
	}

	pairingArgs.PairingCode = ReadStringWithValidation("Enter the pairing code shown on the device: ", validatePairingCode)

	adbClient, err := adb.New()
	if err != nil {
		return fmt.Errorf("failed to initialize ADB: %w", err)
	}
	adb.Client = adbClient

	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	if len(devices) > 0 {
		fmt.Printf("Found existing ADB devices: %v\n", devices)
		if err := disconnect(""); err != nil {
			return err
		}
	}

	pairedPort, err := pairWithDiscovery(&pairingArgs, openPorts)
	if err != nil {
		return err
	}

	debugPort, err := resolveDebugPort(chosenPeer.IP, openPorts, pairedPort)
	if err != nil {
		return fmt.Errorf("unable to locate debug port: %w", err)
	}
	pairingArgs.DebugPort = debugPort

	if err := connect(&pairingArgs); err != nil {
		return err
	}

	if err := validateConnect(); err != nil {
		return err
	}

	if adbpairliteArgs.qf {
		if err := qf(); err != nil {
			return err
		}
	}

	return nil
}

func checkDataPath(ctx context.Context, ip string) error {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return err
	}

	// Confirm wireguard path
	// (should never fail if device is listed in `mesh status`)
	if err := ping(ctx, addr, tailcfg.PingDisco); err != nil {
		return fmt.Errorf("device %s is not reachable on the MESH network: %w", ip, err)
	}

	// Confirm data can be exchanged
	if err := ping(ctx, addr, tailcfg.PingTSMP); err != nil {
		return fmt.Errorf("device %s is reachable on the MESH network, but the data path is broken\n"+
			"Toggle WiFi on the Android device and retry", ip)
	}

	return nil
}

func ping(ctx context.Context, addr netip.Addr, mode tailcfg.PingType) error {
	_, err := localClient.Ping(ctx, addr, mode)
	return err
}

func pairWithDiscovery(args *PairingArgs, openPorts []int) (int, error) {
	if len(openPorts) > maxAutoPairAttempts {
		fmt.Printf("Found %d open ports - too many to try automatically, please identify the pairing port:\n", len(openPorts))
		choice, ok := promptForPort(openPorts)
		if !ok {
			return 0, fmt.Errorf("invalid port selection")
		}
		args.PairPort = openPorts[choice]
		if err := pair(args); err != nil {
			return 0, err
		}
		return args.PairPort, nil
	}

	for _, port := range openPorts {
		args.PairPort = port
		fmt.Printf("Trying port %d...\n", port)
		if err := pair(args); err != nil {
			continue
		}
		return port, nil
	}
	return 0, fmt.Errorf("pairing failed on all discovered ports")
}

func pair(args *PairingArgs) error {
	if adb.Client == nil {
		return fmt.Errorf("adb not initialized")
	}

	if args == nil {
		return fmt.Errorf("invalid pairing args")
	}

	output, err := adb.Client.Exec("pair", net.JoinHostPort(args.Host, strconv.Itoa(args.PairPort)), args.PairingCode)
	if err != nil {
		return fmt.Errorf("ADB pair failed: %w\nOutput: %s", err, string(output))
	}
	fmt.Println("ADB pair successful")
	return nil
}

func resolveDebugPort(ip string, openPorts []int, pairedPort int) (int, error) {
	var candidates []int
	for _, p := range openPorts {
		if p != pairedPort {
			candidates = append(candidates, p)
		}
	}

	switch len(candidates) {
	case 0:
		fmt.Println("Scanning for debug port...")
		return selectPort(ip)
	case 1:
		fmt.Printf("Found debug port: %d\n", candidates[0])
		return candidates[0], nil
	default:
		fmt.Println("Multiple debug port candidates, please select:")
		choice, ok := promptForPort(candidates)
		if !ok {
			return 0, fmt.Errorf("invalid port selection")
		}
		return candidates[choice], nil
	}
}

func connect(args *PairingArgs) error {
	if adb.Client == nil {
		return fmt.Errorf("adb not initialized")
	}
	if args == nil {
		return fmt.Errorf("invalid pairing args")
	}
	fmt.Printf("Connecting to device...\n")
	output, err := adb.Client.Exec("connect", net.JoinHostPort(args.Host, strconv.Itoa(args.DebugPort)))
	if err != nil {
		return fmt.Errorf("ADB connect failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func validateConnect() error {
	checkADBClient()
	fmt.Printf("Validating ADB session...\n")
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}
	switch len(devices) {
	case 0:
		return fmt.Errorf("no devices connected after ADB connect")
	case 1:
		fmt.Printf("Success! Device connected\n")
		for _, d := range devices {
			if strings.HasPrefix(d, "100.") {
				fmt.Printf("Success! Valid MESH network device\n")
				fmt.Printf("You may proceed with forensics acquision\n")
				fmt.Printf("Use mesh adbcollect\n")
				return nil
			}
		}
		return fmt.Errorf("wrong device connected/paried")
	default:
		return fmt.Errorf("multiple devices connected after ADB connect: %v", devices)
	}
}

func selectAndroidPeer(ctx context.Context) (*AndroidPeer, error) {
	peers, err := getAndroidPeers(ctx)
	if err != nil || len(peers) == 0 {
		return nil, fmt.Errorf("unable to find any connected Android clients")
	}

	var chosenPeer AndroidPeer
	if len(peers) > 1 {
		choice, valid := promptForAndroidClient(peers)
		if !valid {
			return nil, fmt.Errorf("invalid client selection")
		}
		chosenPeer = peers[choice]
	} else {
		chosenPeer = peers[0]
		fmt.Printf("found 1 Android device: %s (%s)\n", chosenPeer.HostName, chosenPeer.IP)
	}

	return &chosenPeer, nil
}

func getAndroidPeers(ctx context.Context) ([]AndroidPeer, error) {
	status, err := localClient.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get MESH status: %w", err)
	}

	out := make([]AndroidPeer, 0)
	for _, p := range status.Peers() {
		peer := status.Peer[p]
		if peer.OS != "android" {
			continue
		}

		if len(peer.TailscaleIPs) == 0 {
			continue
		}

		out = append(out, AndroidPeer{
			IP:       peer.TailscaleIPs[0].String(),
			HostName: peer.HostName,
			DNSName:  peer.DNSName,
		})

	}
	return out, nil
}

func promptForAndroidClient(peers []AndroidPeer) (int, bool) {
	fmt.Println("Multiple Android clients found, select an Android client:")
	choices := make([]string, len(peers))
	for i, p := range peers {
		choices[i] = fmt.Sprintf("%s (%s)", p.HostName, p.IP)
	}

	fmt.Print("Select an Android client:\n")
	return promptForSelection(choices)
}

func promptForPort(ports []int) (int, bool) {
	fmt.Println("Multiple open ports found, select a port:")
	choices := make([]string, len(ports))
	for i, p := range ports {
		choices[i] = fmt.Sprintf("%d", p)
	}

	fmt.Print("Select a port:\n")
	return promptForSelection(choices)
}

func promptForSelection(choices []string) (int, bool) {
	for i, c := range choices {
		fmt.Printf(" %d) %s \n", i+1, c)
	}

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(choices) {
		fmt.Println("Invalid selection")
		return -1, false
	}

	return choice - 1, true
}

func selectPort(ip string) (int, error) {
	fmt.Println("Scanning for port...")
	openPorts, err := scanOpenPorts(ip)
	if err != nil {
		return -1, fmt.Errorf("port scan failed: %w", err)
	}
	if len(openPorts) == 0 {
		return -1, fmt.Errorf("no open ports found on %s", ip)
	}

	var pairPort int
	if len(openPorts) == 1 {
		pairPort = openPorts[0]
		fmt.Printf("Found port: %d\n", pairPort)
	} else {
		choice, ok := promptForPort(openPorts)
		if !ok {
			return -1, fmt.Errorf("invalid port selection")
		}

		pairPort = openPorts[choice]
	}
	return pairPort, nil
}

func scanOpenPorts(ip string) ([]int, error) {
	const workers = 1000
	const dialTimeout = 1 * time.Second

	jobs := make(chan int, workers)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var open []int

	for range workers {
		wg.Go(func() {
			for p := range jobs {
				addr := net.JoinHostPort(ip, strconv.Itoa(p))
				conn, err := net.DialTimeout("tcp", addr, dialTimeout)
				if err == nil {
					conn.Close()
					mu.Lock()
					open = append(open, p)
					mu.Unlock()
				}
			}
		})
	}

	for port := 5000; port <= 65535; port++ {
		jobs <- port
	}
	close(jobs)

	wg.Wait()
	sort.Ints(open)
	return open, nil
}

func qf() error {
	devices, err := adb.Client.Devices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	for _, d := range devices {
		if strings.HasPrefix(d, "100.") {
			fmt.Println("Performing forensics acquisition")
			fmt.Println("Starting androidqf")

			acqUUID := uuid.New().String()

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

			outputFolder := filepath.Join(rt.GetExecutableDirectory(), acqUUID)
			acq, err := acquisition.New(outputFolder)
			if err != nil {
				log.FatalExc("Impossible to initialise the acquisition", err)
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
	}
	return nil
}
