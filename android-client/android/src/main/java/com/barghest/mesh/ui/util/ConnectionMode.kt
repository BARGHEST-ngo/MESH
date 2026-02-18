// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.util

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import com.barghest.mesh.R
import com.barghest.mesh.ui.theme.on

sealed class ConnectionMode {
  class NotConnected : ConnectionMode()

  class Derp(val relayName: String) : ConnectionMode()

  class Direct : ConnectionMode()

  @Composable
  fun titleString(): String {
    return when (this) {
      is NotConnected -> stringResource(id = R.string.not_connected)
      is Derp -> stringResource(R.string.relayed_connection, relayName)
      is Direct -> stringResource(R.string.direct_connection)
    }
  }

  fun contentKey(): String {
    return when (this) {
      is NotConnected -> "NotConnected"
      is Derp -> "Derp($relayName)"
      is Direct -> "Direct"
    }
  }

  fun iconDrawable(): Int {
    return when (this) {
      is NotConnected -> R.drawable.xmark_circle
      is Derp -> R.drawable.link_off
      is Direct -> R.drawable.link
    }
  }

  @Composable
  fun color(): Color {
    return when (this) {
      is NotConnected -> MaterialTheme.colorScheme.onPrimary
      is Derp -> MaterialTheme.colorScheme.error
      is Direct -> MaterialTheme.colorScheme.on
    }
  }
}
