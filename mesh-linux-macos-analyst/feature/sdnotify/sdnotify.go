// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

/*
Package sdnotify contains a minimal wrapper around systemd-notify to enable
applications to signal readiness and status to systemd.

This package will only have effect on Linux systems running Tailscale in a
systemd unit with the Type=notify flag set. On other operating systems (or
when running in a Linux distro without being run from inside systemd) this
package will become a no-op.
*/
package sdnotify
