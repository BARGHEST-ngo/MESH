// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.mdm

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.RestrictionsManager
import com.barghest.mesh.App
import com.barghest.mesh.util.TSLog

class MDMSettingsChangedReceiver : BroadcastReceiver() {
  override fun onReceive(context: Context?, intent: Intent?) {
    if (intent?.action == Intent.ACTION_APPLICATION_RESTRICTIONS_CHANGED) {
      TSLog.d("syspolicy", "MDM settings changed")
      val restrictionsManager =
          context?.getSystemService(Context.RESTRICTIONS_SERVICE) as RestrictionsManager
      MDMSettings.update(App.get(), restrictionsManager)
    }
  }
}
