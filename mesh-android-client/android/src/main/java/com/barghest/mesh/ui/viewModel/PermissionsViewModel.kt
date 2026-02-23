// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
