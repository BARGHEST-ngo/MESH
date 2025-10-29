// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

import Foundation

struct Notifications {
    // Stops the virtual machine and saves its state
    static var stop = Notification.Name("io.tailscale.macvmhost.stop")

    // Pauses the virtual machine and exits without saving its state
    static var halt = Notification.Name("io.tailscale.macvmhost.halt")
}
