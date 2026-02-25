// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package libtailscale

import (
	"io"
	"log"
)

// adaptInputStream wraps an [InputStream] into an [io.ReadCloser].
// It launches a goroutine to stream reads into a pipe.
func adaptInputStream(in InputStream) io.ReadCloser {
	if in == nil {
		return nil
	}
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		for {
			b, err := in.Read()
			if err != nil {
				log.Printf("error reading from inputstream: %v", err)
				return
			}
			if b == nil {
				return
			}
			if _, err := w.Write(b); err != nil {
				log.Printf("error writing to pipe: %v", err)
				return
			}
		}
	}()
	return r
}
