// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalUriHandler
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.R
import com.barghest.mesh.ui.Links
import com.barghest.mesh.ui.theme.MeshRadii
import com.barghest.mesh.ui.util.AppVersion
import com.barghest.mesh.ui.theme.accent
import com.barghest.mesh.ui.theme.accentDim
import com.barghest.mesh.ui.theme.meshBg
import com.barghest.mesh.ui.theme.meshBgRaised
import com.barghest.mesh.ui.theme.meshBorder
import com.barghest.mesh.ui.theme.meshBorderHi
import com.barghest.mesh.ui.theme.meshCard
import com.barghest.mesh.ui.theme.meshGreen
import com.barghest.mesh.ui.theme.meshGreenDim
import com.barghest.mesh.ui.theme.meshMuted
import com.barghest.mesh.ui.theme.meshRed
import com.barghest.mesh.ui.theme.meshRedDim
import com.barghest.mesh.ui.theme.meshText
import com.barghest.mesh.ui.theme.meshText2

@Composable
private fun SettingRow(
    icon: String,
    title: String,
    sub: String? = null,
    danger: Boolean = false,
    onClick: (() -> Unit)? = null,
    control: @Composable () -> Unit = {
        MeshIcon("chevR", size = 17.dp, color = if (danger) MaterialTheme.colorScheme.meshRed else MaterialTheme.colorScheme.meshMuted)
    },
) {
    val cs = MaterialTheme.colorScheme
    MeshRow(
        modifier = (if (onClick != null) Modifier.clickable(onClick = onClick) else Modifier)
            .padding(horizontal = 14.dp, vertical = 13.dp),
        icon = icon,
        title = title,
        sub = sub,
        iconTint = if (danger) cs.meshRed else cs.accent,
        iconBg = if (danger) cs.meshRedDim else cs.accentDim(0.13f),
        titleColor = if (danger) cs.meshRed else cs.meshText,
        spacing = 13.dp,
        trailing = control,
    )
}

@Composable
private fun Group(label: String, content: @Composable () -> Unit) {
    val cs = MaterialTheme.colorScheme
    Column(Modifier.padding(bottom = 18.dp)) {
        Eyebrow(label, modifier = Modifier.padding(start = 4.dp, bottom = 8.dp))
        Column(Modifier.fillMaxWidth().background(cs.meshCard, RoundedCornerShape(MeshRadii.card)).border(1.dp, cs.meshBorder, RoundedCornerShape(MeshRadii.card))) {
            content()
        }
    }
}

/** Live Settings wiring. */
data class MeshSettingsLive(
    val controlUrl: String?,
    val onOpenControlServer: () -> Unit,
    val onOpenExitNodes: () -> Unit,
    val onOpenDns: () -> Unit,
    val onOpenSplitTunnel: () -> Unit,
    val onOpenObfuscation: () -> Unit,
    val onForgetEnrollment: () -> Unit,
)

@Composable
fun MeshSettingsScreen(onBack: () -> Unit, live: MeshSettingsLive) {
    val cs = MaterialTheme.colorScheme
    val uriHandler = LocalUriHandler.current

    Column(Modifier.fillMaxSize().background(cs.meshBg)) {
        MeshTopBar(stringResource(R.string.mesh_topbar_settings), onBack = onBack)
        Column(Modifier.weight(1f).verticalScroll(rememberScrollState()).padding(horizontal = 18.dp, vertical = 16.dp)) {
            Group(stringResource(R.string.mesh_set_group_connection)) {
                SettingRow("globe", stringResource(R.string.mesh_set_control_server), live.controlUrl?.takeIf { it.isNotEmpty() } ?: stringResource(R.string.mesh_set_not_configured), onClick = live.onOpenControlServer)
                MeshDivider()
                SettingRow("globe", stringResource(R.string.mesh_set_exit_node), stringResource(R.string.mesh_set_exit_node_sub), onClick = live.onOpenExitNodes)
            }
            Group(stringResource(R.string.mesh_set_group_network)) {
                SettingRow("globe", stringResource(R.string.mesh_set_magic_dns), stringResource(R.string.mesh_set_magic_dns_sub), onClick = live.onOpenDns)
                MeshDivider()
                SettingRow("grid", stringResource(R.string.mesh_set_split), stringResource(R.string.mesh_set_split_sub), onClick = live.onOpenSplitTunnel)
            }
            Group(stringResource(R.string.mesh_set_group_privacy)) {
                SettingRow("eyeOff", stringResource(R.string.mesh_set_obfuscation), stringResource(R.string.mesh_set_obfuscation_sub), onClick = live.onOpenObfuscation)
            }
            Group(stringResource(R.string.mesh_set_group_about)) {
                SettingRow("info", stringResource(R.string.mesh_set_version, AppVersion.Short()), stringResource(R.string.mesh_set_source), onClick = { uriHandler.openUri(Links.SOURCE_REPO_URL) })
                MeshDivider()
                SettingRow("file", stringResource(R.string.mesh_set_license), stringResource(R.string.mesh_set_license_sub), onClick = { uriHandler.openUri(Links.LICENSE_URL) })
            }
            Group(stringResource(R.string.mesh_set_group_danger)) {
                SettingRow("x", stringResource(R.string.mesh_set_forget), stringResource(R.string.mesh_set_forget_sub), danger = true, onClick = live.onForgetEnrollment)
            }
            Text(stringResource(R.string.mesh_set_footer), color = cs.meshMuted, fontSize = 11.5.sp, fontWeight = FontWeight.Medium, modifier = Modifier.fillMaxWidth().padding(top = 8.dp), textAlign = TextAlign.Center)
        }
    }
}

private data class LearnStep(val icon: String, val titleRes: Int, val bodyRes: Int)

@Composable
fun MeshLearnScreen(onBack: () -> Unit, onStart: () -> Unit) {
    val cs = MaterialTheme.colorScheme
    val steps = remember {
        listOf(
            LearnStep("qr", R.string.mesh_learn_s1_title, R.string.mesh_learn_s1_body),
            LearnStep("lock", R.string.mesh_learn_s2_title, R.string.mesh_learn_s2_body),
            LearnStep("eye", R.string.mesh_learn_s3_title, R.string.mesh_learn_s3_body),
            LearnStep("power", R.string.mesh_learn_s4_title, R.string.mesh_learn_s4_body),
        )
    }
    Column(Modifier.fillMaxSize().background(cs.meshBg)) {
        MeshTopBar(stringResource(R.string.mesh_learn_topbar), onBack = onBack)
        Column(Modifier.weight(1f).verticalScroll(rememberScrollState()).padding(horizontal = 20.dp)) {
            Spacer(Modifier.height(18.dp))
            Text(
                stringResource(R.string.mesh_learn_intro),
                color = cs.meshText2, fontSize = 14.5.sp, lineHeight = 22.sp,
            )
            Spacer(Modifier.height(22.dp))
            steps.forEachIndexed { i, s ->
                Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(14.dp)) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Box(
                            Modifier.size(42.dp).background(cs.accentDim(0.14f), RoundedCornerShape(13.dp)).border(1.dp, cs.accentDim(0.3f), RoundedCornerShape(13.dp)),
                            contentAlignment = Alignment.Center,
                        ) { MeshIcon(s.icon, size = 20.dp, color = cs.accent) }
                        if (i < steps.lastIndex) Box(Modifier.padding(vertical = 6.dp).width(1.5.dp).height(22.dp).background(cs.meshBorder))
                    }
                    Column(Modifier.padding(bottom = if (i < steps.lastIndex) 0.dp else 22.dp)) {
                        Text(stringResource(s.titleRes), color = cs.meshText, fontSize = 15.5.sp, fontWeight = FontWeight.SemiBold)
                        Spacer(Modifier.height(4.dp))
                        Text(stringResource(s.bodyRes), color = cs.meshText2, fontSize = 13.sp, lineHeight = 19.sp)
                    }
                }
            }
        }
        MeshBottomBar {
            MeshButton(stringResource(R.string.mesh_start_cta), onStart, full = true, height = 56.dp, icon = "shield")
        }
    }
}

@Composable
fun MeshConfirmSheet(onCancel: () -> Unit, onConfirm: () -> Unit) {
    val cs = MaterialTheme.colorScheme
    Box(Modifier.fillMaxSize()) {
        Box(Modifier.fillMaxSize().background(Color(0x99000000)).clickable(onClick = onCancel))
        Column(
            Modifier
                .align(Alignment.BottomCenter)
                .fillMaxWidth()
                .background(cs.meshBgRaised, RoundedCornerShape(topStart = 22.dp, topEnd = 22.dp))
                .border(1.dp, cs.meshBorder, RoundedCornerShape(topStart = 22.dp, topEnd = 22.dp))
                .padding(start = 20.dp, end = 20.dp, top = 10.dp, bottom = 24.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Box(Modifier.padding(bottom = 18.dp).width(38.dp).height(4.dp).background(cs.meshBorderHi, RoundedCornerShape(2.dp)))
            Box(Modifier.size(54.dp).background(cs.meshRedDim, RoundedCornerShape(16.dp)), contentAlignment = Alignment.Center) {
                MeshIcon("power", size = 26.dp, color = cs.meshRed)
            }
            Spacer(Modifier.height(16.dp))
            Text(stringResource(R.string.mesh_end_confirm_title), color = cs.meshText, fontSize = 19.sp, fontWeight = FontWeight.Bold, textAlign = TextAlign.Center)
            Spacer(Modifier.height(8.dp))
            Text(stringResource(R.string.mesh_end_confirm_body), color = cs.meshText2, fontSize = 13.5.sp, lineHeight = 20.sp, textAlign = TextAlign.Center)
            Spacer(Modifier.height(22.dp))
            MeshButton(stringResource(R.string.mesh_end_cta), onConfirm, full = true, height = 56.dp, icon = "power", variant = MeshButtonVariant.Danger)
            Spacer(Modifier.height(10.dp))
            MeshButton(stringResource(R.string.mesh_keep_session), onCancel, full = true, height = 52.dp, variant = MeshButtonVariant.Ghost)
        }
    }
}
