// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

//go:build (ts_aws || (linux && (arm64 || amd64) && !android)) && !ts_omit_aws

package condregister

import _ "tailscale.com/ipn/store/awsstore"
