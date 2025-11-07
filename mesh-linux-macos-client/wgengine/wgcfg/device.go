// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package wgcfg

import (
	"errors"
	"io"
	"sort"

	"github.com/amnezia-vpn/amneziawg-go/conn"
	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/tun"
	"tailscale.com/types/logger"
)

// NewDevice returns a wireguard-go Device configured for Tailscale use.
func NewDevice(tunDev tun.Device, bind conn.Bind, logger *device.Logger) *device.Device {
	ret := device.NewDevice(tunDev, bind, logger)
	if ret == nil {
		return nil
	}
	ret.DisableSomeRoamingForBrokenMobileSemantics()
	return ret
}

func DeviceConfig(d *device.Device) (*Config, error) {
	r, w := io.Pipe()
	errc := make(chan error, 1)
	go func() {
		errc <- d.IpcGetOperation(w)
		w.Close()
	}()
	cfg, fromErr := FromUAPI(r)
	r.Close()
	getErr := <-errc
	err := errors.Join(getErr, fromErr)
	if err != nil {
		return nil, err
	}
	sort.Slice(cfg.Peers, func(i, j int) bool {
		return cfg.Peers[i].PublicKey.Less(cfg.Peers[j].PublicKey)
	})
	return cfg, nil
}

// ReconfigDevice replaces the existing device configuration with cfg.
func ReconfigDevice(d *device.Device, cfg *Config, logf logger.Logf) (err error) {
	defer func() {
		if err != nil {
			logf("wgcfg.Reconfig failed: %v", err)
		}
	}()

	prev, err := DeviceConfig(d)
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	errc := make(chan error, 1)
	go func() {
		errc <- d.IpcSetOperation(r)
		r.Close()
	}()

	toErr := cfg.ToUAPI(logf, w, prev)
	w.Close()
	setErr := <-errc
	return errors.Join(setErr, toErr)
}
