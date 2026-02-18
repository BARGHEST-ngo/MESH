// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalUriHandler
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import com.barghest.mesh.ui.theme.link

@Composable
fun PrimaryActionButton(onClick: () -> Unit, content: @Composable RowScope.() -> Unit) {
  Button(
      onClick = onClick,
      contentPadding = PaddingValues(vertical = 12.dp),
      modifier = Modifier.fillMaxWidth(),
      content = content)
}

@Composable
fun OpenURLButton(title: String, url: String) {
  val handler = LocalUriHandler.current

  TextButton(onClick = { handler.openUri(url) }) {
    Text(
        title,
        style = MaterialTheme.typography.bodyMedium,
        color = MaterialTheme.colorScheme.link,
        textDecoration = TextDecoration.Underline,
    )
  }
}
