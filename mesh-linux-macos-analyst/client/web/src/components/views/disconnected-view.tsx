// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

import React from "react"
import TailscaleIcon from "src/assets/icons/tailscale-icon.svg?react"

/**
 * DisconnectedView is rendered after node logout.
 */
export default function DisconnectedView() {
  return (
    <>
      <TailscaleIcon className="mx-auto" />
      <p className="mt-12 text-center text-text-muted">
        You logged out of this device. To reconnect it you will have to
        re-authenticate the device from either the Tailscale app or the
        Tailscale command line interface.
      </p>
    </>
  )
}
