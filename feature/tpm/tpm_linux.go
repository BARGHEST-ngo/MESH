// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package tpm

import (
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
)

func open() (transport.TPMCloser, error) {
	tpm, err := linuxtpm.Open("/dev/tpmrm0")
	if err == nil {
		return tpm, nil
	}
	return linuxtpm.Open("/dev/tpm0")
}
