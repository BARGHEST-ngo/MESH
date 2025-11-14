// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !unix

package main

import (
	"errors"
	"os"
	"os/exec"
)

func doExec(cmd string, args []string, env []string) error {
	c := exec.Command(cmd, args[1:]...)
	c.Env = env
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	// Propagate ExitErrors within this func to give us similar semantics to
	// the Unix variant.
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}

	return err
}
