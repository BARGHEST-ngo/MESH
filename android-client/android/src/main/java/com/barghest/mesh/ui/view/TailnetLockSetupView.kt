// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.LocalIndication
import androidx.compose.foundation.clickable
import androidx.compose.foundation.focusable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.text.ClickableText
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalUriHandler
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.tooling.preview.Preview
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.Links
import com.barghest.mesh.ui.theme.defaultTextColor
import com.barghest.mesh.ui.theme.link
import com.barghest.mesh.ui.util.ClipboardValueView
import com.barghest.mesh.ui.util.Lists
import com.barghest.mesh.ui.util.LoadingIndicator
import com.barghest.mesh.ui.util.set
import com.barghest.mesh.ui.viewModel.TailnetLockSetupViewModel
import com.barghest.mesh.ui.viewModel.TailnetLockSetupViewModelFactory

@Composable
fun TailnetLockSetupView(
    backToSettings: BackNavigation,
    model: TailnetLockSetupViewModel = viewModel(factory = TailnetLockSetupViewModelFactory())
) {
  val statusItems by model.statusItems.collectAsState()
  val nodeKey by model.nodeKey.collectAsState()
  val tailnetLockKey by model.tailnetLockKey.collectAsState()
  val tailnetLockTlPubKey = tailnetLockKey.replace("nlpub", "tlpub")

  Scaffold(topBar = { Header(R.string.MESH_lock, onBack = backToSettings) }) { innerPadding ->
    LoadingIndicator.Wrap {
      LazyColumn(modifier = Modifier.padding(innerPadding).fillMaxSize()) {
        item { ExplainerView() }

        items(statusItems) { statusItem ->
          val interactionSource = remember { MutableInteractionSource() }
          ListItem(
              modifier =
                  Modifier.focusable(interactionSource = interactionSource).clickable(
                      interactionSource = interactionSource,
                      indication = LocalIndication.current) {},
              leadingContent = {
                Icon(
                    painter = painterResource(id = statusItem.icon),
                    contentDescription = null,
                    tint = MaterialTheme.colorScheme.onSurfaceVariant)
              },
              headlineContent = { Text(stringResource(statusItem.title)) })
        }

        item {
          // Node key section
          Lists.SectionDivider()
          ClipboardValueView(
              value = nodeKey,
              title = stringResource(R.string.node_key),
              subtitle = stringResource(R.string.node_key_explainer))

          // Tailnet lock key section
          Lists.SectionDivider()
          ClipboardValueView(
              value = tailnetLockTlPubKey,
              title = stringResource(R.string.MESH_lock_key),
              subtitle = stringResource(R.string.MESH_lock_key_explainer))
        }
      }
    }
  }
}

@Composable
private fun ExplainerView() {
  val handler = LocalUriHandler.current

  Lists.MultilineDescription {
    ClickableText(
        explainerText(),
        onClick = { handler.openUri(Links.TAILNET_LOCK_KB_URL) },
        style = MaterialTheme.typography.bodyMedium)
  }
}

@Composable
fun explainerText(): AnnotatedString {
  return buildAnnotatedString {
    withStyle(SpanStyle(color = MaterialTheme.colorScheme.defaultTextColor)) {
      append(stringResource(id = R.string.MESH_lock_explainer))
    }

    pushStringAnnotation(tag = "tailnetLockSupportURL", annotation = Links.TAILNET_LOCK_KB_URL)

    withStyle(
        style =
            SpanStyle(
                color = MaterialTheme.colorScheme.link,
                textDecoration = TextDecoration.Underline)) {
          append(stringResource(id = R.string.learn_more))
        }
    pop()
  }
}

@Composable
@Preview
fun TailnetLockSetupViewPreview() {
  val vm = TailnetLockSetupViewModel()
  vm.nodeKey.set("8BADF00D-EA7-1337-DEAD-BEEF")
  vm.tailnetLockKey.set("C0FFEE-CAFE-50DA")
  TailnetLockSetupView(backToSettings = {}, vm)
}
