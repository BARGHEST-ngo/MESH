// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

import androidx.activity.compose.BackHandler
import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.togetherWith
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.navigationBarsPadding
import androidx.compose.foundation.layout.statusBarsPadding
import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import com.barghest.mesh.ui.theme.meshBg

private enum class MeshScreen { Home, Learn, Adb, Settings }

/** Live backend wiring for the redesigned MESH client. */
data class MeshLiveEnv(
    val connected: Boolean,
    val analyst: Analyst?,        // built from the netmap peer; null => placeholder
    val controlUrl: String?,
    val exitNodeName: String?,    // active exit node display name, or null = none
    val onStartSession: () -> Unit,
    val onEndSession: () -> Unit,
    val onOpenExitNodes: () -> Unit,
    val onOpenControlServer: () -> Unit,
    val onOpenDns: () -> Unit,
    val onOpenSplitTunnel: () -> Unit,
    val onOpenObfuscation: () -> Unit,
    val onForgetEnrollment: () -> Unit,
)

/** Host for the redesigned MESH client; [onExit] fires when backing out of the root. */
@Composable
fun MeshFlow(env: MeshLiveEnv, onExit: () -> Unit = {}) {
    val cs = MaterialTheme.colorScheme
    var screen by remember { mutableStateOf(MeshScreen.Home) }
    var adbReady by remember { mutableStateOf(false) }
    var endConfirm by remember { mutableStateOf(false) }

    val connected = env.connected
    val analyst = env.analyst ?: MeshDefaults.analyst

    fun go(next: MeshScreen) {
        screen = next
    }

    BackHandler(enabled = true) {
        when {
            endConfirm -> endConfirm = false           // dismiss the end-session sheet first
            screen == MeshScreen.Home -> onExit()
            else -> go(MeshScreen.Home)
        }
    }

    val homeActions = remember(env) {
        MeshHomeActions(
            onStart = { env.onStartSession() },
            onEnd = { endConfirm = true },
            onOpenSettings = { go(MeshScreen.Settings) },
            onOpenLearn = { go(MeshScreen.Learn) },
            onOpenExitNodes = { env.onOpenExitNodes() },
            onOpenAdb = { go(MeshScreen.Adb) },
        )
    }
    val settingsLive = remember(env) {
        MeshSettingsLive(
            controlUrl = env.controlUrl,
            onOpenControlServer = env.onOpenControlServer,
            onOpenExitNodes = env.onOpenExitNodes,
            onOpenDns = env.onOpenDns,
            onOpenSplitTunnel = env.onOpenSplitTunnel,
            onOpenObfuscation = env.onOpenObfuscation,
            onForgetEnrollment = env.onForgetEnrollment,
        )
    }

    Column(Modifier.fillMaxSize().background(cs.meshBg)) {
        Box(Modifier.weight(1f).fillMaxWidth().statusBarsPadding().navigationBarsPadding()) {
            AnimatedContent(
                targetState = screen,
                transitionSpec = {
                    // Crossfade — a slide re-records both screen trees each frame and janked.
                    fadeIn(tween(140)) togetherWith fadeOut(tween(140))
                },
                label = "meshScreen",
            ) { s ->
                when (s) {
                    MeshScreen.Home -> MeshHomeScreen(connected, homeActions, analyst = analyst, exitNodeName = env.exitNodeName, adbReady = adbReady)
                    MeshScreen.Learn -> MeshLearnScreen(onBack = { go(MeshScreen.Home) }, onStart = { env.onStartSession() })
                    MeshScreen.Adb -> MeshAdbSetupScreen(
                        onBack = { go(MeshScreen.Home) },
                        onEnabled = { adbReady = true; go(MeshScreen.Home) },
                    )
                    MeshScreen.Settings -> MeshSettingsScreen(
                        onBack = { go(MeshScreen.Home) },
                        live = settingsLive,
                    )
                }
            }
        }
    }

    if (endConfirm) MeshConfirmSheet(onCancel = { endConfirm = false }, onConfirm = { endConfirm = false; env.onEndSession() })
}
