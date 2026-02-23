// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

import Cocoa
import Foundation
import Virtualization
import ArgumentParser

@main
struct HostCli: ParsableCommand {
    static var configuration = CommandConfiguration(
        abstract: "A utility for running virtual machines",
        subcommands: [Run.self],
        defaultSubcommand: Run.self)
}

var config: Config = Config()

extension HostCli {
    struct Run: ParsableCommand {
        @Option var id: String
        @Option var share: String?

        mutating func run() {
            config = Config(id)
            config.sharedDir = share
            print("Running vm with identifier \(id) and sharedDir \(share ?? "<none>")")
            _ = NSApplicationMain(CommandLine.argc, CommandLine.unsafeArgv)
        }
    }
}

