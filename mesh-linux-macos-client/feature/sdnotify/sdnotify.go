// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

/*
Package sdnotify contains a minimal wrapper around systemd-notify to enable
applications to signal readiness and status to systemd.

This package will only have effect on Linux systems running Tailscale in a
systemd unit with the Type=notify flag set. On other operating systems (or
when running in a Linux distro without being run from inside systemd) this
package will become a no-op.
*/
package sdnotify
