// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

package com.barghest.mesh.ui.theme

import androidx.compose.ui.graphics.Color

object LightColors {
    val Background = Color(0xFFFFFFFF)
    val Foreground = Color(0xFF262626)
    val Card = Color(0xFFFAFAFA)
    val CardForeground = Color(0xFF262626)
    val Popover = Color(0xFFFFFFFF)
    val PopoverForeground = Color(0xFF262626)
    val Primary = Color(0xFF4F46E5)
    val PrimaryForeground = Color(0xFFFFFFFF)
    val Secondary = Color(0xFFF3F4F6)
    val SecondaryForeground = Color(0xFF262626)
    val Muted = Color(0xFFF0F1F3)
    val MutedForeground = Color(0xFF808080)
    val Accent = Color(0xFFE5E7EB)
    val AccentForeground = Color(0xFF262626)
    val Destructive = Color(0xFFDC2626)
    val DestructiveForeground = Color(0xFFFAFAFA)
    val Border = Color(0xFFE5E7EB)
    val Input = Color(0xFFFAFAFA)
    val Ring = Color(0xFF4F46E5)
    val Sidebar = Color(0xFFFAFAFA)
    val SidebarForeground = Color(0xFF262626)
    val SidebarPrimary = Color(0xFF4F46E5)
    val SidebarPrimaryForeground = Color(0xFFFFFFFF)
    val SidebarAccent = Color(0xFFF3F4F6)
    val SidebarAccentForeground = Color(0xFF262626)
    val SidebarBorder = Color(0xFFE5E7EB)
    val SidebarRing = Color(0xFF4F46E5)

    val Chart1 = Color(0xFF5B5BFF)
    val Chart2 = Color(0xFF00B4D8)
    val Chart3 = Color(0xFF00D084)
    val Chart4 = Color(0xFFD4AF37)
    val Chart5 = Color(0xFFFF6B35)
}

// Dark warm-neutral palette with a terracotta signal accent.
object DarkColors {
    val Background = Color(0xFF171614) // warm charcoal base
    val Foreground = Color(0xFFF4F2ED) // warm off-white
    val Card = Color(0xFF232220)
    val CardForeground = Color(0xFFF4F2ED)
    val Popover = Color(0xFF1E1D1A) // raised
    val PopoverForeground = Color(0xFFF4F2ED)
    val Primary = Color(0xFFE08453) // terracotta signal accent
    val PrimaryForeground = Color(0xFFFFFFFF)
    val Secondary = Color(0xFF2A2926) // cardHi
    val SecondaryForeground = Color(0xFFF4F2ED)
    val Muted = Color(0xFF1E1D1A)
    val MutedForeground = Color(0xFF7C766C)
    val Accent = Color(0xFF2A2926)
    val AccentForeground = Color(0xFFF4F2ED)
    val Destructive = Color(0xFFFF6F61) // warm coral red
    val DestructiveForeground = Color(0xFFFFFFFF)
    val Border = Color(0xFF34322E)
    val Input = Color(0xFF232220)
    val Ring = Color(0xFFE08453)
    val Sidebar = Color(0xFF1E1D1A)
    val SidebarForeground = Color(0xFFF4F2ED)
    val SidebarPrimary = Color(0xFFE08453)
    val SidebarPrimaryForeground = Color(0xFFFFFFFF)
    val SidebarAccent = Color(0xFF2A2926)
    val SidebarAccentForeground = Color(0xFFF4F2ED)
    val SidebarBorder = Color(0xFF34322E)
    val SidebarRing = Color(0xFFE08453)

    val Chart1 = Color(0xFFE08453)
    val Chart2 = Color(0xFF54C6C4)
    val Chart3 = Color(0xFF4FD39B)
    val Chart4 = Color(0xFFF0B43F)
    val Chart5 = Color(0xFFE9786A)
}

/** MESH design tokens outside Material's ColorScheme slots; surfaced via MeshTokens.kt. */
object MeshPalette {
    val Bg = Color(0xFF171614)
    val BgRaised = Color(0xFF1E1D1A)
    val Card = Color(0xFF232220)
    val CardHi = Color(0xFF2A2926)
    val Border = Color(0xFF34322E)
    val BorderHi = Color(0xFF45433D)
    val Text = Color(0xFFF4F2ED)
    val Text2 = Color(0xFFACA79E)
    val Muted = Color(0xFF7C766C)
    val MonoText = Color(0xFFD3CFC6)

    val Accent = Color(0xFFE08453) // terracotta

    val Green = Color(0xFF4FD39B)
    val Cyan = Color(0xFF54C6C4)
    val Amber = Color(0xFFF0B43F)
    val Red = Color(0xFFFF6F61)
}

val ts_color_light_blue = Color(0xFF4B70CC)
