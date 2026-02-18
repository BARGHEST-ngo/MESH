// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.model.Health
import com.barghest.mesh.ui.theme.success
import com.barghest.mesh.ui.viewModel.HealthViewModel

@Composable
fun HealthView(backToSettings: BackNavigation, model: HealthViewModel = viewModel()) {
  val warnings by model.warnings.collectAsState()

  Scaffold(topBar = { Header(titleRes = R.string.health_warnings, onBack = backToSettings) }) {
      innerPadding ->
    LazyColumn(modifier = Modifier.padding(innerPadding)) {
      if (warnings.isEmpty()) {
        item("allGood") {
          Column(
              horizontalAlignment = Alignment.CenterHorizontally,
              verticalArrangement = Arrangement.spacedBy(16.dp, alignment = Alignment.Top),
              modifier = Modifier.fillMaxWidth().padding(vertical = 16.dp)) {
                Icon(
                    painter = painterResource(id = R.drawable.check_circle),
                    modifier = Modifier.size(48.dp),
                    contentDescription = "A green checkmark",
                    tint = MaterialTheme.colorScheme.success)
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement =
                        Arrangement.spacedBy(2.dp, alignment = Alignment.CenterVertically),
                    modifier = Modifier.fillMaxWidth()) {
                      Text(
                          text = stringResource(R.string.no_issues_found),
                          fontSize = MaterialTheme.typography.titleMedium.fontSize,
                          fontWeight = MaterialTheme.typography.titleMedium.fontWeight)
                      Text(
                          text = stringResource(R.string.tailscale_is_operating_normally),
                          color = MaterialTheme.colorScheme.secondary)
                    }
              }
        }
      }

      items(warnings) { HealthWarningView(it) }
    }
  }
}

@Composable
fun HealthWarningView(warning: Health.UnhealthyState) {
  Box(modifier = Modifier.background(color = MaterialTheme.colorScheme.surfaceContainerLow)) {
    Box(
        modifier =
            Modifier.padding(start = 16.dp, end = 16.dp, top = 8.dp, bottom = 8.dp)
                .clip(shape = RoundedCornerShape(10.dp, 10.dp, 10.dp, 10.dp))
                .fillMaxWidth()) {
          ListItem(
              colors = warning.Severity.listItemColors(),
              headlineContent = {
                if (warning.Title.isNotEmpty()) {
                  Text(
                      warning.Title,
                      style = MaterialTheme.typography.titleMedium,
                  )
                }
              },
              supportingContent = {
                Text(warning.Text, style = MaterialTheme.typography.bodyMedium)
              })
        }
  }
}
