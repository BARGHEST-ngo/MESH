// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh

import android.os.Build
import android.util.Log
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch
import org.barghest.meshnet.meshnetmobile.Client
import org.barghest.meshnet.meshnetmobile.Listener
import org.barghest.meshnet.meshnetmobile.Meshnetmobile

object MeshnetSession : Listener {
  private const val TAG = "MeshnetSession"

  private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
  @Volatile private var client: Client? = null

  fun connect(serverAddr: String, pin: String) {
    scope.launch {
      disconnectCurrent()
      val c = Meshnetmobile.newClient(serverAddr, pin, Build.MODEL, this@MeshnetSession)
      client = c
      try {
        // Blocks until the handshake finishes or fails.
        c.connect()
        Log.d(TAG, "connected, client id ${c.id()}")
      } catch (e: Exception) {
        Log.e(TAG, "connect failed: ${e.message}")
      }
    }
  }

  fun disconnect() {
    scope.launch { disconnectCurrent() }
  }

  private fun disconnectCurrent() {
    val c = client ?: return
    client = null
    try {
      c.disconnect()
    } catch (e: Exception) {
      Log.e(TAG, "disconnect failed: ${e.message}")
    }
  }

  // Called from a Go goroutine, not the main thread.
  override fun onStateChange(state: String, err: String) {
    if (err.isEmpty()) Log.d(TAG, "state: $state") else Log.e(TAG, "state: $state: $err")
  }

  // Called from a Go goroutine, not the main thread.
  override fun onMessage(topic: String, data: ByteArray?) {
    Log.d(TAG, "message: $topic (${data?.size ?: 0} bytes)")
  }
}
