// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import androidx.annotation.StringRes
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.R
import com.barghest.mesh.ui.localapi.Client
import com.barghest.mesh.ui.model.Ipn
import com.barghest.mesh.ui.model.Tailcfg
import com.barghest.mesh.ui.notifier.Notifier
import com.barghest.mesh.ui.theme.off
import com.barghest.mesh.ui.theme.success
import com.barghest.mesh.ui.util.set
import com.barghest.mesh.util.TSLog
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.combine
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch

class DNSSettingsViewModelFactory : ViewModelProvider.Factory {
  @Suppress("UNCHECKED_CAST")
  override fun <T : ViewModel> create(modelClass: Class<T>): T {
    return DNSSettingsViewModel() as T
  }
}

class DNSSettingsViewModel : IpnViewModel() {
  val enablementState: StateFlow<DNSEnablementState> =
      MutableStateFlow(DNSEnablementState.NOT_RUNNING)
  val dnsConfig: StateFlow<Tailcfg.DNSConfig?> = MutableStateFlow(null)

  init {
    viewModelScope.launch {
      Notifier.netmap
          .combine(Notifier.prefs) { netmap, prefs -> Pair(netmap, prefs) }
          .stateIn(viewModelScope)
          .collect { (netmap, prefs) ->
            TSLog.d("DNSSettingsViewModel", "prefs: CorpDNS=" + prefs?.CorpDNS.toString())
            prefs?.let {
              if (it.CorpDNS) {
                enablementState.set(DNSEnablementState.ENABLED)
              } else {
                enablementState.set(DNSEnablementState.DISABLED)
              }
            } ?: run { enablementState.set(DNSEnablementState.NOT_RUNNING) }
            netmap?.let { dnsConfig.set(netmap.DNS) }
          }
    }
  }

  fun toggleCorpDNS(callback: (Result<Ipn.Prefs>) -> Unit) {
    val prefs =
        Notifier.prefs.value
            ?: run {
              callback(Result.failure(Exception("no prefs")))
              return@toggleCorpDNS
            }

    val prefsOut = Ipn.MaskedPrefs()
    prefsOut.CorpDNS = !prefs.CorpDNS
    Client(viewModelScope).editPrefs(prefsOut, callback)
  }
}

enum class DNSEnablementState(
    @StringRes val title: Int,
    @StringRes val caption: Int,
    val symbolDrawable: Int,
    val tint: @Composable () -> Color
) {
  NOT_RUNNING(
      R.string.not_running,
      R.string.mesh_is_not_running_this_device_is_using_the_system_dns_resolver,
      R.drawable.xmark_circle,
      { MaterialTheme.colorScheme.off }),
  ENABLED(
      R.string.using_mesh_dns,
      R.string.this_device_is_using_tailscale_to_resolve_dns_names,
      R.drawable.check_circle,
      { MaterialTheme.colorScheme.success }),
  DISABLED(
      R.string.not_using_mesh_dns,
      R.string.this_device_is_using_the_system_dns_resolver,
      R.drawable.xmark_circle,
      { MaterialTheme.colorScheme.error })
}
