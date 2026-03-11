// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.
package com.barghest.mesh

import android.content.Context
import android.os.Build
import android.util.Log

object AppSourceChecker {

  const val TAG = "AppSourceChecker"

  fun getInstallSource(context: Context): String {
    val packageManager = context.packageManager
    val packageName = context.packageName
    Log.d(TAG, "Package name: $packageName")

    val installerPackageName =
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
          packageManager.getInstallSourceInfo(packageName).installingPackageName
        } else {
          @Suppress("deprecation") packageManager.getInstallerPackageName(packageName)
        }

    Log.d(TAG, "Installer package name: $installerPackageName")

    return when (installerPackageName) {
      "com.android.vending" -> "googleplay"
      "org.fdroid.fdroid" -> "fdroid"
      "com.amazon.venezia" -> "amazon"
      null -> "unknown"
      else -> "unknown($installerPackageName)"
    }
  }
}
