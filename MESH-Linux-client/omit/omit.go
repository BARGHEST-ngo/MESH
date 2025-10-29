// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package omit provides consts to access Tailscale ts_omit_FOO build tags.
// They're often more convenient to eliminate some away locally with a const
// rather than using build tags.
package omit

import "errors"

// Err is an error that can be returned by functions in this package.
var Err = errors.New("feature not linked into binary per ts_omit build tag")
