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
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.App
import com.barghest.mesh.ui.model.Health
import com.barghest.mesh.ui.util.set
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class HealthViewModel : ViewModel() {
  val warnings: StateFlow<List<Health.UnhealthyState>> = MutableStateFlow(listOf())

  init {
    viewModelScope.launch {
      App.get().healthNotifier?.currentWarnings?.collect { set -> warnings.set(set.sorted()) }
    }
  }
}
