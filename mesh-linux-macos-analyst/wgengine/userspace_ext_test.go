// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package wgengine_test

import (
	"testing"

	"github.com/amnezia-vpn/amneziawg-go/tun"
	"tailscale.com/net/tstun"
	"tailscale.com/tsd"
	"tailscale.com/tstest"
	"tailscale.com/types/logger"
	"tailscale.com/wgengine"
	"tailscale.com/wgengine/router"
)

func TestIsNetstack(t *testing.T) {
	sys := tsd.NewSystem()
	e, err := wgengine.NewUserspaceEngine(
		tstest.WhileTestRunningLogger(t),
		wgengine.Config{
			SetSubsystem:  sys.Set,
			HealthTracker: sys.HealthTracker.Get(),
			Metrics:       sys.UserMetricsRegistry(),
			EventBus:      sys.Bus.Get(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	if !sys.IsNetstack() {
		t.Errorf("IsNetstack = false; want true")
	}
}

func TestIsNetstackRouter(t *testing.T) {
	tests := []struct {
		name              string
		conf              wgengine.Config
		setNetstackRouter bool
		want              bool
	}{
		{
			name: "no_netstack",
			conf: wgengine.Config{
				Tun:    newFakeOSTUN(),
				Router: newFakeOSRouter(),
			},
			want: false,
		},
		{
			name: "netstack",
			conf: wgengine.Config{},
			want: true,
		},
		{
			name: "hybrid_netstack",
			conf: wgengine.Config{
				Tun:    newFakeOSTUN(),
				Router: newFakeOSRouter(),
			},
			setNetstackRouter: true,
			want:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := tsd.NewSystem()
			if tt.setNetstackRouter {
				sys.NetstackRouter.Set(true)
			}
			conf := tt.conf
			conf.SetSubsystem = sys.Set
			conf.HealthTracker = sys.HealthTracker.Get()
			conf.Metrics = sys.UserMetricsRegistry()
			conf.EventBus = sys.Bus.Get()
			e, err := wgengine.NewUserspaceEngine(logger.Discard, conf)
			if err != nil {
				t.Fatal(err)
			}
			defer e.Close()
			if got := sys.IsNetstackRouter(); got != tt.want {
				t.Errorf("IsNetstackRouter = %v; want %v", got, tt.want)
			}
		})
	}
}

func newFakeOSRouter() router.Router {
	return someRandoOSRouter{router.NewFake(logger.Discard)}
}

type someRandoOSRouter struct {
	router.Router
}

func newFakeOSTUN() tun.Device {
	return someRandoOSTUN{tstun.NewFake()}
}

type someRandoOSTUN struct {
	tun.Device
}

// Name returns something that is not FakeTUN.
func (t someRandoOSTUN) Name() (string, error) { return "some_os_tun0", nil }
