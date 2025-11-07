// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.material3.Switch
import androidx.compose.runtime.Composable

@Composable
fun TintedSwitch(checked: Boolean, onCheckedChange: ((Boolean) -> Unit)?, enabled: Boolean = true) {
  Switch(checked = checked, onCheckedChange = onCheckedChange, enabled = enabled)
}
