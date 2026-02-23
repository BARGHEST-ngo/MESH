// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.
package com.barghest.mesh.util

object FeatureFlags {

  // Map to hold the feature flags
  private val flags: MutableMap<String, Boolean> = mutableMapOf()

  fun initialize(defaults: Map<String, Boolean>) {
    flags.clear()
    flags.putAll(defaults)
  }

  fun enable(feature: String) {
    flags[feature] = true
  }

  fun disable(feature: String) {
    flags[feature] = false
  }

  fun isEnabled(feature: String): Boolean {
    return flags[feature] ?: false
  }
}
