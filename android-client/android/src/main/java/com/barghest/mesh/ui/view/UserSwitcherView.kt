// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import android.content.Intent
import android.net.Uri

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.text.ClickableText
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.util.Lists

import com.barghest.mesh.ui.util.set
import com.barghest.mesh.ui.viewModel.UserSwitcherViewModel

data class UserSwitcherNav(
    val backToSettings: BackNavigation,
    val onNavigateHome: () -> Unit
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun UserSwitcherView(nav: UserSwitcherNav, viewModel: UserSwitcherViewModel = viewModel()) {
  val currentUser by viewModel.loggedInUser.collectAsState()
  var showDeleteDialog by remember { mutableStateOf(false) }
  val context = LocalContext.current
  val netmapState by viewModel.netmap.collectAsState()
  val capabilityIsOwner = ""
  val isOwner = netmapState?.hasCap(capabilityIsOwner) == true

  Scaffold(
      topBar = {
        Header(
            R.string.accounts,
            onBack = nav.backToSettings)
      }) { innerPadding ->
        Column(
            modifier = Modifier.padding(innerPadding).fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(8.dp)) {
              val showErrorDialog by viewModel.errorDialog.collectAsState()

              // Show the error overlay if need be
              showErrorDialog?.let {
                ErrorDialog(type = it, action = { viewModel.errorDialog.set(null) })
              }

              LazyColumn {
                item {
                  // Log out button (only if user is logged in)
                  if (currentUser != null) {
                    Setting.Text(
                        R.string.log_out,
                        destructive = true,
                        onClick = {
                          viewModel.logout {
                            it.onSuccess { nav.onNavigateHome() }
                                .onFailure {
                                  viewModel.errorDialog.set(ErrorDialogType.LOGOUT_FAILED)
                                }
                          }
                        })
                    Lists.ItemDivider()
                  }
                }
              }
            }
      }

  if (showDeleteDialog) {
    AlertDialog(
        onDismissRequest = { showDeleteDialog = false },
        title = { Text(text = stringResource(R.string.delete_MESH)) },
        text = {
          if (isOwner) {
            OwnerDeleteDialogText {
              val uri = Uri.parse("https://login.tailscale.com/admin/settings/general")
              context.startActivity(Intent(Intent.ACTION_VIEW, uri))
            }
          } else {
            Text(stringResource(R.string.request_deletion_nonowner))
          }
        },
        confirmButton = {
          TextButton(
              onClick = {
                val intent =
                    Intent(Intent.ACTION_VIEW, Uri.parse("https://tailscale.com/contact/support"))
                context.startActivity(intent)
                showDeleteDialog = false
              }) {
                Text(text = stringResource(R.string.contact_support))
              }
        },
        dismissButton = {
          TextButton(onClick = { showDeleteDialog = false }) {
            Text(text = stringResource(R.string.cancel))
          }
        })
  }
}



@Composable
fun OwnerDeleteDialogText(onSettingsClick: () -> Unit) {
  val part1 = stringResource(R.string.request_deletion_owner_part1)
  val part2a = stringResource(R.string.request_deletion_owner_part2a)
  val part2b = stringResource(R.string.request_deletion_owner_part2b)

  val annotatedText = buildAnnotatedString {
    append(part1 + " ")

    pushStringAnnotation(
        tag = "settings", annotation = "https://login.tailscale.com/admin/settings/general")
    withStyle(style = SpanStyle(color = MaterialTheme.colorScheme.primary)) {
      append("Settings > General")
    }
    pop()

    append(" $part2a\n\n") // newline after "Delete tailnet."
    append(part2b)
  }

  val context = LocalContext.current
  ClickableText(
      text = annotatedText,
      style = MaterialTheme.typography.bodyMedium.copy(color = MaterialTheme.colorScheme.onSurface),
      onClick = { offset ->
        annotatedText
            .getStringAnnotations(tag = "settings", start = offset, end = offset)
            .firstOrNull()
            ?.let { annotation ->
              val intent = Intent(Intent.ACTION_VIEW, Uri.parse(annotation.item))
              context.startActivity(intent)
            }
      })
}



@Composable
@Preview
fun UserSwitcherViewPreview() {
  val vm = UserSwitcherViewModel()
  val nav =
      UserSwitcherNav(
          backToSettings = {},
          onNavigateHome = {})
  UserSwitcherView(nav, vm)
}
