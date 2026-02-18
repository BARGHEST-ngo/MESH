// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import android.net.Uri
import androidx.activity.result.ActivityResultLauncher
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.Button
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.barghest.mesh.R
import com.barghest.mesh.ui.theme.exitNodeToggleButton
import com.barghest.mesh.ui.util.Lists
import com.barghest.mesh.ui.util.friendlyDirName
import com.barghest.mesh.ui.viewModel.PermissionsViewModel
import com.barghest.mesh.util.TSLog

@Composable
fun TaildropDirView(
    backToPermissionsView: BackNavigation,
    openDirectoryLauncher: ActivityResultLauncher<Uri?>,
    permissionsViewModel: PermissionsViewModel
) {
    Scaffold(
        topBar = {
            Header(titleRes = R.string.taildrop_dir_access, onBack = backToPermissionsView)
        }) { innerPadding ->
        LazyColumn(modifier = Modifier.padding(innerPadding)) {
            item {
                ListItem(
                    headlineContent = {
                        Text(
                            stringResource(R.string.taildrop_dir_access),
                            style = MaterialTheme.typography.titleMedium
                        )
                    },
                    supportingContent = {
                        Text(
                            text = stringResource(R.string.permission_taildrop_dir),
                            style = MaterialTheme.typography.bodyMedium
                        )
                    })
            }

            item("divider0") { Lists.SectionDivider() }

            item {
                val currentDir by permissionsViewModel.currentDir.collectAsState()
                TSLog.d("TaildropDirView", "currentDir in UI: $currentDir")
                val displayPath = currentDir?.let { friendlyDirName(it) } ?: "No access"

                ListItem(
                    headlineContent = {
                        Text(
                            text = stringResource(R.string.dir_access),
                            style = MaterialTheme.typography.titleMedium
                        )
                    }
                )
            }
        }
    }
}
