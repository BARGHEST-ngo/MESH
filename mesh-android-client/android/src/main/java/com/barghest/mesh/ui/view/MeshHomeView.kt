// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.R
import com.barghest.mesh.ui.model.Ipn
import com.barghest.mesh.ui.theme.AppTheme
import com.barghest.mesh.ui.util.AppVersion
import com.barghest.mesh.ui.viewModel.MainViewModel
import com.barghest.mesh.ui.viewModel.AppViewModel
import androidx.compose.ui.tooling.preview.Preview
import kotlinx.coroutines.flow.emptyFlow

data class MeshHomeNavigation(
    val onNavigateToSettings: () -> Unit,
    val onNavigateToAuthKey: () -> Unit,
    val onNavigateToCustomControl: () -> Unit,
    val onNavigateToMainView: () -> Unit,
    val onNavigateToADBSetup: () -> Unit
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MeshHomeView(
    navigation: MeshHomeNavigation,
    viewModel: MainViewModel,
    appViewModel: AppViewModel,
    isPresented: Boolean = true,
    onDismiss: () -> Unit = {}
) {
    val context = LocalContext.current
    val state by viewModel.ipnState.collectAsState(initial = Ipn.State.NoState)
    val user by viewModel.loggedInUser.collectAsState(initial = null)
    val isOn by viewModel.vpnToggleState.collectAsState(initial = false)
    val stateVal by viewModel.stateRes.collectAsState(initial = R.string.placeholder)
    val stateStr = stringResource(id = stateVal)
    val prefs by viewModel.prefs.collectAsState(initial = null)
    val controlURL = prefs?.ControlURL
    val hasCustomControl = !controlURL.isNullOrEmpty()

    val sheetState = rememberModalBottomSheetState(skipPartiallyExpanded = false)

    if (isPresented) {
        ModalBottomSheet(
            onDismissRequest = onDismiss,
            sheetState = sheetState,
            scrimColor = MaterialTheme.colorScheme.scrim.copy(alpha = 0.32f)
        ) {
            val lazyListState = rememberLazyListState()
            val isScrolled by remember {
                derivedStateOf {
                    lazyListState.firstVisibleItemIndex > 0 ||
                            lazyListState.firstVisibleItemScrollOffset > 0
                }
            }
            val topBarAlpha by animateFloatAsState(
                if (isScrolled) 0.95f else 0f,
                label = "topBarAlpha"
            )

            Box(Modifier.fillMaxWidth()) {
                LazyColumn(
                    state = lazyListState,
                    contentPadding = PaddingValues(top = 80.dp, bottom = 24.dp)
                ) {
                    item { MeshHeaderSection() }

                    item {
                        ConnectionStatusSection(
                            state = state,
                            isOn = isOn,
                            stateStr = stateStr,
                            user = user?.NetworkProfile?.tailnetNameForDisplay()
                        )
                    }

                    item {
                        ControlServerSection(
                            hasCustomControl = hasCustomControl,
                            controlURL = controlURL ?: "",
                            prefs = prefs?.toString() ?: ""
                        )
                    }

                    item {
                        ConnectionStepsSection(
                            navigation = navigation,
                            isConnected = user != null && state == Ipn.State.Running
                        )
                    }

                    item { FooterSectionModern() }
                }
                Box(
                    Modifier
                        .align(Alignment.TopCenter)
                        .fillMaxWidth()
                        .height(64.dp)
                        .background(MaterialTheme.colorScheme.background.copy(alpha = topBarAlpha))
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(horizontal = 16.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        IconButton(onClick = navigation.onNavigateToSettings) {
                            Icon(
                                Icons.Default.Settings,
                                contentDescription = stringResource(R.string.settings_description),
                                tint = MaterialTheme.colorScheme.onBackground
                            )
                        }
                        TextButton(
                            onClick = {
                                onDismiss()
                            }
                        ) {
                            Text(
                                text = stringResource(R.string.close_plain),
                                fontFamily = FontFamily.Monospace,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.onBackground
                            )
                        }
                    }
                }
            }
        }
    }
}


@Composable
private fun MeshHeaderSection() {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 8.dp)
    ) {
        Text(
            text = buildAnnotatedString {
                append(stringResource(R.string.app_name_mesh))
                withStyle(style = SpanStyle(fontSize = 10.sp)) {
                    append(" v${AppVersion.Short()}")
                }
            },
            fontFamily = FontFamily.Monospace,
            fontWeight = FontWeight.Bold,
            fontSize = 28.sp,
            color = MaterialTheme.colorScheme.onBackground
        )
        Text(
            text = stringResource(R.string.tagline_mesh),
            fontFamily = FontFamily.Monospace,
            fontSize = 12.sp,
            color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.7f)
        )
    }
}

@Composable
private fun ConnectionStatusSection(
    state: Ipn.State,
    isOn: Boolean,
    stateStr: String,
    user: String?
) {
    val color = when {
        state == Ipn.State.Running -> if (isOn) Color(0xFF248A3D) else Color(0xFFFF9500)
        else -> MaterialTheme.colorScheme.error
    }

    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 8.dp),
        color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.25f),
        shape = RoundedCornerShape(12.dp)
    ) {
        Column(
            Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Surface(color = color, shape = CircleShape) { Box(Modifier.size(8.dp)) }
                Spacer(Modifier.width(8.dp))
                Text(
                    text = if (state == Ipn.State.Running) "[ CONNECTED ]" else "[ DISCONNECTED ]",
                    fontFamily = FontFamily.Monospace,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onBackground
                )
            }
            if (!user.isNullOrEmpty()) {
                Text(
                    text = user,
                    fontFamily = FontFamily.Monospace,
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                )
            }
            Text(
                text = stateStr,
                fontFamily = FontFamily.Monospace,
                fontSize = 11.sp,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
            )
        }
    }
}

@Composable
private fun ControlServerSection(
    hasCustomControl: Boolean,
    controlURL: String,
    prefs: String
) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 8.dp),
        color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.25f),
        shape = RoundedCornerShape(12.dp)
    ) {
        Column(
            Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Surface(
                    color = if (hasCustomControl) Color(0xFF248A3D) else MaterialTheme.colorScheme.error,
                    shape = CircleShape
                ) { Box(Modifier.size(8.dp)) }
                Spacer(Modifier.width(8.dp))
                Text(
                    text = if (hasCustomControl) stringResource(R.string.control_configured) else stringResource(
                        R.string.control_not_set
                    ),
                    fontFamily = FontFamily.Monospace,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onBackground
                )
            }
            if (hasCustomControl) {
                Text(
                    text = controlURL,
                    fontFamily = FontFamily.Monospace,
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                )
            }
        }
    }
}

@Composable
private fun ConnectionStepsSection(
    navigation: MeshHomeNavigation,
    isConnected: Boolean
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        Text(
            text = if (isConnected) stringResource(R.string.connection_options)
            else stringResource(R.string.connect_steps_title),
            fontFamily = FontFamily.Monospace,
            fontWeight = FontWeight.Medium,
            color = MaterialTheme.colorScheme.onBackground
        )

        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            ConnectionStepButton(
                text = stringResource(R.string.step_1_enable_adb),
                onClick = navigation.onNavigateToADBSetup
            )
            ConnectionStepButton(
                text = stringResource(R.string.step_2_control_server),
                onClick = navigation.onNavigateToCustomControl
            )
            ConnectionStepButton(
                text = stringResource(R.string.step_3_auth_key),
                onClick = navigation.onNavigateToAuthKey
            )
        }
    }
}

@Composable
private fun ColumnScope.ConnectionStepButton(text: String, onClick: () -> Unit) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .height(64.dp)
            .border(
                width = 1.dp,
                color = Color(0xFF494A4D),
                shape = RoundedCornerShape(10.dp)
            ),
        color = Color(0xFF2E2E2E),
        shape = RoundedCornerShape(10.dp),
        onClick = onClick
    ) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Text(
                text = text,
                fontFamily = FontFamily.Monospace,
                fontSize = 11.sp,
                color = MaterialTheme.colorScheme.onBackground
            )
        }
    }
}


@Composable
private fun FooterSectionModern() {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 24.dp, vertical = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        Text(
            text = stringResource(R.string.footer_brand),
            fontSize = 11.sp,
            fontFamily = FontFamily.Monospace,
            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
        )
        Text(
            text = "Debug Settings",
            fontSize = 11.sp,
            fontFamily = FontFamily.Monospace,
            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
        )
    }
}

@Preview
@Composable
fun MeshHomeViewPreview() {
    AppTheme {
        val appViewModel = AppViewModel(LocalContext.current.applicationContext as android.app.Application, emptyFlow())
        MeshHomeView(
            isPresented = true,
            onDismiss = {},
            navigation = MeshHomeNavigation(
                onNavigateToSettings = {},
                onNavigateToAuthKey = {},
                onNavigateToCustomControl = {},
                onNavigateToMainView = {},
                onNavigateToADBSetup = {}
            ),
            viewModel = MainViewModel(appViewModel),
            appViewModel = appViewModel
        )
    }
}


