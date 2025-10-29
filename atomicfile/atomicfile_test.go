// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !js && !windows

package atomicfile

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDoesNotOverwriteIrregularFiles(t *testing.T) {
	// Per tailscale/tailscale#7658 as one example, almost any imagined use of
	// atomicfile.Write should likely not attempt to overwrite an irregular file
	// such as a device node, socket, or named pipe.

	const filename = "TestDoesNotOverwriteIrregularFiles"
	var path string
	// macOS private temp does not allow unix socket creation, but /tmp does.
	if runtime.GOOS == "darwin" {
		path = filepath.Join("/tmp", filename)
		t.Cleanup(func() { os.Remove(path) })
	} else {
		path = filepath.Join(t.TempDir(), filename)
	}

	// The least troublesome thing to make that is not a file is a unix socket.
	// Making a null device sadly requires root.
	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	err = WriteFile(path, []byte("hello"), 0644)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "is not a regular file") {
		t.Fatalf("unexpected error: %v", err)
	}
}
