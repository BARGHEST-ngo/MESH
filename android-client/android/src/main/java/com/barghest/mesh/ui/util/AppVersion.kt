// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.util

import com.barghest.mesh.BuildConfig

class AppVersion {
  companion object {
    // Returns the short version of the build version, which is what users typically expect.
    // For instance, if the build version is "1.75.80-t8fdffb8da-g2daeee584df",
    // this function returns "1.75.80".
    fun Short(): String {
      // Split the full version string by hyphen (-)
      val parts = BuildConfig.VERSION_NAME.split("-")
      // Return only the part before the first hyphen
      return parts[0]
    }
  }
}
