// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// execTailscale replaces the current process with the co-located tailscale
// binary, passing args directly. The tailscale binary must live in the same
// directory as the mesh-analyst binary.
func execTailscale(args []string) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	binary := filepath.Join(filepath.Dir(self), "tailscale")
	if _, err := os.Stat(binary); err != nil {
		return fmt.Errorf("tailscale binary not found at %s", binary)
	}
	return syscall.Exec(binary, append([]string{binary}, args...), os.Environ())
}
