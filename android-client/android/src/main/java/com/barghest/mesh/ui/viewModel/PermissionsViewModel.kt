// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import androidx.lifecycle.ViewModel
import com.barghest.mesh.TaildropDirectoryStore
import com.barghest.mesh.util.TSLog
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

class PermissionsViewModel : ViewModel() {
  private val _currentDir =
      MutableStateFlow<String?>(TaildropDirectoryStore.loadSavedDir()?.toString())
  val currentDir: StateFlow<String?> = _currentDir

  fun refreshCurrentDir() {
    val newUri = TaildropDirectoryStore.loadSavedDir()?.toString()
    TSLog.d("PermissionsViewModel", "refreshCurrentDir: $newUri")
    _currentDir.value = newUri
  }
}
