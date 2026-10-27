// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package session

import (
	"fmt"
	"io"
	"sync"

	meshnet "github.com/BARGHEST-ngo/MESHnet"
)

type Session struct {
	srv *meshnet.Server
	out io.Writer

	mu      sync.Mutex
	devices map[string]*device
}

type device struct {
	info meshnet.ClientInfo
}

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

func (s *Session) connected(c meshnet.ClientInfo) {
	s.mu.Lock()
	s.devices[c.ID] = &device{info: c}
	s.mu.Unlock()
	fmt.Fprintf(s.out, "device connected: %s\n", c.Name)
}

func (s *Session) disconnected(c meshnet.ClientInfo, err error) {
	s.mu.Lock()
	delete(s.devices, c.ID)
	s.mu.Unlock()
	if err != nil {
		fmt.Fprintf(s.out, "device disconnected: %s (%v)\n", c.Name, err)
		return
	}
	fmt.Fprintf(s.out, "device disconnected: %s\n", c.Name)
}
