// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.util

import android.Manifest
import android.content.pm.ApplicationInfo
import android.content.pm.PackageManager

data class InstalledApp(val name: String, val packageName: String)

class InstalledAppsManager(
    val packageManager: PackageManager,
) {
  fun fetchInstalledApps(): List<InstalledApp> {
    return packageManager
        .getInstalledApplications(PackageManager.GET_META_DATA)
        .filter(appIsIncluded)
        .map {
          InstalledApp(
              name = it.loadLabel(packageManager).toString(),
              packageName = it.packageName,
          )
        }
        .sortedBy { it.name }
  }

  private val appIsIncluded: (ApplicationInfo) -> Boolean = { app ->
    app.packageName != "com.barghest.com.barghest.com.barghest.com.barghest.mesh" &&
        // Only show apps that can access the Internet
        packageManager.checkPermission(Manifest.permission.INTERNET, app.packageName) ==
            PackageManager.PERMISSION_GRANTED
  }
}
