// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.view

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.statusBarsPadding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardCapitalization
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.R
import com.barghest.mesh.ui.theme.meshBg
import com.barghest.mesh.ui.theme.meshMuted
import com.barghest.mesh.ui.theme.meshText
import com.barghest.mesh.ui.theme.meshText2
import com.barghest.mesh.ui.util.set
import com.barghest.mesh.ui.view.mesh.Eyebrow
import com.barghest.mesh.ui.view.mesh.MeshButton
import com.barghest.mesh.ui.view.mesh.MeshTopBar
import com.barghest.mesh.ui.viewModel.LoginWithAuthKeyViewModel
import com.barghest.mesh.ui.viewModel.LoginWithCustomControlURLViewModel

data class LoginViewStrings(
    var title: String,
    var explanation: String,
    var inputTitle: String,
    var placeholder: String,
)

@Composable
fun LoginWithCustomControlURLView(
    onNavigateHome: () -> Unit,
    backToSettings: () -> Unit,
    onNavigateToAuthKey: () -> Unit,
    viewModel: LoginWithCustomControlURLViewModel = LoginWithCustomControlURLViewModel()
) {
    val error by viewModel.errorDialog.collectAsState()
    val successMessage by viewModel.successMessage.collectAsState()

    // Show success dialog when control plane is added
    successMessage?.let { message ->
        AlertDialog(
            onDismissRequest = { viewModel.clearSuccessMessage() },
            title = { Text("Success") },
            text = {
                Column {
                    Text(message)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        "Now use a pre-auth key to connect to your MESH network.",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
            },
            confirmButton = {
                TextButton(onClick = { viewModel.proceedToAuthKey(onNavigateToAuthKey) }) {
                    Text("Use Pre-Auth Key")
                }
            },
            dismissButton = {
                TextButton(onClick = { viewModel.clearSuccessMessage() }) { Text("Later") }
            }
        )
    }

    error?.let { ErrorDialog(type = it, action = { viewModel.errorDialog.set(null) }) }

    LoginScaffold(
        onBack = backToSettings,
        strings = LoginViewStrings(
            title = stringResource(id = R.string.custom_control_menu),
            explanation = stringResource(id = R.string.custom_control_menu_desc),
            inputTitle = stringResource(id = R.string.custom_control_url_title),
            placeholder = stringResource(id = R.string.custom_control_placeholder),
        ),
        onSubmitAction = { viewModel.setControlURL(it, onNavigateHome) },
    )
}

@Composable
fun LoginWithAuthKeyView(
    onNavigateHome: () -> Unit,
    backToSettings: () -> Unit,
    viewModel: LoginWithAuthKeyViewModel = LoginWithAuthKeyViewModel()
) {
    val error by viewModel.errorDialog.collectAsState()
    error?.let { ErrorDialog(type = it, action = { viewModel.errorDialog.set(null) }) }

    LoginScaffold(
        onBack = backToSettings,
        strings = LoginViewStrings(
            title = stringResource(id = R.string.auth_key_title),
            explanation = stringResource(id = R.string.auth_key_explanation),
            inputTitle = stringResource(id = R.string.auth_key_input_title),
            placeholder = stringResource(id = R.string.auth_key_placeholder),
        ),
        onSubmitAction = { viewModel.setAuthKey(it, onNavigateHome) },
    )
}

/** Shared login form for the control-server and auth-key screens. */
@Composable
private fun LoginScaffold(
    onBack: () -> Unit,
    strings: LoginViewStrings,
    onSubmitAction: (String) -> Unit,
) {
    val cs = MaterialTheme.colorScheme
    var textVal by remember { mutableStateOf("") }

    Column(Modifier.fillMaxSize().background(cs.meshBg).statusBarsPadding().imePadding()) {
        MeshTopBar(strings.inputTitle, onBack = onBack)
        Column(
            Modifier.fillMaxWidth().weight(1f).verticalScroll(rememberScrollState()).padding(horizontal = 20.dp),
        ) {
            Spacer(Modifier.height(20.dp))
            Text(strings.title, color = cs.meshText, fontSize = 21.sp, fontWeight = FontWeight.SemiBold, letterSpacing = (-0.3).sp)
            Spacer(Modifier.height(6.dp))
            Text(strings.explanation, color = cs.meshText2, fontSize = 13.5.sp, lineHeight = 20.sp)
            Spacer(Modifier.height(20.dp))
            Eyebrow(strings.inputTitle)
            Spacer(Modifier.height(10.dp))
            OutlinedTextField(
                modifier = Modifier.fillMaxWidth(),
                value = textVal,
                onValueChange = { textVal = it },
                singleLine = true,
                textStyle = MaterialTheme.typography.bodyMedium,
                placeholder = { Text(strings.placeholder, color = cs.meshMuted, style = MaterialTheme.typography.bodySmall) },
                colors = OutlinedTextFieldDefaults.colors(
                    focusedTextColor = cs.onSurface,
                    unfocusedTextColor = cs.onSurface,
                    focusedContainerColor = cs.surfaceContainerHigh,
                    unfocusedContainerColor = cs.surfaceContainerHigh,
                    focusedBorderColor = cs.primary,
                    unfocusedBorderColor = cs.outline,
                    cursorColor = cs.primary,
                ),
                keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.None, imeAction = ImeAction.Go),
                keyboardActions = KeyboardActions(onGo = { onSubmitAction(textVal) }),
            )
            Spacer(Modifier.height(20.dp))
            MeshButton(stringResource(id = R.string.add_account_short), { onSubmitAction(textVal) }, full = true, height = 56.dp)
        }
    }
}
