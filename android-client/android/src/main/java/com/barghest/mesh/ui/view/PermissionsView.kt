// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.model.Permissions
import com.barghest.mesh.ui.util.friendlyDirName
import com.barghest.mesh.ui.util.itemsWithDividers
import com.barghest.mesh.ui.viewModel.PermissionsViewModel

@Composable
fun PermissionsView(
    backToSettings: BackNavigation,
    navToTaildropDirView: () -> Unit,
    navToNotificationsView: () -> Unit,
    permissionsViewModel: PermissionsViewModel = viewModel()
) {
  val permissions = Permissions.withGrantedStatus

  Scaffold(topBar = { Header(titleRes = R.string.permissions, onBack = backToSettings) }) {
      innerPadding ->
    LazyColumn(modifier = Modifier.padding(innerPadding)) {
      // Existing Android runtime permissions
      itemsWithDividers(permissions) { (permission, granted) ->
        ListItem(
            modifier = Modifier.clickable { navToNotificationsView() },
            leadingContent = {
              Icon(
                  painterResource(R.drawable.baseline_notifications_none_24),
                  tint = MaterialTheme.colorScheme.onSurfaceVariant,
                  modifier = Modifier.size(24.dp),
                  contentDescription =
                      stringResource(if (granted) R.string.ok else R.string.warning))
            },
            headlineContent = {
              Text(stringResource(permission.title), style = MaterialTheme.typography.titleMedium)
            },
            supportingContent = {
              if (granted) Text(stringResource(R.string.on)) else Text(stringResource(R.string.off))
            })
      }

      item {
        ListItem(
            modifier = Modifier.clickable { navToTaildropDirView() },
            leadingContent = {
              Icon(
                  painterResource(R.drawable.baseline_drive_folder_upload_24),
                  tint = MaterialTheme.colorScheme.onSurfaceVariant,
                  modifier = Modifier.size(24.dp),
                  contentDescription = stringResource(R.string.taildrop_dir))
            },
            headlineContent = {
              Text(
                  stringResource(R.string.taildrop_dir_access),
                  style = MaterialTheme.typography.titleMedium)
            },
            supportingContent = {
              val displayPath =
                  permissionsViewModel.currentDir.collectAsState().value?.let {
                    friendlyDirName(it)
                  } ?: "No access"

              Text(displayPath)
            })
      }
    }
  }
}
