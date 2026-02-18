// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalUriHandler
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.tooling.preview.Preview
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.BuildConfig
import com.barghest.mesh.R
import com.barghest.mesh.mdm.AlwaysNeverUserDecides
import com.barghest.mesh.mdm.MDMSettings
import com.barghest.mesh.mdm.ShowHide
import com.barghest.mesh.ui.Links
import com.barghest.mesh.ui.theme.link
import com.barghest.mesh.ui.theme.listItem
import com.barghest.mesh.ui.util.AndroidTVUtil
import com.barghest.mesh.ui.util.AndroidTVUtil.isAndroidTV
import com.barghest.mesh.ui.util.AppVersion
import com.barghest.mesh.ui.util.Lists
import com.barghest.mesh.ui.util.set
import com.barghest.mesh.ui.viewModel.AppViewModel
import com.barghest.mesh.ui.viewModel.SettingsNav
import com.barghest.mesh.ui.viewModel.SettingsViewModel

@Composable
fun SettingsView(
    settingsNav: SettingsNav,
    viewModel: SettingsViewModel = viewModel(),
    appViewModel: AppViewModel = viewModel()
) {
  val handler = LocalUriHandler.current

  val user by viewModel.loggedInUser.collectAsState()
  val isAdmin by viewModel.isAdmin.collectAsState()
  val tailnetLockEnabled by viewModel.tailNetLockEnabled.collectAsState()
  val corpDNSEnabled by viewModel.corpDNSEnabled.collectAsState()
  val isVPNPrepared by appViewModel.vpnPrepared.collectAsState()
  val showTailnetLock by MDMSettings.manageTailnetLock.flow.collectAsState()

  Scaffold(
      topBar = {
        Header(titleRes = R.string.settings_title, onBack = settingsNav.onNavigateBackHome)
      }) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding).verticalScroll(rememberScrollState())) {
          if (isVPNPrepared) {
            UserView(
                profile = user,
                actionState = UserActionState.NAV,
                onClick = settingsNav.onNavigateToUserSwitcher)
          }

          if (isAdmin && !isAndroidTV()) {
            Lists.ItemDivider()
            AdminTextView { handler.openUri(Links.ADMIN_URL) }
          }

          Lists.SectionDivider()
          Setting.Text(
              R.string.dns_settings,
              subtitle =
                  corpDNSEnabled?.let {
                    stringResource(
                        if (it) R.string.using_mesh_dns else R.string.not_using_mesh_dns)
                  },
              onClick = settingsNav.onNavigateToDNSSettings)

          Lists.ItemDivider()
          Setting.Text(
              R.string.split_tunneling,
              subtitle = stringResource(R.string.exclude_certain_apps_from_using_tailscale),
              onClick = settingsNav.onNavigateToSplitTunneling)

          Lists.ItemDivider()
          Setting.Text(
              R.string.awg_settings,
              subtitle = stringResource(R.string.awg_settings_subtitle),
              onClick = settingsNav.onNavigateToAWGSettings)

          if (showTailnetLock.value == ShowHide.Hide) {
              Lists.ItemDivider()
              Setting.Text(
                  R.string.MESH_lock,
                  subtitle =
                      tailnetLockEnabled?.let {
                          stringResource(if (it) R.string.enabled else R.string.disabled)
                      },
                  onClick = settingsNav.onNavigateToTailnetLock
              )
          }
          Lists.SectionDivider()
          Lists.ItemDivider()
          Setting.Text(
              R.string.about_mesh,
              subtitle = "${stringResource(id = R.string.version)} ${AppVersion.Short()}",
              onClick = settingsNav.onNavigateToAbout)
        }
      }
}

object Setting {
  @Composable
  fun Text(
      titleRes: Int = 0,
      title: String? = null,
      subtitle: String? = null,
      destructive: Boolean = false,
      enabled: Boolean = true,
      onClick: (() -> Unit)? = null
  ) {
    var modifier: Modifier = Modifier
    if (enabled) {
      onClick?.let { modifier = modifier.clickable(onClick = it) }
    }
    ListItem(
        modifier = modifier,
        colors = MaterialTheme.colorScheme.listItem,
        headlineContent = {
          Text(
              title ?: stringResource(titleRes),
              style = MaterialTheme.typography.bodyMedium,
              color = if (destructive) MaterialTheme.colorScheme.error else Color.Unspecified)
        },
        supportingContent =
            subtitle?.let {
              {
                Text(
                    it,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
              }
            })
  }

  @Composable
  fun Switch(
      titleRes: Int = 0,
      title: String? = null,
      isOn: Boolean,
      enabled: Boolean = true,
      onToggle: (Boolean) -> Unit = {}
  ) {
    ListItem(
        colors = MaterialTheme.colorScheme.listItem,
        headlineContent = {
          Text(
              title ?: stringResource(titleRes),
              style = MaterialTheme.typography.bodyMedium,
          )
        },
        trailingContent = {
          TintedSwitch(checked = isOn, onCheckedChange = onToggle, enabled = enabled)
        })
  }
}

@Composable
fun AdminTextView(onNavigateToAdminConsole: () -> Unit) {
  Text(
      text = ""
  )
}

@Preview
@Composable
fun SettingsPreview() {
  val vm = SettingsViewModel()
  vm.corpDNSEnabled.set(true)
  vm.tailNetLockEnabled.set(true)
  vm.isAdmin.set(true)
  vm.managedByOrganization.set("")
  SettingsView(SettingsNav({}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}), vm)
}
