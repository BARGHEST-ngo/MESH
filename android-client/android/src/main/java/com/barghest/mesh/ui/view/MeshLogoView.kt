// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import com.barghest.mesh.R

@Composable
fun TailscaleLogoView(
    animated: Boolean = false,
    usesOnBackgroundColors: Boolean = false,
    modifier: Modifier
) {
      Image(
          painter = painterResource(id = R.drawable.ic_launcher_foreground),
          contentDescription = "MESH Logo",
          modifier = modifier.fillMaxWidth().fillMaxHeight().fillMaxSize())
  }