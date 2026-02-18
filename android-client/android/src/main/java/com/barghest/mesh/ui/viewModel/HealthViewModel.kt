// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
