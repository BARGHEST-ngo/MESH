// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.theme

import androidx.compose.material3.ColorScheme
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

/** MESH design tokens (accent, semantic colours, neutral steps) as [ColorScheme] extensions. */

val ColorScheme.accent: Color get() = MeshPalette.Accent

fun ColorScheme.accentDim(alpha: Float = 0.16f): Color = MeshPalette.Accent.copy(alpha = alpha)

val ColorScheme.meshGreen: Color get() = MeshPalette.Green
val ColorScheme.meshCyan: Color get() = MeshPalette.Cyan
val ColorScheme.meshAmber: Color get() = MeshPalette.Amber
val ColorScheme.meshRed: Color get() = MeshPalette.Red

val ColorScheme.meshGreenDim: Color get() = MeshPalette.Green.copy(alpha = 0.15f)
val ColorScheme.meshCyanDim: Color get() = MeshPalette.Cyan.copy(alpha = 0.15f)
val ColorScheme.meshAmberDim: Color get() = MeshPalette.Amber.copy(alpha = 0.15f)
val ColorScheme.meshRedDim: Color get() = MeshPalette.Red.copy(alpha = 0.15f)

val ColorScheme.meshBg: Color get() = MeshPalette.Bg
val ColorScheme.meshBgRaised: Color get() = MeshPalette.BgRaised
val ColorScheme.meshCard: Color get() = MeshPalette.Card
val ColorScheme.meshCardHi: Color get() = MeshPalette.CardHi
val ColorScheme.meshBorder: Color get() = MeshPalette.Border
val ColorScheme.meshBorderHi: Color get() = MeshPalette.BorderHi
val ColorScheme.meshText: Color get() = MeshPalette.Text
val ColorScheme.meshText2: Color get() = MeshPalette.Text2
val ColorScheme.meshMuted: Color get() = MeshPalette.Muted
val ColorScheme.meshMonoText: Color get() = MeshPalette.MonoText

fun Color.dim(alpha: Float = 0.15f): Color = copy(alpha = alpha)

object MeshRadii {
    val sm = 12.dp
    val md = 16.dp
    val card = 20.dp
    val pill = 999.dp
}
