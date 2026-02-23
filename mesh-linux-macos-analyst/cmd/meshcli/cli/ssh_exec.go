// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

//go:build !js && !windows

package cli

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func findSSH() (string, error) {
	return exec.LookPath("ssh")
}

func execSSH(ssh string, argv []string) error {
	if err := syscall.Exec(ssh, argv, os.Environ()); err != nil {
		return err
	}
	return errors.New("unreachable")
}
