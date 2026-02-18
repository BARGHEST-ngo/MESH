// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.
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
