// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.util

import androidx.annotation.StringRes
import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import com.barghest.mesh.R

// Convenience wrapper for passing formatted strings to Composables
class ComposableStringFormatter(
    @StringRes val stringRes: Int = R.string.template,
    private vararg val params: Any
) {

  // Convenience constructor for passing a non-formatted string directly
  constructor(string: String) : this(stringRes = R.string.template, string)

  // Returns the fully formatted string
  @Composable fun getString(): String = stringResource(id = stringRes, *params)
}
