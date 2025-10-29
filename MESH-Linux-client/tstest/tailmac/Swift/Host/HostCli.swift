// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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

