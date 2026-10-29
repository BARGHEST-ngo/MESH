// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	meshnet "github.com/BARGHEST-ngo/MESHnet"
	"github.com/tailscale/tailcat"
)

type Session struct {
	srv *meshnet.Server
	out io.Writer

	mu      sync.Mutex
	devices map[string]*device
}

type device struct {
	info    meshnet.ClientInfo
	bridges []net.Listener
}

const streamOpenTimeout = 10 * time.Second

func New(out io.Writer) *Session {
	s := &Session{
		devices: make(map[string]*device),
		out:     out,
	}
	s.srv = meshnet.NewServer(meshnet.ServerConfig{
		OnClientConnected:    s.connected,
		OnClientDisconnected: s.disconnected,
	})
	return s
}

func (s *Session) Start() (string, error) {
	return s.srv.Start()
}

func (s *Session) Stop() error {
	return s.srv.Stop()
}

func (s *Session) NewPairing() (string, error) {
	return s.srv.NewPairing()
}

func (s *Session) Clients() []meshnet.ClientInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]meshnet.ClientInfo, 0)
	for _, d := range s.devices {
		out = append(out, meshnet.ClientInfo{
			ID:   d.info.ID,
			Name: d.info.Name,
		})
	}
	return out
}

// BridgeADB lets local adb reach a port on the device.
// It listens on a local port and turns each connection into an "adb:<devicePort>" stream to the device
// The bridge closes when the device disconnects or the session stops.
func (s *Session) BridgeADB(clientID string, devicePort int) (localAddr string, err error) {
	if devicePort < 1 || devicePort > 65535 {
		return "", fmt.Errorf("invalid device port %d", devicePort)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("adb bridge: %w", err)
	}

	s.mu.Lock()
	d := s.devices[clientID]
	if d == nil {
		s.mu.Unlock()
		ln.Close()
		return "", fmt.Errorf("device %s is not connected", clientID)
	}
	d.bridges = append(d.bridges, ln)
	s.mu.Unlock()

	go s.serveBridge(ln, clientID, fmt.Sprintf("adb:%d", devicePort))
	return ln.Addr().String(), nil
}

func (s *Session) serveBridge(ln net.Listener, clientId, target string) {
	for {
		local, err := ln.Accept()
		if err != nil {
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), streamOpenTimeout)
			stream, err := s.srv.RequestStream(ctx, clientId, target)
			cancel()
			if err != nil {
				fmt.Fprintf(s.out, "adb bridge: %v\n", err)
				local.Close()
				return
			}
			tailcat.ProxyConns(local, stream)
		}()
	}
}

func (s *Session) connected(c meshnet.ClientInfo) {
	s.mu.Lock()
	s.devices[c.ID] = &device{info: c}
	s.mu.Unlock()
	fmt.Fprintf(s.out, "device connected: %s\n", c.Name)
}

func (s *Session) disconnected(c meshnet.ClientInfo, err error) {
	s.mu.Lock()
	d := s.devices[c.ID]
	delete(s.devices, c.ID)
	s.mu.Unlock()
	if d != nil {
		for _, ln := range d.bridges {
			ln.Close()
		}
	}

	if err != nil {
		fmt.Fprintf(s.out, "device disconnected: %s (%v)\n", c.Name, err)
		return
	}
	fmt.Fprintf(s.out, "device disconnected: %s\n", c.Name)
}
