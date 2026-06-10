// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.ui.theme.MeshMono
import com.barghest.mesh.ui.theme.MeshRadii
import com.barghest.mesh.ui.theme.accent
import com.barghest.mesh.ui.theme.accentDim
import com.barghest.mesh.ui.theme.dim
import com.barghest.mesh.ui.theme.meshBg
import com.barghest.mesh.ui.theme.meshBgRaised
import com.barghest.mesh.ui.theme.meshBorder
import com.barghest.mesh.ui.theme.meshBorderHi
import com.barghest.mesh.ui.theme.meshCard
import com.barghest.mesh.ui.theme.meshCardHi
import com.barghest.mesh.ui.theme.meshMonoText
import com.barghest.mesh.ui.theme.meshMuted
import com.barghest.mesh.ui.theme.meshRed
import com.barghest.mesh.ui.theme.meshRedDim
import com.barghest.mesh.ui.theme.meshText
import com.barghest.mesh.ui.theme.meshText2

/** Status dot with an optional pulsing halo. */
@Composable
fun StatusDot(color: Color, size: Dp = 9.dp, pulse: Boolean = false) {
    Box(contentAlignment = Alignment.Center, modifier = Modifier.size(size)) {
        if (pulse) {
            val t = rememberInfiniteTransition(label = "dotPulse")
            val scale by t.animateFloat(
                initialValue = 1f, targetValue = 2.6f,
                animationSpec = infiniteRepeatable(tween(1800), RepeatMode.Restart), label = "scale",
            )
            val fade by t.animateFloat(
                initialValue = 0.5f, targetValue = 0f,
                animationSpec = infiniteRepeatable(tween(1800), RepeatMode.Restart), label = "fade",
            )
            Box(
                Modifier
                    .size(size)
                    .graphicsLayer { scaleX = scale; scaleY = scale; alpha = fade }
                    .clip(CircleShape)
                    .background(color),
            )
        }
        Box(Modifier.size(size).clip(CircleShape).background(color))
    }
}

@Composable
fun Eyebrow(text: String, modifier: Modifier = Modifier, color: Color = MaterialTheme.colorScheme.meshMuted) {
    Text(
        text = text,
        modifier = modifier,
        fontFamily = MaterialTheme.typography.bodyMedium.fontFamily,
        fontSize = 12.sp,
        fontWeight = FontWeight.SemiBold,
        letterSpacing = 0.2.sp,
        color = color,
    )
}

/** Shared card surface: rounded clip + fill + hairline border. */
fun Modifier.meshCard(
    shape: androidx.compose.ui.graphics.Shape,
    fill: Color,
    border: Color,
): Modifier = clip(shape).background(fill).border(1.dp, border, shape)

@Composable
fun MeshCard(
    modifier: Modifier = Modifier,
    onClick: (() -> Unit)? = null,
    pad: Dp = 16.dp,
    border: Color = MaterialTheme.colorScheme.meshBorder,
    content: @Composable () -> Unit,
) {
    val shape = RoundedCornerShape(MeshRadii.card)
    var base = modifier.meshCard(shape, MaterialTheme.colorScheme.meshCard, border)
    if (onClick != null) base = base.clickable(onClick = onClick)
    Box(base.padding(pad)) { content() }
}

/** Standard list row: 34dp icon square + title/sub column + trailing slot (default chevron). */
@Composable
fun MeshRow(
    icon: String,
    title: String,
    sub: String? = null,
    modifier: Modifier = Modifier,
    iconTint: Color = MaterialTheme.colorScheme.accent,
    iconBg: Color = MaterialTheme.colorScheme.accentDim(0.14f),
    titleColor: Color = MaterialTheme.colorScheme.meshText,
    spacing: Dp = 12.dp,
    trailing: @Composable () -> Unit = {
        MeshIcon("chevR", size = 18.dp, color = MaterialTheme.colorScheme.meshMuted)
    },
) {
    val cs = MaterialTheme.colorScheme
    Row(
        modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(spacing),
    ) {
        Box(
            Modifier.size(34.dp).clip(RoundedCornerShape(10.dp)).background(iconBg),
            contentAlignment = Alignment.Center,
        ) { MeshIcon(icon, size = 18.dp, color = iconTint) }
        Column(Modifier.weight(1f)) {
            Text(title, color = titleColor, fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
            if (sub != null) Text(sub, color = cs.meshText2, fontSize = 12.sp)
        }
        trailing()
    }
}

/** Full-width hairline in the border colour. */
@Composable
fun MeshDivider() {
    Box(Modifier.fillMaxWidth().height(1.dp).background(MaterialTheme.colorScheme.meshBorder))
}

/** Bottom action bar: a hairline divider above a padded slot for the screen's primary button(s). */
@Composable
fun MeshBottomBar(content: @Composable () -> Unit) {
    Column(Modifier.fillMaxWidth().background(MaterialTheme.colorScheme.meshBg)) {
        MeshDivider()
        Box(Modifier.padding(horizontal = 20.dp, vertical = 14.dp)) { content() }
    }
}

enum class MeshButtonVariant { Primary, Secondary, Ghost, Danger }

@Composable
fun MeshButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    variant: MeshButtonVariant = MeshButtonVariant.Primary,
    full: Boolean = false,
    icon: String? = null,
    height: Dp = 52.dp,
    enabled: Boolean = true,
) {
    val cs = MaterialTheme.colorScheme
    val shape = RoundedCornerShape(MeshRadii.md)
    val (bg, fg, borderColor) = when (variant) {
        MeshButtonVariant.Primary -> Triple(cs.accent, Color.White, Color.Transparent)
        MeshButtonVariant.Secondary -> Triple(cs.meshCard, cs.meshText, cs.meshBorderHi)
        MeshButtonVariant.Ghost -> Triple(Color.Transparent, cs.meshText2, cs.meshBorder)
        MeshButtonVariant.Danger -> Triple(cs.meshRedDim, cs.meshRed, cs.meshRed)
    }
    var m = modifier
        .then(if (full) Modifier.fillMaxWidth() else Modifier)
        .height(height)
        .clip(shape)
        .background(bg)
        .border(1.dp, borderColor, shape)
        .alpha(if (enabled) 1f else 0.45f)
    if (enabled) m = m.clickable(onClick = onClick)
    Row(
        modifier = m.padding(horizontal = 20.dp),
        horizontalArrangement = Arrangement.Center,
        verticalAlignment = Alignment.CenterVertically,
    ) {
        if (icon != null) {
            MeshIcon(icon, size = 18.dp, color = fg)
            Spacer(Modifier.width(9.dp))
        }
        Text(text, color = fg, fontSize = 15.5.sp, fontWeight = FontWeight.SemiBold)
    }
}

@Composable
fun MeshChip(
    text: String,
    modifier: Modifier = Modifier,
    color: Color = MaterialTheme.colorScheme.meshText2,
    bg: Color = MaterialTheme.colorScheme.meshBgRaised,
    border: Color = MaterialTheme.colorScheme.meshBorder,
    icon: String? = null,
) {
    Row(
        modifier = modifier
            .clip(RoundedCornerShape(MeshRadii.pill))
            .background(bg)
            .border(1.dp, border, RoundedCornerShape(MeshRadii.pill))
            .padding(horizontal = 9.dp, vertical = 4.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(5.dp),
    ) {
        if (icon != null) MeshIcon(icon, size = 12.dp, color = color)
        Text(text, color = color, fontSize = 11.sp, fontWeight = FontWeight.SemiBold, letterSpacing = 0.2.sp)
    }
}

/** Eyebrow + title + sub; the eyebrow is hidden when it would just echo the title. */
@Composable
fun SectionHead(
    title: String? = null,
    eyebrow: String? = null,
    sub: String? = null,
    action: (@Composable () -> Unit)? = null,
) {
    val cs = MaterialTheme.colorScheme
    val showEyebrow = eyebrow != null &&
        (title == null || !eyebrow.equals(title, ignoreCase = true))
    Row(
        modifier = Modifier.fillMaxWidth().padding(bottom = 12.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.Top,
    ) {
        Column(Modifier.weight(1f)) {
            if (showEyebrow) {
                Eyebrow(eyebrow!!)
                Spacer(Modifier.height(6.dp))
            }
            if (title != null) {
                Text(title, color = cs.meshText, fontSize = 16.sp, fontWeight = FontWeight.SemiBold, letterSpacing = (-0.1).sp)
            }
            if (sub != null) {
                Spacer(Modifier.height(3.dp))
                Text(sub, color = cs.meshText2, fontSize = 13.sp, lineHeight = 18.sp)
            }
        }
        if (action != null) {
            Spacer(Modifier.width(12.dp))
            action()
        }
    }
}

@Composable
fun MeshTopBar(
    label: String,
    onBack: (() -> Unit)? = null,
    action: (@Composable () -> Unit)? = null,
) {
    val cs = MaterialTheme.colorScheme
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .background(cs.background)
            .padding(horizontal = 14.dp, vertical = 10.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        if (onBack != null) {
            Box(
                Modifier
                    .size(38.dp)
                    .clip(RoundedCornerShape(12.dp))
                    .background(cs.meshCard)
                    .border(1.dp, cs.meshBorder, RoundedCornerShape(12.dp))
                    .clickable(onClick = onBack),
                contentAlignment = Alignment.Center,
            ) { MeshIcon("chevL", size = 19.dp, color = cs.meshText2) }
        } else {
            Spacer(Modifier.width(6.dp))
        }
        Text(
            label,
            modifier = Modifier.weight(1f),
            color = cs.meshText,
            fontSize = 16.5.sp,
            fontWeight = FontWeight.SemiBold,
            letterSpacing = (-0.2).sp,
        )
        if (action != null) action()
    }
}

/** Monospace technical value (IPs, byte counts, fingerprints). */
@Composable
fun MonoText(
    text: String,
    color: Color = MaterialTheme.colorScheme.meshMonoText,
    fontSize: androidx.compose.ui.unit.TextUnit = 12.sp,
    fontWeight: FontWeight = FontWeight.Normal,
    modifier: Modifier = Modifier,
) {
    Text(text, modifier = modifier, color = color, fontFamily = MeshMono, fontSize = fontSize, fontWeight = fontWeight)
}
