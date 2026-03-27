// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.util

import android.provider.Settings
import android.content.Context
import android.os.Build
import androidx.annotation.RequiresApi
import java.time.LocalDate
import java.time.temporal.ChronoUnit

object SecurityChecker {

fun isDevOptionsEnabled(context: Context): Boolean {
    return try {
        Settings.Global.getInt(
            context.contentResolver,
            Settings.Global.DEVELOPMENT_SETTINGS_ENABLED,
            0
        ) != 0
    } catch (e: Exception) {
        false
    }
}

@RequiresApi(Build.VERSION_CODES.O)
fun isSPLstale(thresholdDays: Int): Boolean {
	if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) return false
	try {
		val patch = Build.VERSION.SECURITY_PATCH
		val today = LocalDate.now()
		val patchDate = LocalDate.parse(patch)
		val daysSince = ChronoUnit.DAYS.between(patchDate, today)
		return daysSince > thresholdDays
	} catch (e: Exception) {
		return false
	}
}
}
