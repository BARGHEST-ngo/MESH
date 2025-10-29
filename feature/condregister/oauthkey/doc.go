// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Package oauthkey registers support for OAuth key resolution
// if it's not disabled via the ts_omit_oauthkey build tag.
// Currently (2025-09-19), tailscaled does not need OAuth key
// resolution, only the CLI and tsnet do, so this package is
// pulled out separately to avoid linking OAuth packages into
// tailscaled.
package oauthkey
