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

import com.barghest.mesh.ui.model.Ipn

class AdvertisedRoutesHelper {
  companion object {
    fun exitNodeOnFromPrefs(prefs: Ipn.Prefs): Boolean {
      var v4 = false
      var v6 = false
      prefs.AdvertiseRoutes?.forEach {
        if (it == "0.0.0.0/0") {
          v4 = true
        }
        if (it == "::/0") {
          v6 = true
        }
      }
      return v4 && v6
    }
  }
}
