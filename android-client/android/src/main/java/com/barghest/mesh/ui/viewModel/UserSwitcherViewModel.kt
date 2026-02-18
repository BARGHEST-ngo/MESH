// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import com.barghest.mesh.ui.view.ErrorDialogType
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

class UserSwitcherViewModel : IpnViewModel() {

  // Set to a non-null value to show the appropriate error dialog
  val errorDialog: StateFlow<ErrorDialogType?> = MutableStateFlow(null)

  // True if we should render the kebab menu
  val showHeaderMenu: StateFlow<Boolean> = MutableStateFlow(false)
}
