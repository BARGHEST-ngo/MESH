// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.util.Lists
import com.barghest.mesh.ui.util.LoadingIndicator
import com.barghest.mesh.ui.util.flag
import com.barghest.mesh.ui.util.itemsWithDividers
import com.barghest.mesh.ui.viewModel.ExitNodePickerNav
import com.barghest.mesh.ui.viewModel.ExitNodePickerViewModel
import com.barghest.mesh.ui.viewModel.ExitNodePickerViewModelFactory

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MullvadExitNodePicker(
    countryCode: String,
    nav: ExitNodePickerNav,
    model: ExitNodePickerViewModel = viewModel(factory = ExitNodePickerViewModelFactory(nav))
) {
    val mullvadExitNodes by model.mullvadExitNodesByCountryCode.collectAsState()
    val bestAvailableByCountry by model.mullvadBestAvailableByCountry.collectAsState()

    mullvadExitNodes[countryCode]?.toList()?.let { nodes ->
        val any = nodes.first()

        LoadingIndicator.Wrap {
            Scaffold(
                topBar = {
                    Header(
                        title = { Text("${countryCode.flag()} ${any.country}") },
                        onBack = nav.onNavigateBackToMullvad)
                }) { innerPadding ->
                LazyColumn(modifier = Modifier.padding(innerPadding)) {
                    if (nodes.size > 1) {
                        val bestAvailableNode = bestAvailableByCountry[countryCode]!!
                        item {
                            ExitNodeItem(
                                model,
                                ExitNodePickerViewModel.ExitNode(
                                    id = bestAvailableNode.id,
                                    label = stringResource(R.string.best_available),
                                    online = bestAvailableNode.online,
                                    selected = false,
                                ))
                            Lists.SectionDivider()
                        }
                    }

                    itemsWithDividers(nodes) { node -> ExitNodeItem(model, node) }
                }
            }
        }
    }
}