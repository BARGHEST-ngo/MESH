// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package cli

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func findSSH() (string, error) {
	// use C:\Windows\System32\OpenSSH\ssh.exe since unexpected behavior
	// occurred with ssh.exe provided by msys2/cygwin and other environments.
	if systemRoot := os.Getenv("SystemRoot"); systemRoot != "" {
		exe := filepath.Join(systemRoot, "System32", "OpenSSH", "ssh.exe")
		if st, err := os.Stat(exe); err == nil && !st.IsDir() {
			return exe, nil
		}
	}
	return exec.LookPath("ssh")
}

func execSSH(ssh string, argv []string) error {
	// Don't use syscall.Exec on Windows, it's not fully implemented.
	cmd := exec.Command(ssh, argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var ee *exec.ExitError
	err := cmd.Run()
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}
	return err
}
