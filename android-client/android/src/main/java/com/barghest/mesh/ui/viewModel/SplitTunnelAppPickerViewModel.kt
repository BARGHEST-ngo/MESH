// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.App
import com.barghest.mesh.mdm.MDMSettings
import com.barghest.mesh.mdm.SettingState
import com.barghest.mesh.ui.util.InstalledApp
import com.barghest.mesh.ui.util.InstalledAppsManager
import com.barghest.mesh.ui.util.set
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class SplitTunnelAppPickerViewModel : ViewModel() {
  val installedAppsManager = InstalledAppsManager(packageManager = App.get().packageManager)
  val excludedPackageNames: StateFlow<List<String>> = MutableStateFlow(listOf())
  val installedApps: StateFlow<List<InstalledApp>> = MutableStateFlow(listOf())
  val mdmExcludedPackages: StateFlow<SettingState<String?>> = MDMSettings.excludedPackages.flow
  val mdmIncludedPackages: StateFlow<SettingState<String?>> = MDMSettings.includedPackages.flow

  private var saveJob: Job? = null

  init {
    installedApps.set(installedAppsManager.fetchInstalledApps())
    excludedPackageNames.set(
        App.get()
            .disallowedPackageNames()
            .intersect(installedApps.value.map { it.packageName }.toSet())
            .toList())
  }

  fun exclude(packageName: String) {
    if (excludedPackageNames.value.contains(packageName)) return
    excludedPackageNames.set(excludedPackageNames.value + packageName)
    debounceSave()
  }

  fun unexclude(packageName: String) {
    excludedPackageNames.set(excludedPackageNames.value - packageName)
    debounceSave()
  }

  private fun debounceSave() {
    saveJob?.cancel()
    saveJob =
        viewModelScope.launch {
          delay(500) // Wait to batch multiple rapid updates
          App.get().updateUserDisallowedPackageNames(excludedPackageNames.value)
        }
  }
}
