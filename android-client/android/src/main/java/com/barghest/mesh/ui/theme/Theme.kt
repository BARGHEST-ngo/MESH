// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.ButtonColors
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.ColorScheme
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ListItemColors
import androidx.compose.material3.ListItemDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.TextFieldColors
import androidx.compose.material3.TopAppBarColors
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.material3.Typography
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.DisposableEffect
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.TextUnit
import androidx.compose.ui.unit.sp
import com.google.accompanist.systemuicontroller.rememberSystemUiController

private const val BASE = 16 // sp

private val AppTypography =
    Typography(
        displayLarge = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Bold, fontSize = (BASE + 20).sp, lineHeight = (BASE + 28).sp),
        displayMedium = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Bold, fontSize = (BASE + 14).sp, lineHeight = (BASE + 22).sp),
        displaySmall = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.SemiBold, fontSize = (BASE + 10).sp, lineHeight = (BASE + 18).sp),

        headlineLarge = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.SemiBold, fontSize = (BASE + 8).sp, lineHeight = (BASE + 16).sp),
        headlineMedium = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE + 6).sp, lineHeight = (BASE + 14).sp),
        headlineSmall = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE + 4).sp, lineHeight = (BASE + 12).sp),

        titleLarge = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE + 4).sp, lineHeight = (BASE + 10).sp),
        titleMedium = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = 18.sp, lineHeight = 26.sp),
        titleSmall = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE + 1).sp, lineHeight = (BASE + 8).sp),

        bodyLarge = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Normal, fontSize = (BASE + 2).sp, lineHeight = (BASE + 8).sp),
        bodyMedium = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Normal, fontSize = 16.sp, lineHeight = 26.sp),
        bodySmall = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Normal, fontSize = (BASE - 2).sp, lineHeight = (BASE + 4).sp),

        labelLarge = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE - 1).sp, lineHeight = (BASE + 4).sp),
        labelMedium = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Medium, fontSize = (BASE - 2).sp, lineHeight = (BASE + 3).sp),
        labelSmall = TextStyle(fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Normal, fontSize = (BASE - 4).sp, lineHeight = (BASE + 1).sp),
    )

@Composable
fun AppTheme(useDarkTheme: Boolean = true, content: @Composable () -> Unit) {
    val colorScheme = if (useDarkTheme) createDarkColorScheme() else createLightColorScheme()

    @Suppress("deprecation")
    val systemUiController = rememberSystemUiController()
    DisposableEffect(systemUiController, useDarkTheme) {
        systemUiController.setStatusBarColor(color = colorScheme.surfaceContainer)
        systemUiController.setNavigationBarColor(color = if (useDarkTheme) DarkColors.Background else LightColors.Background)
        onDispose {}
    }

    MaterialTheme(colorScheme = colorScheme, typography = AppTypography, content = content)
}

@Composable
private fun createLightColorScheme(): ColorScheme {
    return lightColorScheme(
        primary = LightColors.Primary,
        onPrimary = LightColors.PrimaryForeground,
        primaryContainer = LightColors.Secondary,
        onPrimaryContainer = LightColors.SecondaryForeground,
        secondary = LightColors.Secondary,
        onSecondary = LightColors.SecondaryForeground,
        secondaryContainer = LightColors.Accent,
        onSecondaryContainer = LightColors.AccentForeground,
        tertiary = LightColors.Accent,
        onTertiary = LightColors.AccentForeground,
        tertiaryContainer = LightColors.Muted,
        onTertiaryContainer = LightColors.MutedForeground,
        error = LightColors.Destructive,
        onError = LightColors.DestructiveForeground,
        errorContainer = LightColors.Destructive,
        onErrorContainer = LightColors.DestructiveForeground,
        background = LightColors.Background,
        onBackground = LightColors.Foreground,
        surface = LightColors.Card,
        onSurface = LightColors.CardForeground,
        surfaceVariant = LightColors.Popover,
        onSurfaceVariant = LightColors.PopoverForeground,
        outline = LightColors.Border,
        outlineVariant = LightColors.Border,
        scrim = Color(0x00000000)
    )
}

@Composable
private fun createDarkColorScheme(): ColorScheme {
    return darkColorScheme(
        primary = DarkColors.Primary,
        onPrimary = DarkColors.PrimaryForeground,
        primaryContainer = DarkColors.Secondary,
        onPrimaryContainer = DarkColors.SecondaryForeground,
        secondary = DarkColors.Secondary,
        onSecondary = DarkColors.SecondaryForeground,
        secondaryContainer = DarkColors.Accent,
        onSecondaryContainer = DarkColors.AccentForeground,
        tertiary = DarkColors.Accent,
        onTertiary = DarkColors.AccentForeground,
        tertiaryContainer = DarkColors.Muted,
        onTertiaryContainer = DarkColors.MutedForeground,
        error = DarkColors.Destructive,
        onError = DarkColors.DestructiveForeground,
        errorContainer = DarkColors.Destructive,
        onErrorContainer = DarkColors.DestructiveForeground,
        background = DarkColors.Background,
        onBackground = DarkColors.Foreground,
        surface = DarkColors.Card,
        onSurface = DarkColors.CardForeground,
        surfaceVariant = DarkColors.Popover,
        onSurfaceVariant = DarkColors.PopoverForeground,
        outline = DarkColors.Border,
        outlineVariant = DarkColors.Border,
        scrim = Color(0x00000000)
    )
}



val ColorScheme.warning: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFFE5A300) else Color(0xFFCC8A00)

val ColorScheme.onWarning: Color get() = Color(0xFF000000)

val ColorScheme.warningContainer: Color
    @Composable
    get() = if (isSystemInDarkTheme()) Color(0xFF352A00) else Color(0xFFFFF8E1)

val ColorScheme.onWarningContainer: Color
    @Composable
    get() = if (isSystemInDarkTheme()) Color(0xFFFFE08A) else Color(0xFF4A3800)

val ColorScheme.success: Color
    @Composable
    get() = if (isSystemInDarkTheme()) Color(0xFF2ECB10) else Color(0xFF0F8A0F)

val ColorScheme.onSuccess: Color get() = Color(0xFF000000)

val ColorScheme.successContainer: Color
    @Composable
    get() = if (isSystemInDarkTheme()) Color(0xFF0E2A0E) else Color(0xFFE6F5E6)

val ColorScheme.onSuccessContainer: Color
    @Composable
    get() = if (isSystemInDarkTheme()) Color(0xFF9BFF8F) else Color(0xFF0A4D0A)

val ColorScheme.on: Color get() = Color(0xFFFFFFFF)

val ColorScheme.off: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFF303030) else Color(0xFFE0E0E0)

val ColorScheme.link: Color get() = onPrimaryContainer

val ColorScheme.customError: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFFFF5555) else Color(0xFFCC0000)

val ColorScheme.customErrorContainer: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFF4A0E11) else Color(0xFFFFEAEA)

val ColorScheme.listItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = default.containerColor,
            headlineColor = default.headlineColor,
            leadingIconColor = MaterialTheme.colorScheme.onPrimaryContainer,
            overlineColor = default.overlineColor,
            supportingTextColor = default.supportingTextColor,
            trailingIconColor = MaterialTheme.colorScheme.onPrimaryContainer,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.titledListItem: ListItemColors
    @Composable
    get() {
        val default = listItem
        return ListItemColors(
            containerColor = default.containerColor,
            headlineColor = default.headlineColor,
            leadingIconColor = default.leadingIconColor,
            overlineColor = MaterialTheme.colorScheme.onSurface,
            supportingTextColor = default.supportingTextColor,
            trailingIconColor = default.trailingIconColor,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.disabledListItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = default.containerColor,
            headlineColor = MaterialTheme.colorScheme.disabled,
            leadingIconColor = default.leadingIconColor,
            overlineColor = default.overlineColor,
            supportingTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
            trailingIconColor = default.trailingIconColor,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.surfaceContainerListItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = MaterialTheme.colorScheme.surfaceContainer,
            headlineColor = MaterialTheme.colorScheme.onSurface,
            leadingIconColor = MaterialTheme.colorScheme.onSurface,
            overlineColor = MaterialTheme.colorScheme.onSurfaceVariant,
            supportingTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
            trailingIconColor = MaterialTheme.colorScheme.onSurface,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.primaryListItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = MaterialTheme.colorScheme.primary,
            headlineColor = MaterialTheme.colorScheme.onPrimary,
            leadingIconColor = MaterialTheme.colorScheme.onPrimary,
            overlineColor = MaterialTheme.colorScheme.onPrimary.copy(alpha = 0.7f),
            supportingTextColor = MaterialTheme.colorScheme.onPrimary,
            trailingIconColor = MaterialTheme.colorScheme.onPrimary,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.warningListItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = MaterialTheme.colorScheme.warning,
            headlineColor = MaterialTheme.colorScheme.onPrimary,
            leadingIconColor = MaterialTheme.colorScheme.onPrimary,
            overlineColor = MaterialTheme.colorScheme.onPrimary.copy(alpha = 0.8f),
            supportingTextColor = MaterialTheme.colorScheme.onPrimary,
            trailingIconColor = MaterialTheme.colorScheme.onPrimary,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

val ColorScheme.errorListItem: ListItemColors
    @Composable
    get() {
        val default = ListItemDefaults.colors()
        return ListItemColors(
            containerColor = MaterialTheme.colorScheme.customError,
            headlineColor = MaterialTheme.colorScheme.onPrimary,
            leadingIconColor = MaterialTheme.colorScheme.onPrimary,
            overlineColor = MaterialTheme.colorScheme.onPrimary.copy(alpha = 0.8f),
            supportingTextColor = MaterialTheme.colorScheme.onPrimary,
            trailingIconColor = MaterialTheme.colorScheme.onPrimary,
            disabledHeadlineColor = default.disabledHeadlineColor,
            disabledLeadingIconColor = default.disabledLeadingIconColor,
            disabledTrailingIconColor = default.disabledTrailingIconColor)
    }

@OptIn(ExperimentalMaterial3Api::class)
val ColorScheme.topAppBar: TopAppBarColors
    @Composable
    get() =
        TopAppBarDefaults.topAppBarColors()
            .copy(
                containerColor = MaterialTheme.colorScheme.surfaceContainer,
                navigationIconContentColor = MaterialTheme.colorScheme.onSurface,
                titleContentColor = MaterialTheme.colorScheme.onSurface,
            )

val ColorScheme.secondaryButton: ButtonColors
    @Composable
    get() {
        val defaults = ButtonDefaults.buttonColors()
        return if (isSystemInDarkTheme()) {
            ButtonColors(
                containerColor = Color(0xFF5B7FFF),
                contentColor = Color(0xFFFFFFFF),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        } else {
            ButtonColors(
                containerColor = Color(0xFF5B7FFF),
                contentColor = Color(0xFFFFFFFF),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        }
    }

val ColorScheme.errorButton: ButtonColors
    @Composable
    get() {
        val defaults = ButtonDefaults.buttonColors()
        return if (isSystemInDarkTheme()) {
            ButtonColors(
                containerColor = Color(0xFFFF5555),
                contentColor = Color(0xFF000000),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        } else {
            ButtonColors(
                containerColor = Color(0xFFCC0000),
                contentColor = Color(0xFFFFFFFF),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        }
    }

val ColorScheme.warningButton: ButtonColors
    @Composable
    get() {
        val defaults = ButtonDefaults.buttonColors()
        return if (isSystemInDarkTheme()) {
            ButtonColors(
                containerColor = Color(0xFFE5A300),
                contentColor = Color(0xFF000000),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        } else {
            ButtonColors(
                containerColor = Color(0xFFCC8A00),
                contentColor = Color(0xFFFFFFFF),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        }
    }

val ColorScheme.defaultTextColor: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFFFFFFFF) else Color(0xFFFFFFFF)

val ColorScheme.logoBackground: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFFFFFFFF) else Color(0xFF0A0A0A)

val ColorScheme.standaloneLogoDotEnabled: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFFFFFFFF) else Color(0xFF000000)

val ColorScheme.standaloneLogoDotDisabled: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0x66FFFFFF) else Color(0x661F1E1E)

val ColorScheme.onBackgroundLogoDotEnabled: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0xFF141414) else Color(0xFFFFFFFF)

val ColorScheme.onBackgroundLogoDotDisabled: Color
    @Composable get() = if (isSystemInDarkTheme()) Color(0x66141414) else Color(0x66FFFFFF)

val ColorScheme.exitNodeToggleButton: ButtonColors
    @Composable
    get() {
        val defaults = ButtonDefaults.buttonColors()
        return if (isSystemInDarkTheme()) {
            ButtonColors(
                containerColor = Color(0xFF303030),
                contentColor = Color(0xFFFFFFFF),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        } else {
            ButtonColors(
                containerColor = Color(0xFFEDEBEA),
                contentColor = Color(0xFF000000),
                disabledContainerColor = defaults.disabledContainerColor,
                disabledContentColor = defaults.disabledContentColor)
        }
    }

val ColorScheme.disabled: Color get() = Color(0xFF9E9E9E)

@OptIn(ExperimentalMaterial3Api::class)
val ColorScheme.searchBarColors: TextFieldColors
    @Composable
    get() =
        OutlinedTextFieldDefaults.colors(
            focusedLeadingIconColor = MaterialTheme.colorScheme.onSurface,
            unfocusedLeadingIconColor = MaterialTheme.colorScheme.onSurface,
            focusedTextColor = MaterialTheme.colorScheme.onSurface,
            unfocusedTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
            disabledTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
            focusedContainerColor = MaterialTheme.colorScheme.surfaceContainer,
            unfocusedContainerColor = MaterialTheme.colorScheme.surfaceContainer,
            disabledContainerColor = MaterialTheme.colorScheme.surfaceContainer,
            focusedBorderColor = Color.Transparent,
            unfocusedBorderColor = Color.Transparent
        )

val TextStyle.short: TextStyle get() = copy(lineHeight = 20.sp)
val Typography.minTextSize: TextUnit get() = 10.sp
