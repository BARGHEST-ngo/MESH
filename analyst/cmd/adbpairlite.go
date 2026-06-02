package cmd

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type AndroidPeer struct {
	IP       string
	HostName string
	DNSName  string
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

		fmt.Printf("chosen Android device: <%s> (%s)\n", chosenPeer.HostName, chosenPeer.IP)
		fmt.Printf("(review these instructions!)\n\n")
		fmt.Println("prompt user to allow wireless debugging and open the pairing dialog...")
	}

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
	fmt.Println("Available clients:")
	for i, p := range peers {
		fmt.Printf("  %d) %s (%s)\n", i+1, p.HostName, p.IP)
	}

	fmt.Print("\nSelect an Android client: ")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(peers) {
		fmt.Println("Invalid selection")
		return -1, false
	}

	return choice - 1, true
}

func scanOpenPorts(ctx context.Context, ip string) ([]int, error) {
	const workers = 1000
	const dialTimeout = 500 * time.Millisecond

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

			addr := fmt.Sprintf("%s:%d", ip, p)
			conn, err := net.DialTimeout("tcp", addr, dialTimeout)
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
