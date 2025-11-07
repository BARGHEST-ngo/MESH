// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
