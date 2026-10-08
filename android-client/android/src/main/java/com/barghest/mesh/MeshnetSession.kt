// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh

import android.os.Build
import android.util.Log
import com.barghest.mesh.ui.model.Ipn
import java.util.concurrent.atomic.AtomicInteger
import com.barghest.mesh.ui.notifier.Notifier
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch
import org.barghest.meshnet.meshnetmobile.Client
import org.barghest.meshnet.meshnetmobile.Listener
import org.barghest.meshnet.meshnetmobile.Meshnetmobile
import org.barghest.meshnet.meshnetmobile.Stream

// MeshnetSession holds the app's single MESHnet client. IPNService owns its lifetime:
// the UI calls prepare() with a scanned address and PIN, then starts the service,
// which calls start(). The address and the PIN are secrets: they are only ever held
// in memory here, never put in an Intent and never logged.
object MeshnetSession {
  private const val TAG = "MeshnetSession"
  private const val HELLO_TOPIC = "hello"

  private class Pairing(val serverAddr: String, val pin: String)

  private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
  @Volatile private var client: Client? = null
  private var pending: Pairing? = null
  // Counts connects, so callbacks from a client that has been replaced are ignored.
  private val generation = AtomicInteger()

  // True from start() until the session ends.
  @Volatile
  var isActive = false
    private set

  // Set by IPNService. Called with each MESHnet state ("Connecting", "Connected",
  // "Disconnected", "Errored") from a Go goroutine, not the main thread.
  @Volatile var onState: ((String) -> Unit)? = null

  // prepare stores the pairing details for the next start(). A PIN works once, so
  // start() consumes them.
  @Synchronized
  fun prepare(serverAddr: String, pin: String) {
    pending = Pairing(serverAddr, pin)
  }

  @Synchronized
  private fun takePending(): Pairing? {
    val p = pending
    pending = null
    return p
  }

  // start connects using the details given to prepare(). It returns false, and does
  // nothing, if there are none. Progress is reported through onState.
  fun start(): Boolean {
    val p = takePending() ?: return false
    isActive = true
    connect(p.serverAddr, p.pin)
    return true
  }

  private fun connect(serverAddr: String, pin: String) {
    val gen = generation.incrementAndGet()
    scope.launch {
      disconnectCurrent()
      val listener =
          object : Listener {
            override fun onStateChange(state: String, err: String) {
              if (gen == generation.get()) onStateChange(gen, state, err)
            }

            // gomobile passes an empty Go slice as null.
            override fun onMessage(topic: String, data: ByteArray?) {
              Log.d(TAG, "message: $topic (${data?.size ?: 0} bytes)")
            }

            override fun onStreamRequest(target: String, stream: Stream) {
              if (gen != generation.get()) {
                stream.close()
                return
              }
              try {
                val port = allowedAdbPort(target)
                if (port == null) {
                  Log.w(TAG, "refusing stream for $target")
                  return
                }
                Log.d(TAG, "stream for $target, proxying to adb")
                Meshnetmobile.proxy(stream, "127.0.0.1:$port")// returns when either side closes
              } catch (e: Exception) {
                Log.e(TAG, "stream for $target: ${e.message}")
              } finally {
                stream.close()
              }
            }
          }
      val c = Meshnetmobile.newClient(serverAddr, pin, Build.MODEL, listener)
      client = c
      try {
        // Blocks until the handshake finishes or fails.
        c.connect()
        Log.d(TAG, "connected, client id ${c.id()}")
        // A first message so the analyst's side can see the channel works.
        c.send(HELLO_TOPIC, "Hello from ${Build.MODEL}".toByteArray())
        Log.d(TAG, "sent $HELLO_TOPIC message")
      } catch (e: Exception) {
        Log.e(TAG, "connect failed: ${e.message}")
      }
    }
  }

  private fun allowedAdbPort(target: String): Int? {
    if (!target.startsWith("adb:")) {
      return null
    }

    val port = target.removePrefix("adb:").toIntOrNull()
    if (port == null || port !in 1..65535) {
      return null
    }

    return port
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
  private fun onStateChange(gen: Int, state: String, err: String) {
    if (err.isEmpty()) Log.d(TAG, "state: $state") else Log.e(TAG, "state: $state: $err")
    if (state != "Connecting" && state != "Connected") isActive = false
    publishState(state)
    onState?.invoke(state)
  }

  // FAKE(meshnet): the screens still read the old tailscale state, so translate the
  // MESHnet state into it.
  private fun publishState(state: String) {
    Notifier.setState(
        when (state) {
          "Connecting" -> Ipn.State.Starting
          "Connected" -> Ipn.State.Running
          else -> Ipn.State.Stopped // Disconnected, Errored
        })
  }
}
