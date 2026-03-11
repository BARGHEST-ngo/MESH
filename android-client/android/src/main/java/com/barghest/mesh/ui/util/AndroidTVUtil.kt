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

import android.content.pm.PackageManager
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.unit.dp
import com.barghest.mesh.UninitializedApp
import com.barghest.mesh.ui.util.AndroidTVUtil.isAndroidTV

object AndroidTVUtil {
  private val FEATURE_FIRETV = "amazon.hardware.fire_tv"

  fun isAndroidTV(): Boolean {
    val pm = UninitializedApp.get().packageManager
    return (pm.hasSystemFeature(@Suppress("deprecation") PackageManager.FEATURE_TELEVISION) ||
        pm.hasSystemFeature(PackageManager.FEATURE_LEANBACK) ||
        pm.hasSystemFeature(FEATURE_FIRETV))
  }
}

// Applies a letterbox effect iff we're running on Android TV to reduce the overall width
// of the UI.
fun Modifier.universalFit(): Modifier {
  return when (isAndroidTV()) {
    true -> this.padding(horizontal = 150.dp, vertical = 10.dp).clip(RoundedCornerShape(10.dp))
    false -> this
  }
}
