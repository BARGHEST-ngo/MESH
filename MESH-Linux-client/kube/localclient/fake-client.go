// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package localclient

import (
	"context"
	"fmt"

	"tailscale.com/ipn"
)

type FakeLocalClient struct {
	FakeIPNBusWatcher
}

func (f *FakeLocalClient) WatchIPNBus(ctx context.Context, mask ipn.NotifyWatchOpt) (IPNBusWatcher, error) {
	return &f.FakeIPNBusWatcher, nil
}

func (f *FakeLocalClient) CertPair(ctx context.Context, domain string) ([]byte, []byte, error) {
	return nil, nil, fmt.Errorf("CertPair not implemented")
}

type FakeIPNBusWatcher struct {
	NotifyChan chan ipn.Notify
}

func (f *FakeIPNBusWatcher) Close() error {
	return nil
}

func (f *FakeIPNBusWatcher) Next() (ipn.Notify, error) {
	return <-f.NotifyChan, nil
}
