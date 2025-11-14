// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import android.annotation.SuppressLint
import android.content.ActivityNotFoundException
import android.content.Context
import android.content.Intent
import android.os.Build
import android.provider.Settings
import android.widget.Toast
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp

@Composable
fun ADBSetup(
    onNavigateHome: () -> Unit
) {
    val context = LocalContext.current
    val scroll = rememberScrollState()
    var showTroubleshooting by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            Header(
                title = (@androidx.compose.runtime.Composable { Text("Enable ADB") } as @Composable () -> Unit),
                onBack = onNavigateHome
            )
        },
        bottomBar = {
            Surface(tonalElevation = 2.dp) {
                Row(
                    modifier = Modifier
                        .navigationBarsPadding()
                        .imePadding()
                        .fillMaxWidth()
                        .padding(16.dp),
                    horizontalArrangement = Arrangement.spacedBy(12.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    OutlinedButton(
                        modifier = Modifier.weight(1f),
                        onClick = { showTroubleshooting = !showTroubleshooting }
                    ) {
                        Text(if (showTroubleshooting) "[ HIDE ]" else "[ TROUBLESHOOT ]")
                    }
                    Button(
                        modifier = Modifier.weight(1f),
                        onClick = { openAdbWireless(context) }
                    ) {
                        Text("[ OPEN DEV OPTIONS]")
                    }
                }
            }
        }
    ) { paddingInsets ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingInsets)
                .verticalScroll(scroll)
                .padding(horizontal = 20.dp, vertical = 16.dp),
            horizontalAlignment = Alignment.Start
        ) {
            Text(
                text = "ADB Setup",
                style = MaterialTheme.typography.headlineMedium.copy(fontWeight = FontWeight.Bold),
                color = MaterialTheme.colorScheme.onSurface
            )
            Spacer(Modifier.height(6.dp))
            Text(
                text = "Enable Developer options and Wireless debugging to connect over Wi-Fi.",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(Modifier.height(16.dp))
            ElevatedCard(
                shape = RoundedCornerShape(16.dp)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    verticalAlignment = Alignment.Top
                ) {
                    Text(
                        text = "On some devices the exact menu names differ slightly (e.g., “About phone” → “Software information”). The flow below is the same.",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
            }
            Spacer(Modifier.height(20.dp))
            StepItem(
                stepNumber = 1,
                title = "Open system Settings",
                body = "Go to Settings → About phone → Software information."
            )
            Divider(modifier = Modifier.padding(vertical = 8.dp))
            StepItem(
                stepNumber = 2,
                title = "Enable Developer options",
                body = "Find “Build number” and tap it 7 times. Enter your PIN if prompted. You’ll see “Developer mode has been enabled”."
            )
            Divider(modifier = Modifier.padding(vertical = 8.dp))
            StepItem(
                stepNumber = 3,
                title = "Open Developer options",
                body = "Use the button below to jump there directly. If it doesn’t open the exact page, navigate manually: Settings → System → Developer options."
            )
            Divider(modifier = Modifier.padding(vertical = 8.dp))
            StepItem(
                stepNumber = 4,
                title = "Turn on Wireless debugging",
                body = "Inside Developer options, enable “Wireless debugging”. Confirm any prompts."
            )
            Divider(modifier = Modifier.padding(vertical = 8.dp))
            StepItem(
                stepNumber = 5,
                title = "Return here",
                body = "Come back and continue the setup."
            )
            Spacer(Modifier.height(12.dp))
            AssistiveNote(
                text = "Only enable Wireless debugging when you intend to use it. Turn it off afterwards for better security."
            )
            Spacer(Modifier.height(12.dp))
            AssistiveAccordion(
                expanded = showTroubleshooting,
                onToggle = { showTroubleshooting = !showTroubleshooting },
                title = "Troubleshooting & OEM notes"
            ) {
                Text(
                    "• The Developer options menu might live under “System” on some devices.\n" +
                            "• If “Wireless debugging” is missing, your device may not support it or your OEM has hidden it.\n" +
                            "• Wireless debugging requires Android 11 (API 30) or newer.\n" +
                            "• If the direct link fails, open Settings and search for “Developer options”.",
                    style = MaterialTheme.typography.bodyMedium
                )
            }
            Spacer(Modifier.height(96.dp))
        }
    }
}

@Composable
private fun StepItem(
    stepNumber: Int,
    title: String,
    body: String
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp)
            .semantics { contentDescription = "Step $stepNumber: $title" },
        verticalAlignment = Alignment.Top
    ) {
        Surface(
            modifier = Modifier.size(28.dp),
            shape = CircleShape,
            color = MaterialTheme.colorScheme.primaryContainer,
            tonalElevation = 2.dp
        ) {
            Box(contentAlignment = Alignment.Center, modifier = Modifier.fillMaxSize()) {
                Text(
                    text = stepNumber.toString(),
                    style = MaterialTheme.typography.labelLarge.copy(fontWeight = FontWeight.Bold),
                    color = MaterialTheme.colorScheme.onPrimaryContainer
                )
            }
        }
        Spacer(Modifier.width(12.dp))
        Column(Modifier.weight(1f)) {
            Text(
                text = title,
                style = MaterialTheme.typography.titleMedium.copy(fontWeight = FontWeight.SemiBold)
            )
            Spacer(Modifier.height(4.dp))
            Text(
                text = body,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
        }
    }
}

@Composable
private fun AssistiveNote(
    text: String
) {
    Surface(
        shape = RoundedCornerShape(12.dp),
        tonalElevation = 1.dp,
        modifier = Modifier.fillMaxWidth()
    ) {
        Row(modifier = Modifier.padding(12.dp), verticalAlignment = Alignment.Top) {
            Text(text, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
        }
    }
}

@Composable
private fun AssistiveAccordion(
    expanded: Boolean,
    onToggle: () -> Unit,
    title: String,
    content: @Composable () -> Unit
) {
    OutlinedCard(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(14.dp)
    ) {
        Column(Modifier.fillMaxWidth()) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickableWithoutRipple(onToggle)
                    .padding(horizontal = 12.dp, vertical = 14.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(title, style = MaterialTheme.typography.titleSmall)
            }
            AnimatedVisibility(visible = expanded) {
                Column(Modifier.padding(horizontal = 12.dp, vertical = 8.dp)) {
                    content()
                    Spacer(Modifier.height(8.dp))
                }
            }
        }
    }
}

@SuppressLint("SuspiciousModifierThen")
private fun Modifier.clickableWithoutRipple(onClick: () -> Unit): Modifier =
    this.then(
        clickable(
            indication = null,
            interactionSource = MutableInteractionSource(),
            onClick = onClick
        )
    )

private fun openAdbWireless(context: Context): Boolean {
    try {
        val intents = buildList {
            add(Intent(Settings.ACTION_APPLICATION_DEVELOPMENT_SETTINGS).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK))
            add(Intent(Settings.ACTION_SETTINGS).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK))
        }

        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) {
            Toast.makeText(context, "Wireless debugging requires Android 11 or higher", Toast.LENGTH_LONG).show()
        }

        for (intent in intents) {
            val resolved = context.packageManager.resolveActivity(intent, 0)
            if (resolved != null) {
                context.startActivity(intent)
                if (intent.action == Settings.ACTION_SETTINGS) {
                    Toast.makeText(
                        context,
                        "Navigate to System → Developer options → Wireless debugging",
                        Toast.LENGTH_LONG
                    ).show()
                }
                return true
            }
        }

        Toast.makeText(
            context,
            "Unable to open settings. Please enable Wireless debugging manually.",
            Toast.LENGTH_LONG
        ).show()
        return false
    } catch (e: SecurityException) {
        Toast.makeText(
            context,
            "Unable to open settings due to a permission issue. Please enable manually.",
            Toast.LENGTH_LONG
        ).show()
        return false
    } catch (e: ActivityNotFoundException) {
        Toast.makeText(
            context,
            "Settings activity not found. Please enable manually.",
            Toast.LENGTH_LONG
        ).show()
        return false
    }
}

@Preview(showBackground = true, name = "ADB Setup – Light", widthDp = 360)
@Composable
private fun PreviewAdbSetupLight() {
    MaterialTheme {
        ADBSetup(onNavigateHome = {})
    }
}

@Preview(showBackground = true, name = "ADB Setup – Dark", widthDp = 360)
@Composable
private fun PreviewAdbSetupDark() {
    MaterialTheme(colorScheme = darkColorScheme()) {
        ADBSetup(onNavigateHome = {})
    }
}
