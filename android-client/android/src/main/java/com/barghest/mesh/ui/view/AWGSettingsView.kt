// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import android.widget.Toast
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.barghest.mesh.R
import com.barghest.mesh.ui.theme.listItem
import com.barghest.mesh.ui.theme.off
import com.barghest.mesh.ui.theme.success
import com.barghest.mesh.ui.util.Lists
import com.barghest.mesh.ui.viewModel.AWGConfig
import com.barghest.mesh.ui.viewModel.AWGSettingsViewModel
import com.barghest.mesh.ui.viewModel.AWGSettingsViewModelFactory

@Composable
fun AWGSettingsView(
    backToSettings: BackNavigation,
    model: AWGSettingsViewModel = viewModel(factory = AWGSettingsViewModelFactory())
) {
    val config by model.config.collectAsState()
    val validationError by model.validationError.collectAsState()
    val saveStatus by model.saveStatus.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(saveStatus) {
        when (saveStatus) {
            "saved" -> {
                Toast.makeText(context, R.string.awg_restart_required, Toast.LENGTH_LONG).show()
                model.clearSaveStatus()
            }
            "error" -> {
                Toast.makeText(context, R.string.awg_validation_error, Toast.LENGTH_SHORT).show()
                model.clearSaveStatus()
            }
        }
    }

    Scaffold(topBar = { Header(R.string.awg_settings, onBack = backToSettings) }) { innerPadding ->
        Column(
            modifier = Modifier
                .padding(innerPadding)
                .verticalScroll(rememberScrollState())
        ) {
            // Status indicator
            ListItem(
                colors = MaterialTheme.colorScheme.listItem,
                leadingContent = {
                    Icon(
                        painter = painterResource(
                            if (config.isObfuscationEnabled()) R.drawable.check_circle
                            else R.drawable.xmark_circle
                        ),
                        contentDescription = null,
                        tint = if (config.isObfuscationEnabled())
                            MaterialTheme.colorScheme.success
                        else MaterialTheme.colorScheme.off
                    )
                },
                headlineContent = {
                    Text(
                        stringResource(
                            if (config.isObfuscationEnabled()) R.string.awg_obfuscation_enabled
                            else R.string.awg_obfuscation_disabled
                        ),
                        style = MaterialTheme.typography.titleMedium
                    )
                },
                supportingContent = {
                    Text(stringResource(R.string.awg_settings_subtitle))
                }
            )

            // Validation error
            validationError?.let { error ->
                Lists.ItemDivider()
                ListItem(
                    colors = MaterialTheme.colorScheme.listItem,
                    headlineContent = {
                        Text(error, color = MaterialTheme.colorScheme.error)
                    }
                )
            }

            // Junk Packets Section
            Lists.SectionDivider(stringResource(R.string.awg_junk_packets))
            AWGNumberField(
                label = stringResource(R.string.awg_jc),
                description = stringResource(R.string.awg_jc_desc),
                value = config.jc,
                onValueChange = { model.updateConfig(config.copy(jc = it)) }
            )
            AWGNumberField(
                label = stringResource(R.string.awg_jmin),
                description = stringResource(R.string.awg_jmin_desc),
                value = config.jmin,
                onValueChange = { model.updateConfig(config.copy(jmin = it)) }
            )
            AWGNumberField(
                label = stringResource(R.string.awg_jmax),
                description = stringResource(R.string.awg_jmax_desc),
                value = config.jmax,
                onValueChange = { model.updateConfig(config.copy(jmax = it)) }
            )

            // Handshake Padding Section
            Lists.SectionDivider(stringResource(R.string.awg_handshake_padding))
            AWGNumberField(
                label = stringResource(R.string.awg_s1),
                description = stringResource(R.string.awg_s1_desc),
                value = config.s1,
                onValueChange = { model.updateConfig(config.copy(s1 = it)) }
            )
            AWGNumberField(
                label = stringResource(R.string.awg_s2),
                description = stringResource(R.string.awg_s2_desc),
                value = config.s2,
                onValueChange = { model.updateConfig(config.copy(s2 = it)) }
            )

            // Message Type Section
            Lists.SectionDivider(stringResource(R.string.awg_message_types))
            Text(
                stringResource(R.string.awg_h_desc),
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
            )
            AWGHValuesRow(config, model)

            // Action Buttons
            Spacer(modifier = Modifier.height(24.dp))
            Row(modifier = Modifier.padding(horizontal = 16.dp)) {
                OutlinedButton(
                    onClick = { model.resetToDefaults() },
                    modifier = Modifier.weight(1f)
                ) {
                    Text(stringResource(R.string.awg_reset_defaults))
                }
                Spacer(modifier = Modifier.width(16.dp))
                Button(
                    onClick = { model.saveConfig() },
                    modifier = Modifier.weight(1f),
                    enabled = validationError == null
                ) {
                    Text(stringResource(R.string.awg_save))
                }
            }
            Spacer(modifier = Modifier.height(16.dp))
        }
    }
}

@Composable
private fun AWGNumberField(
    label: String,
    description: String,
    value: Int,
    onValueChange: (Int) -> Unit
) {
    ListItem(
        colors = MaterialTheme.colorScheme.listItem,
        headlineContent = { Text(label, style = MaterialTheme.typography.bodyMedium) },
        supportingContent = {
            Text(description, style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant)
        },
        trailingContent = {
            OutlinedTextField(
                value = value.toString(),
                onValueChange = { newValue ->
                    onValueChange(newValue.toIntOrNull() ?: 0)
                },
                modifier = Modifier.width(80.dp),
                singleLine = true,
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number)
            )
        }
    )
}

@Composable
private fun AWGHValuesRow(config: AWGConfig, model: AWGSettingsViewModel) {
    Row(modifier = Modifier.padding(horizontal = 16.dp)) {
        AWGSmallNumberField(
            label = stringResource(R.string.awg_h1),
            value = config.h1,
            onValueChange = { model.updateConfig(config.copy(h1 = it)) },
            modifier = Modifier.weight(1f)
        )
        Spacer(modifier = Modifier.width(8.dp))
        AWGSmallNumberField(
            label = stringResource(R.string.awg_h2),
            value = config.h2,
            onValueChange = { model.updateConfig(config.copy(h2 = it)) },
            modifier = Modifier.weight(1f)
        )
        Spacer(modifier = Modifier.width(8.dp))
        AWGSmallNumberField(
            label = stringResource(R.string.awg_h3),
            value = config.h3,
            onValueChange = { model.updateConfig(config.copy(h3 = it)) },
            modifier = Modifier.weight(1f)
        )
        Spacer(modifier = Modifier.width(8.dp))
        AWGSmallNumberField(
            label = stringResource(R.string.awg_h4),
            value = config.h4,
            onValueChange = { model.updateConfig(config.copy(h4 = it)) },
            modifier = Modifier.weight(1f)
        )
    }
}

@Composable
private fun AWGSmallNumberField(
    label: String,
    value: Int,
    onValueChange: (Int) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier) {
        Text(label, style = MaterialTheme.typography.labelSmall)
        OutlinedTextField(
            value = value.toString(),
            onValueChange = { newValue ->
                onValueChange(newValue.toIntOrNull() ?: 1)
            },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number)
        )
    }
}

