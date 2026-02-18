// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view

import android.os.Build
import androidx.compose.foundation.clickable
import androidx.compose.foundation.focusable
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalClipboardManager
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.AnnotatedString
import com.barghest.mesh.R
import com.barghest.mesh.ui.theme.listItem
import com.barghest.mesh.ui.util.AndroidTVUtil.isAndroidTV
import com.barghest.mesh.ui.util.itemsWithDividers

data class DeviceInfoItem(
    val title: String,
    val value: String
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DeviceInfoView(onNavigateBack: () -> Unit) {
    val clipboardManager = LocalClipboardManager.current
    
    val deviceInfo = buildList {
        add(DeviceInfoItem("Device Model", Build.MODEL))
        add(DeviceInfoItem("Manufacturer", Build.MANUFACTURER))
        add(DeviceInfoItem("Android Version", Build.VERSION.RELEASE))
        add(DeviceInfoItem("Build Number", Build.DISPLAY))
        add(DeviceInfoItem("API Level", Build.VERSION.SDK_INT.toString()))
        
        // Security patch level is only available on API 23+
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            add(DeviceInfoItem("Security Patch", Build.VERSION.SECURITY_PATCH))
        }
    }
    
    Scaffold(
        topBar = {
            Header(
                titleRes = R.string.device_information,
                onBack = onNavigateBack
            )
        }
    ) { innerPadding ->
        LazyColumn(
            modifier = Modifier.padding(innerPadding)
        ) {
            itemsWithDividers(deviceInfo, key = { "device_info_${it.title}" }) { item ->
                DeviceInfoRow(
                    title = item.title,
                    value = item.value,
                    onCopy = {
                        if (!isAndroidTV()) {
                            clipboardManager.setText(AnnotatedString(item.value))
                        }
                    }
                )
            }
        }
    }
}

@Composable
fun DeviceInfoRow(title: String, value: String, onCopy: () -> Unit) {
    val modifier = if (isAndroidTV()) {
        Modifier.focusable(false)
    } else {
        Modifier.clickable { onCopy() }
    }
    
    ListItem(
        modifier = modifier,
        colors = MaterialTheme.colorScheme.listItem,
        headlineContent = { Text(text = title) },
        supportingContent = { Text(text = value) },
        trailingContent = {
            if (!isAndroidTV()) {
                Icon(
                    painter = painterResource(id = R.drawable.clipboard),
                    contentDescription = "Copy to clipboard"
                )
            }
        }
    )
}

