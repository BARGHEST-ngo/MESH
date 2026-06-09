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
import androidx.compose.foundation.rememberScrollState
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.ui.theme.MeshMono
import com.barghest.mesh.ui.theme.MeshRadii
import com.barghest.mesh.ui.theme.accent
import com.barghest.mesh.ui.theme.accentDim
import com.barghest.mesh.ui.theme.meshAmber
import com.barghest.mesh.ui.theme.meshAmberDim
import com.barghest.mesh.ui.theme.meshBg
import com.barghest.mesh.ui.theme.meshBorder
import com.barghest.mesh.ui.theme.meshCard
import com.barghest.mesh.ui.theme.meshGreen
import com.barghest.mesh.ui.theme.meshText
import com.barghest.mesh.ui.theme.meshText2
import android.content.Context
import android.content.Intent
import android.os.Build
import android.provider.Settings
import android.widget.Toast
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import com.barghest.mesh.R

private data class AdbStep(val n: Int, val title: String, val body: String)

/** Opens Developer options (Wireless debugging), falling back to the top-level Settings app. */
private fun openDeveloperSettings(context: Context) {
    if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) {
        Toast.makeText(context, context.getString(R.string.mesh_adb_toast_need_r), Toast.LENGTH_LONG).show()
    }
    val intents = listOf(
        Intent(Settings.ACTION_APPLICATION_DEVELOPMENT_SETTINGS).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK),
        Intent(Settings.ACTION_SETTINGS).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK),
    )
    for (intent in intents) {
        if (context.packageManager.resolveActivity(intent, 0) != null) {
            context.startActivity(intent)
            if (intent.action == Settings.ACTION_SETTINGS) {
                Toast.makeText(context, context.getString(R.string.mesh_adb_toast_navigate), Toast.LENGTH_LONG).show()
            }
            return
        }
    }
    Toast.makeText(context, context.getString(R.string.mesh_adb_toast_unable), Toast.LENGTH_LONG).show()
}

@Composable
fun MeshAdbSetupScreen(
    onBack: () -> Unit,
    onEnabled: () -> Unit,
) {
    val cs = MaterialTheme.colorScheme
    val context = LocalContext.current
    var done by remember { mutableStateOf(setOf<Int>()) }
    val steps = listOf(
        AdbStep(1, stringResource(R.string.mesh_adb_s1_title), stringResource(R.string.mesh_adb_s1_body)),
        AdbStep(2, stringResource(R.string.mesh_adb_s2_title), stringResource(R.string.mesh_adb_s2_body)),
        AdbStep(3, stringResource(R.string.mesh_adb_s3_title), stringResource(R.string.mesh_adb_s3_body)),
        AdbStep(4, stringResource(R.string.mesh_adb_s4_title), stringResource(R.string.mesh_adb_s4_body)),
        AdbStep(5, stringResource(R.string.mesh_adb_s5_title), stringResource(R.string.mesh_adb_s5_body)),
    )

    Column(Modifier.fillMaxSize().background(cs.meshBg)) {
        MeshTopBar(stringResource(R.string.mesh_adb_topbar), onBack = onBack)
        Column(Modifier.weight(1f).verticalScroll(rememberScrollState()).padding(horizontal = 20.dp)) {
            Spacer(Modifier.height(18.dp))
            Text(stringResource(R.string.mesh_adb_title), color = cs.meshText, fontSize = 21.sp, fontWeight = FontWeight.SemiBold, letterSpacing = (-0.3).sp)
            Spacer(Modifier.height(6.dp))
            Text(
                stringResource(R.string.mesh_adb_intro),
                color = cs.meshText2, fontSize = 13.5.sp, lineHeight = 20.sp,
            )
            Spacer(Modifier.height(16.dp))

            Callout("info", cs.accent, cs.accentDim(0.08f), cs.accentDim(0.25f),
                stringResource(R.string.mesh_adb_callout_menus))
            Spacer(Modifier.height(18.dp))

            Column(verticalArrangement = Arrangement.spacedBy(10.dp)) {
                steps.forEach { s ->
                    val ok = s.n in done
                    Row(
                        Modifier
                            .fillMaxWidth()
                            .background(cs.meshCard, RoundedCornerShape(MeshRadii.card))
                            .border(1.dp, if (ok) cs.meshGreen.copy(alpha = 0.4f) else cs.meshBorder, RoundedCornerShape(MeshRadii.card))
                            .clickable { done = if (ok) done - s.n else done + s.n }
                            .padding(14.dp),
                        horizontalArrangement = Arrangement.spacedBy(13.dp),
                        verticalAlignment = Alignment.Top,
                    ) {
                        Box(
                            Modifier.size(28.dp).background(if (ok) cs.meshGreen else cs.accentDim(0.16f), RoundedCornerShape(9.dp)),
                            contentAlignment = Alignment.Center,
                        ) {
                            if (ok) MeshIcon("check", size = 16.dp, color = Color(0xFF0C0D0F), stroke = 2.4f)
                            else Text("${s.n}", color = cs.accent, fontFamily = MeshMono, fontSize = 13.sp, fontWeight = FontWeight.Bold)
                        }
                        Column(Modifier.weight(1f)) {
                            Text(s.title, color = cs.meshText, fontSize = 14.5.sp, fontWeight = FontWeight.SemiBold)
                            Spacer(Modifier.height(2.dp))
                            Text(s.body, color = cs.meshText2, fontSize = 12.5.sp, lineHeight = 18.sp)
                        }
                    }
                }
            }

            Spacer(Modifier.height(16.dp))
            Callout("alert", cs.meshAmber, cs.meshAmberDim, cs.meshAmber.copy(alpha = 0.3f),
                stringResource(R.string.mesh_adb_callout_warn))
            Spacer(Modifier.height(20.dp))
        }

        Column(Modifier.fillMaxWidth().background(cs.meshBg)) {
            Box(Modifier.fillMaxWidth().height(1.dp).background(cs.meshBorder))
            Row(Modifier.padding(horizontal = 20.dp, vertical = 14.dp), horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                MeshButton(stringResource(R.string.mesh_adb_open_settings), { openDeveloperSettings(context) }, modifier = Modifier.weight(1f), variant = MeshButtonVariant.Secondary, height = 56.dp)
                MeshButton(stringResource(R.string.mesh_adb_its_on), onEnabled, modifier = Modifier.weight(1f), height = 56.dp, icon = "check")
            }
        }
    }
}

@Composable
private fun Callout(icon: String, iconColor: Color, bg: Color, border: Color, text: String) {
    val cs = MaterialTheme.colorScheme
    Row(
        Modifier
            .fillMaxWidth()
            .background(bg, RoundedCornerShape(12.dp))
            .border(1.dp, border, RoundedCornerShape(12.dp))
            .padding(13.dp),
        horizontalArrangement = Arrangement.spacedBy(10.dp),
        verticalAlignment = Alignment.Top,
    ) {
        MeshIcon(icon, size = 17.dp, color = iconColor, modifier = Modifier.padding(top = 1.dp))
        Text(text, color = cs.meshText2, fontSize = 12.5.sp, lineHeight = 19.sp)
    }
}
