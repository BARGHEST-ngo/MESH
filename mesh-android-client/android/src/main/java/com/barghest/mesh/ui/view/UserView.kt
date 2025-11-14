// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.offset
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowRight
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.ListItemColors
import androidx.compose.material3.ListItemDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.barghest.mesh.R
import com.barghest.mesh.ui.model.IpnLocal
import com.barghest.mesh.ui.theme.minTextSize
import com.barghest.mesh.ui.theme.short
import com.barghest.mesh.ui.util.AutoResizingText

// Used to decorate UserViews.
// NONE indicates no decoration
// CURRENT indicates the user is the current user and will be "checked"
// SWITCHING indicates the user is being switched to and will be "loading"
// NAV will show a chevron
enum class UserActionState {
  CURRENT,
  SWITCHING,
  NAV,
  NONE
}

@Composable
fun UserView(
    profile: IpnLocal.LoginProfile?,
    onClick: (() -> Unit)? = null,
    colors: ListItemColors = ListItemDefaults.colors(),
    actionState: UserActionState = UserActionState.NONE,
) {
  Box {
    var modifier: Modifier = Modifier
    onClick?.let { modifier = modifier.clickable { it() } }
    profile?.let {
      ListItem(
          modifier = modifier,
          colors = colors,
          leadingContent = { },
          headlineContent = {
            AutoResizingText(
                text = profile.UserProfile.LoginName,
                style = MaterialTheme.typography.titleMedium.short,
                minFontSize = MaterialTheme.typography.minTextSize,
                overflow = TextOverflow.Ellipsis)
          },
          supportingContent = {
            Column {
              AutoResizingText(
                  text = profile.NetworkProfile?.tailnetNameForDisplay() ?: "",
                  style = MaterialTheme.typography.bodyMedium.short,
                  minFontSize = MaterialTheme.typography.minTextSize,
                  overflow = TextOverflow.Ellipsis)

              profile.customControlServerHostname()?.let {
                AutoResizingText(
                    text = it,
                    style = MaterialTheme.typography.bodyMedium.short,
                    minFontSize = MaterialTheme.typography.minTextSize,
                    overflow = TextOverflow.Ellipsis)
              }
            }
          },
          trailingContent = {
            when (actionState) {
              UserActionState.CURRENT -> CheckedIndicator()
              UserActionState.SWITCHING -> SimpleActivityIndicator(size = 26)
              UserActionState.NAV ->
                  Icon(
                      Icons.AutoMirrored.Filled.KeyboardArrowRight, null, Modifier.offset(x = 6.dp))
              UserActionState.NONE -> Unit
            }
          })
    }
        ?: run {
          ListItem(
              modifier = modifier,
              colors = colors,
              headlineContent = {
                Text(
                    text = stringResource(id = R.string.accounts),
                    style = MaterialTheme.typography.titleMedium)
              })
        }
  }
}
