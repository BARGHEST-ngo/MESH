package cmd

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/BARGHEST-ngo/androidqf_mesh/adb"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type AndroidPeer struct {
	IP       string
	HostName string
	DNSName  string
}

type AdbPairArgs struct {
	Host        string
	PairPort    int
	DebugPort   int
	PairingCode string
}

const maxAutoPairAttempts = 5

func AdbPairLiteCmd() *ffcli.Command {
	fs := flag.NewFlagSet("adbpairlite", flag.ContinueOnError)

	return &ffcli.Command{
		Name:       "adbpairlite",
		ShortUsage: "mesh adbpairlite [flags]",
		FlagSet:    fs,
		Exec:       runAdbPairLite,
	}
}

func runAdbPairLite(ctx context.Context, args []string) error {
	ensureAdbConf()
	pairingArgs := AdbPairArgs{}
	//TODO: Allow manual pairing via command flags
	if len(args) == 0 {
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
			if err := disc(""); err != nil {
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

		fmt.Printf("Success! Valid MESH network device\n")
		fmt.Printf("You may proceed with forensics acquisition\n")
	}

	return nil
}

func pairWithDiscovery(args *AdbPairArgs, openPorts []int) (int, error) {
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

func pair(args *AdbPairArgs) error {
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
	if err := saveHost(args.Host); err != nil {
		return fmt.Errorf("failed to save host config: %w", err)
	}
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

func connect(args *AdbPairArgs) error {
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

	if err := saveHostport(strconv.Itoa(args.DebugPort)); err != nil {
		return fmt.Errorf("failed to save hostport config: %w", err)
	}
	return nil
}

func disc(serial string) error {
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
