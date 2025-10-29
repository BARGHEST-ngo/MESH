// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build !linux || ts_omit_listenrawdisco

package magicsock

import (
	"errors"
	"fmt"
	"io"
)

func (c *Conn) listenRawDisco(family string) (io.Closer, error) {
	return nil, fmt.Errorf("raw disco listening not supported on this OS: %w", errors.ErrUnsupported)
}
