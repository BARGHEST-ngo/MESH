// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
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
import androidx.compose.ui.res.stringResource
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
              var showResetConfirm by remember { mutableStateOf(false) }

              // Show the error overlay if need be
              showErrorDialog?.let {
                ErrorDialog(type = it, action = { viewModel.errorDialog.set(null) })
              }

              // Confirm before resetting as it erases keys/profiles/prefs
              if (showResetConfirm) {
                AlertDialog(
                    onDismissRequest = { showResetConfirm = false },
                    title = { Text(stringResource(R.string.reset_mesh_confirm_title)) },
                    text = { Text(stringResource(R.string.reset_mesh_confirm_message)) },
                    confirmButton = {
                      TextButton(
                          onClick = {
                            showResetConfirm = false
                            viewModel.resetAuth { result ->
                              result
                                  .onSuccess { nav.onNavigateHome() }
                                  .onFailure {
                                    viewModel.errorDialog.set(ErrorDialogType.RESET_FAILED)
                                  }
                            }
                          }) {
                            Text(
                                stringResource(R.string.reset),
                                color = MaterialTheme.colorScheme.error)
                          }
                    },
                    dismissButton = {
                      TextButton(onClick = { showResetConfirm = false }) {
                        Text(stringResource(R.string.cancel))
                      }
                    })
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

                  Setting.Text(
                      R.string.reset_mesh,
                      subtitle = stringResource(R.string.reset_mesh_subtitle),
                      destructive = true,
                      onClick = { showResetConfirm = true })
                  Lists.ItemDivider()
                }
              }
            }
      }

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
