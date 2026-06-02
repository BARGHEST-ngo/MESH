package cmd

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

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
	if len(args) == 0 {
		fmt.Println("Starting automatic pairing...")

		peers, err := getAndroidPeers(ctx)
		if err != nil || len(peers) == 0 {
			return fmt.Errorf("unable to find any connected Android clients")
		}

		var chosenPeer AndroidPeer
		if len(peers) > 1 {
			choice, valid := promptForAndroidClient(peers)
			if !valid {
				return fmt.Errorf("invalid client selection")
			}
			chosenPeer = peers[choice]
		} else {
			chosenPeer = peers[0]
			fmt.Printf("found 1 Android device: <%s> (%s)\n", chosenPeer.HostName, chosenPeer.IP)
		}

		pairingArgs.Host = chosenPeer.IP

		fmt.Printf("chosen Android device: <%s> (%s)\n", chosenPeer.HostName, chosenPeer.IP)
		fmt.Printf("(review these instructions!)\n\n")
		fmt.Println("prompt user to allow wireless debugging and open the pairing dialog...")
		ReadString("Press Enter when the pairing dialog is open...")

		fmt.Println("Scanning for pairing port...")
		openPorts, err := scanOpenPorts(chosenPeer.IP)
		openPorts = append(openPorts, 5001)
		if err != nil {
			return fmt.Errorf("port scan failed: %w", err)
		}
		if len(openPorts) == 0 {
			return fmt.Errorf("no open ports found on %s — is the pairing dialog open?", chosenPeer.IP)
		}

		var pairPort int
		if len(openPorts) == 1 {
			pairPort = openPorts[0]
			fmt.Printf("Found pairing port: %d\n", pairPort)
		} else {
			choice, ok := promptForPairingPort(openPorts)
			if !ok {
				return fmt.Errorf("invalid pairing port selection")
			}

			pairPort = openPorts[choice]
		}
		fmt.Printf("pairPort: %d\n", pairPort)
		pairingArgs.PairPort = pairPort

		pairingCode := ReadStringWithValidation("Enter the pairing code shown on the device: ", validatePairingCode)
		pairingArgs.PairingCode = pairingCode
	}

	fmt.Print(pairingArgs)

	return nil
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

func promptForPairingPort(ports []int) (int, bool) {
	fmt.Println("Multiple open ports found, select the pairing port:")
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

	fmt.Print("\nSelect an option: ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(choices) {
		fmt.Println("Invalid selection")
		return -1, false
	}

	return choice - 1, true
}

func scanOpenPorts(ip string) ([]int, error) {
	const workers = 1000
	const dialTimeout = 500 * time.Millisecond

	dialer := &net.Dialer{Timeout: dialTimeout}
	tlsCfg := &tls.Config{InsecureSkipVerify: true}

	sem := make(chan struct{}, workers)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var open []int

	for port := 5000; port <= 65535; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			addr := net.JoinHostPort(ip, strconv.Itoa(p))
			conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsCfg)
			if err == nil {
				conn.Close()
				mu.Lock()
				open = append(open, p)
				mu.Unlock()
			}
		}(port)
	}
	wg.Wait()
	sort.Ints(open)
	return open, nil
}
