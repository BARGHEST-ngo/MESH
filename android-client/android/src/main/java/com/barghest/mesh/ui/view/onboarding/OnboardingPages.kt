// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view.onboarding

import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Check
import androidx.compose.material.icons.outlined.Lock
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.Font
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.barghest.mesh.R

internal val JetBrainsMonoFamily = FontFamily(
    Font(R.font.jetbrains_mono_regular, FontWeight.Normal),
    Font(R.font.jetbrains_mono_bold, FontWeight.Bold)
)

@Composable
internal fun WelcomePage() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color.Black)
            .padding(horizontal = 24.dp, vertical = 32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Top
    ) {
        Spacer(Modifier.weight(1f))
        Image(
            painter = painterResource(id = R.drawable.mesh_playstore),
            contentDescription = "MESH Logo",
            modifier = Modifier.size(160.dp)
        )
        Spacer(Modifier.height(24.dp))
        Spacer(Modifier.height(12.dp))
        Text(
            text = "MESH is a secure, peer-to-peer networking tool that enables " +
                    "consensual Android device analysis. It is designed for digital " +
                    "forensics analysts and human rights investigators working directly " +
                    "with device owners.\n\n" +
                    "MESH is developed by the civil society group BARGHEST. " +
                    "It is fully open-source and independently audited.",
            style = MaterialTheme.typography.bodyLarge.copy(
                fontFamily = JetBrainsMonoFamily
            ),
            color = Color.White,
            textAlign = TextAlign.Center
        )
        Spacer(Modifier.weight(2f))
    }
}

@Composable
internal fun HowItWorksPage() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 24.dp, vertical = 32.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(Modifier.height(16.dp))
        Text(
            text = "How MESH Works",
            style = MaterialTheme.typography.headlineMedium.copy(fontWeight = FontWeight.Bold),
            color = MaterialTheme.colorScheme.onSurface
        )
        Spacer(Modifier.height(8.dp))
        Text(
            text = "MESH establishes a private, encrypted connection between " +
                    "your device and a trusted analyst for the duration of a session.",
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            textAlign = TextAlign.Center
        )
        Spacer(Modifier.height(24.dp))
        BulletItem("End-to-end encrypted using WireGuard")
        BulletItem("Peer-to-peer by default, with encrypted HTTPS relays when direct connection is unavailable")
        BulletItem("No device data is accessed without your explicit consent")
        BulletItem("Connections are temporary and exist only during an active session")
        BulletItem("You control when the connection starts and ends")
        BulletItem("The MESH app does not collect any data")
    }
}

@Composable
internal fun ConsentPage() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 24.dp, vertical = 32.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "You’re in control",
            style = MaterialTheme.typography.headlineMedium.copy(fontWeight = FontWeight.Bold),
            color = MaterialTheme.colorScheme.onSurface
        )
        Spacer(Modifier.height(8.dp))
        Text(
            text = "When you enable ADB access, a trusted analyst can assist " +
                    "with examining your device. Before proceeding, you should " +
                    "discuss and agree on what data will be accessed.\n\n" +
                    "With your consent, an analyst may:",
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            textAlign = TextAlign.Center
        )
        Spacer(Modifier.height(24.dp))
        BulletItem("Review selected files, application data, and system logs with you")
        BulletItem("Create a backup of installed applications")
        BulletItem("Install investigative tools required for the session")
        BulletItem("Run diagnostic commands on your device")
        Spacer(Modifier.height(24.dp))
        Text(
            text = "You decide when ADB is enabled or disabled. " +
                    "ADB access should remain off except during an agreed " +
                    "and active investigation session.",
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            textAlign = TextAlign.Center
        )
    }
}

@Composable
private fun BulletItem(text: String) {
    Row(
        modifier = Modifier.padding(vertical = 6.dp),
        verticalAlignment = Alignment.Top
    ) {
        Text(
            text = "\u2022",
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.surfaceContainer
        )
        Spacer(Modifier.width(12.dp))
        Text(
            text = text,
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )
    }
}