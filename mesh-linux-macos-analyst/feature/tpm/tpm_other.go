// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux && !windows

package tpm

import (
	"errors"

	"github.com/google/go-tpm/tpm2/transport"
)

func open() (transport.TPMCloser, error) {
	return nil, errors.New("TPM not supported on this platform")
}
