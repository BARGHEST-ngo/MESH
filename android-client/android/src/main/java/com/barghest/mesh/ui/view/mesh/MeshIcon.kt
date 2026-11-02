// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.size
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.CornerRadius
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.PathFillType
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.StrokeJoin
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.graphics.drawscope.scale
import androidx.compose.ui.graphics.vector.PathParser
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

/** Bespoke 24×24 stroke-icon set drawn via PathParser; [stroke] is in 24-unit terms. */
@Composable
fun MeshIcon(
    name: String,
    size: Dp = 22.dp,
    color: Color = Color.Unspecified,
    stroke: Float = 1.7f,
    modifier: Modifier = Modifier,
) {
    val tint = color.takeIf { it != Color.Unspecified } ?: Color(0xFFF4F2ED)
    Canvas(modifier = modifier.size(size)) {
        val unit = this.size.minDimension / 24f
        scale(unit, unit, pivot = Offset.Zero) {
            drawMeshIcon(name, tint, stroke)
        }
    }
}

// Parse each path once and reuse it; re-parsing on every draw was a source of jank.
private val pathCache = java.util.concurrent.ConcurrentHashMap<String, Path>()

private fun pp(d: String): Path = pathCache.getOrPut(d) { PathParser().parsePathString(d).toPath() }

/** Draws one named glyph in 24-unit coordinates. */
private fun DrawScope.drawMeshIcon(name: String, c: Color, w: Float) {
    val s = Stroke(width = w, cap = StrokeCap.Round, join = StrokeJoin.Round)

    fun path(d: String) = drawPath(pp(d), c, style = s)
    fun line(x1: Float, y1: Float, x2: Float, y2: Float) =
        drawLine(c, Offset(x1, y1), Offset(x2, y2), strokeWidth = w, cap = StrokeCap.Round)
    fun circle(cx: Float, cy: Float, r: Float, fill: Boolean = false) =
        if (fill) drawCircle(c, r, Offset(cx, cy))
        else drawCircle(c, r, Offset(cx, cy), style = s)
    fun rrect(x: Float, y: Float, ww: Float, hh: Float, rx: Float) =
        drawRoundRect(c, Offset(x, y), Size(ww, hh), CornerRadius(rx, rx), style = s)

    when (name) {
        "shield" -> path("M12 3l7 3v5c0 4.4-3 7.7-7 9-4-1.3-7-4.6-7-9V6l7-3z")
        "lock" -> { rrect(5f, 11f, 14f, 9f, 2f); path("M8 11V8a4 4 0 0 1 8 0v3") }
        "check" -> path("M5 12.5l4.5 4.5L19 7")
        "x" -> { line(6f, 6f, 18f, 18f); line(18f, 6f, 6f, 18f) }
        "chevL" -> path("M15 5l-7 7 7 7")
        "chevR" -> path("M9 5l7 7-7 7")
        "chevD" -> path("M5 9l7 7 7-7")
        "gear" -> {
            // Filled cog — a bare circle with radial ticks reads as a brightness toggle.
            val g = pp(
                "M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61" +
                    "l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54" +
                    "c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94" +
                    "l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94" +
                    "s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96" +
                    "c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54" +
                    "c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61" +
                    "l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z",
            )
            g.fillType = PathFillType.EvenOdd
            drawPath(g, c)
        }
        "phone" -> { rrect(6.5f, 2.5f, 11f, 19f, 2.5f); path("M10.5 18.5h3") }
        "laptop" -> { rrect(4f, 5f, 16f, 11f, 1.5f); path("M2.5 19.5h19") }
        "grid" -> {
            rrect(4f, 4f, 6.5f, 6.5f, 1.2f); rrect(13.5f, 4f, 6.5f, 6.5f, 1.2f)
            rrect(4f, 13.5f, 6.5f, 6.5f, 1.2f); rrect(13.5f, 13.5f, 6.5f, 6.5f, 1.2f)
        }
        "file" -> { path("M7 3h7l4 4v14H7z"); path("M14 3v4h4") }
        "scan" -> {
            path("M4 8V5.5A1.5 1.5 0 0 1 5.5 4H8M16 4h2.5A1.5 1.5 0 0 1 20 5.5V8M20 16v2.5a1.5 1.5 0 0 1-1.5 1.5H16M8 20H5.5A1.5 1.5 0 0 1 4 18.5V16")
            path("M4 12h16")
        }
        "link" -> {
            path("M9 15l6-6"); path("M11 7l1-1a4 4 0 0 1 6 6l-1 1"); path("M13 17l-1 1a4 4 0 0 1-6-6l1-1")
        }
        "activity", "pulse" -> path("M3 12h4l2.5 7 5-14L17 12h4")
        "qr" -> {
            rrect(4f, 4f, 6f, 6f, 1f); rrect(14f, 4f, 6f, 6f, 1f); rrect(4f, 14f, 6f, 6f, 1f)
            path("M14 14h3v3M20 14v6M14 20h3")
        }
        "key" -> { circle(8f, 8f, 4f); path("M11 11l8 8M16 16l2-2M18 18l2-2") }
        "eye" -> { path("M2.5 12S6 5.5 12 5.5 21.5 12 21.5 12 18 18.5 12 18.5 2.5 12 2.5 12z"); circle(12f, 12f, 3f) }
        "eyeOff" -> {
            path("M4 4l16 16")
            path("M9.5 9.6A3 3 0 0 0 12 15a3 3 0 0 0 2.4-1.2")
            path("M6.5 6.7C4 8.3 2.5 12 2.5 12s3.5 6.5 9.5 6.5c1.6 0 3-.4 4.2-1M9.8 5.8A8 8 0 0 1 12 5.5c6 0 9.5 6.5 9.5 6.5a16 16 0 0 1-2 2.7")
        }
        "alert" -> { path("M12 3l9.5 17h-19L12 3z"); path("M12 10v4.5M12 17.6v.01") }
        "power" -> { path("M12 3v9"); path("M6.5 7a8 8 0 1 0 11 0") }
        "wifi" -> { path("M2.5 9.5a14 14 0 0 1 19 0M5.5 13a9 9 0 0 1 13 0M8.5 16.5a4.5 4.5 0 0 1 7 0"); path("M12 20.5v.01") }
        "user" -> { circle(12f, 8f, 4f); path("M4.5 20a7.5 7.5 0 0 1 15 0") }
        "relay" -> {
            circle(12f, 12f, 2.5f); circle(4f, 6f, 2f); circle(20f, 6f, 2f); circle(12f, 20f, 2f)
            path("M10 11l-4.5-3.5M14 11l4.5-3.5M12 14.5v3.5")
        }
        "info" -> { circle(12f, 12f, 9f); path("M12 11v5M12 8v.01") }
        "bell" -> { path("M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6"); path("M10 19a2 2 0 0 0 4 0") }
        "globe" -> { circle(12f, 12f, 9f); path("M3 12h18M12 3c2.5 2.5 2.5 15 0 18M12 3c-2.5 2.5-2.5 15 0 18") }
        else -> Unit
    }
}
