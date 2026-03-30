// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.viewModel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.ui.localapi.Client
import com.barghest.mesh.ui.model.Ipn
import com.barghest.mesh.ui.model.StableNodeID
import com.barghest.mesh.ui.notifier.Notifier
import com.barghest.mesh.ui.util.LoadingIndicator
import com.barghest.mesh.ui.util.set
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.combine
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch

data class ExitNodePickerNav(
    val onNavigateBackHome: () -> Unit,
    val onNavigateBackToExitNodes: () -> Unit,
    val onNavigateToRunAsExitNode: () -> Unit,
)

class ExitNodePickerViewModelFactory(private val nav: ExitNodePickerNav) :
    ViewModelProvider.Factory {
  @Suppress("UNCHECKED_CAST")
  override fun <T : ViewModel> create(modelClass: Class<T>): T {
    return ExitNodePickerViewModel(nav) as T
  }
}

class ExitNodePickerViewModel(private val nav: ExitNodePickerNav) : IpnViewModel() {
  data class ExitNode(
      val id: StableNodeID? = null,
      val label: String,
      val online: StateFlow<Boolean>,
      val selected: Boolean,
      val city: String = ""
  )

  val tailnetExitNodes: StateFlow<List<ExitNode>> = MutableStateFlow(emptyList())
  val anyActive: StateFlow<Boolean> = MutableStateFlow(false)

  init {
    viewModelScope.launch {
      Notifier.netmap
          .combine(Notifier.prefs) { netmap, prefs -> Pair(netmap, prefs) }
          .stateIn(viewModelScope)
          .collect { (netmap, prefs) ->
            val exitNodeId = prefs?.activeExitNodeID ?: prefs?.selectedExitNodeID
            netmap?.Peers?.let { peers ->
              val allNodes =
                  peers
                      .filter { it.isExitNode }
                      .map {
                        ExitNode(
                            id = it.StableID,
                            label = it.displayName,
                            online = MutableStateFlow(it.Online ?: false),
                            selected = it.StableID == exitNodeId,
                            city = it.Hostinfo.Location?.City ?: "",
                        )
                      }

              tailnetExitNodes.set(allNodes.sortedWith { a, b -> a.label.compareTo(b.label) })
              anyActive.set(allNodes.any { it.selected })
            }
          }
    }
  }

  fun setExitNode(node: ExitNode) {
    LoadingIndicator.start()
    val prefsOut = Ipn.MaskedPrefs()
    prefsOut.ExitNodeID = node.id

    Client(viewModelScope).editPrefs(prefsOut) {
      nav.onNavigateBackHome()
      LoadingIndicator.stop()
    }
  }

  fun toggleAllowLANAccess(callback: (Result<Ipn.Prefs>) -> Unit) {
    val prefs =
        Notifier.prefs.value
            ?: run {
              callback(Result.failure(Exception("no prefs")))
              return@toggleAllowLANAccess
            }

    val prefsOut = Ipn.MaskedPrefs()
    prefsOut.ExitNodeAllowLANAccess = !prefs.ExitNodeAllowLANAccess
    Client(viewModelScope).editPrefs(prefsOut, callback)
  }
}

val List<ExitNodePickerViewModel.ExitNode>.selected
  get() = this.any { it.selected }
