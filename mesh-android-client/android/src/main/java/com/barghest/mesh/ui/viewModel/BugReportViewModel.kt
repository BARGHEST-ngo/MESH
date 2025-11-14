// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.ui.localapi.Client
import com.barghest.mesh.ui.util.set
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

class BugReportViewModel : ViewModel() {
  val bugReportID: StateFlow<String> = MutableStateFlow("")

  init {
    Client(viewModelScope).bugReportId { result ->
      result
          .onSuccess { bugReportID.set(it.trim()) }
          .onFailure { bugReportID.set("(Error fetching ID)") }
    }
  }
}
