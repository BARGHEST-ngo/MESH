// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build ios || android || ts_omit_webclient

package ipnlocal

import (
	"errors"
	"net"
)

const webClientPort = 5252

type webClient struct{}

func (b *LocalBackend) ConfigureWebClient(any) {}

func (b *LocalBackend) webClientGetOrInit() error {
	return errors.New("not implemented")
}

func (b *LocalBackend) webClientShutdown() {}

func (b *LocalBackend) handleWebClientConn(c net.Conn) error {
	return errors.New("not implemented")
}
func (b *LocalBackend) updateWebClientListenersLocked() {}
