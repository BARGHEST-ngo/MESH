// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package com.barghest.mesh.ui.view.mesh

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.animateDpAsState
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.togetherWith
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.navigationBarsPadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.statusBarsPadding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.ColorFilter
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.barghest.mesh.R
import com.barghest.mesh.ui.view.onboarding.QuizDialog
import com.barghest.mesh.ui.view.onboarding.WarningDialog
import com.barghest.mesh.ui.theme.accent
import com.barghest.mesh.ui.theme.accentDim
import com.barghest.mesh.ui.theme.meshBg
import com.barghest.mesh.ui.theme.meshBorderHi
import com.barghest.mesh.ui.theme.meshGreen
import com.barghest.mesh.ui.theme.meshMuted
import com.barghest.mesh.ui.theme.meshText
import com.barghest.mesh.ui.theme.meshText2

@Composable
private fun OnbBullet(text: String) {
    val cs = MaterialTheme.colorScheme
    Row(
        Modifier.fillMaxWidth().padding(vertical = 9.dp),
        horizontalArrangement = Arrangement.spacedBy(12.dp),
        verticalAlignment = Alignment.Top,
    ) {
        Box(
            Modifier
                .padding(top = 1.dp)
                .size(22.dp)
                .background(cs.accentDim(0.16f), RoundedCornerShape(7.dp)),
            contentAlignment = Alignment.Center,
        ) { MeshIcon("check", size = 14.dp, color = cs.accent, stroke = 2.2f) }
        Text(text, color = cs.meshText2, fontSize = 14.5.sp, lineHeight = 21.sp)
    }
}

@Composable
fun MeshOnboarding(pendingLink: Boolean = false, onDone: () -> Unit) {
    val cs = MaterialTheme.colorScheme
    var page by remember { mutableIntStateOf(0) }
    val pages = 3

    // Comprehension pop-quizzes gating the "how it works" and consent pages, ported from
    // the legacy onboarding (true/false → feedback; a wrong answer shows a warning, then
    // proceeds either way).
    var activeQuiz by remember { mutableIntStateOf(0) }
    var showQuiz by remember { mutableStateOf(false) }
    var showWarning by remember { mutableStateOf(false) }
    var warningMessage by remember { mutableStateOf("") }

    val quizConsent = stringResource(R.string.mesh_quiz_consent)
    val quizAdb = stringResource(R.string.mesh_quiz_adb)
    val warnConsent = stringResource(R.string.mesh_warn_consent)
    val warnAdb = stringResource(R.string.mesh_warn_adb)

    fun proceedAfterQuiz() {
        if (activeQuiz == 2) onDone() else page = 2
    }

    if (showQuiz) {
        QuizDialog(
            title = stringResource(R.string.mesh_quiz_title),
            question = if (activeQuiz == 1) quizAdb else quizConsent,
            onAnswer = { answeredTrue ->
                showQuiz = false
                if (answeredTrue) {
                    proceedAfterQuiz()
                } else {
                    warningMessage = if (activeQuiz == 1) warnAdb else warnConsent
                    showWarning = true
                }
            },
        )
    }
    if (showWarning) {
        WarningDialog(message = warningMessage, onDismiss = { showWarning = false; proceedAfterQuiz() })
    }

    Column(
        Modifier.fillMaxSize().background(cs.meshBg).statusBarsPadding().navigationBarsPadding(),
    ) {
        Spacer(Modifier.height(14.dp))
        AnimatedContent(
            targetState = page,
            transitionSpec = { fadeIn(tween(350)) togetherWith fadeOut(tween(150)) },
            modifier = Modifier.weight(1f),
            label = "onbPage",
        ) { p ->
            when (p) {
                0 -> Welcome()
                1 -> HowItWorks()
                else -> Consent()
            }
        }

        Column(Modifier.padding(start = 26.dp, end = 26.dp, top = 16.dp, bottom = 26.dp)) {
            Row(
                Modifier.fillMaxWidth().padding(bottom = 18.dp),
                horizontalArrangement = Arrangement.spacedBy(7.dp, Alignment.CenterHorizontally),
            ) {
                repeat(pages) { i ->
                    val w by animateDpAsState(if (i == page) 24.dp else 5.dp, label = "ind")
                    Box(
                        Modifier
                            .height(5.dp)
                            .width(w)
                            .background(
                                if (i == page) cs.accent else cs.meshBorderHi,
                                RoundedCornerShape(3.dp),
                            ),
                    )
                }
            }
            if (pendingLink && page == pages - 1) {
                Text(
                    stringResource(R.string.mesh_onb_pending_notice),
                    color = cs.accent,
                    fontSize = 12.5.sp,
                    fontWeight = FontWeight.Medium,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.fillMaxWidth().padding(bottom = 12.dp),
                )
            }
            MeshButton(
                text = when (page) {
                    0 -> stringResource(R.string.mesh_onb_get_started)
                    pages - 1 -> stringResource(R.string.mesh_onb_understand)
                    else -> stringResource(R.string.mesh_onb_continue)
                },
                onClick = {
                    when (page) {
                        0 -> page = 1
                        1 -> { activeQuiz = 1; showQuiz = true }
                        else -> { activeQuiz = 2; showQuiz = true }
                    }
                },
                full = true,
                height = 56.dp,
                icon = if (page == pages - 1) "check" else "chevR",
            )
            if (page > 0) {
                Text(
                    stringResource(R.string.mesh_back),
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(top = 12.dp)
                        .clickable { page-- },
                    color = cs.meshMuted,
                    fontSize = 12.5.sp,
                    fontWeight = FontWeight.Medium,
                    textAlign = TextAlign.Center,
                )
            }
        }
    }
}

@Composable
private fun Welcome() {
    val cs = MaterialTheme.colorScheme
    Column(
        Modifier.fillMaxSize().padding(horizontal = 30.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        Image(
            painter = painterResource(R.drawable.mesh_wordmark),
            contentDescription = "MESH",
            colorFilter = ColorFilter.tint(cs.meshText),
            modifier = Modifier.fillMaxWidth(0.82f),
        )
        Spacer(Modifier.height(20.dp))
        Eyebrow(stringResource(R.string.mesh_onb_welcome_eyebrow))
        Spacer(Modifier.height(22.dp))
        Text(
            stringResource(R.string.mesh_onb_welcome_body),
            color = cs.meshText2, fontSize = 15.sp, lineHeight = 24.sp, textAlign = TextAlign.Center,
            modifier = Modifier.widthIn(max = 320.dp),
        )
        Spacer(Modifier.height(22.dp))
        Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(7.dp)) {
            StatusDot(cs.meshGreen, size = 7.dp)
            Text(stringResource(R.string.mesh_onb_welcome_tag), color = cs.meshMuted, fontSize = 11.5.sp, fontWeight = FontWeight.Medium)
        }
    }
}

/** Eyebrow + large title + body — the shared header for the How-it-works / Consent pages. */
@Composable
private fun OnbHeader(eyebrow: String, title: String, body: String) {
    val cs = MaterialTheme.colorScheme
    Eyebrow(eyebrow, color = cs.accent)
    Spacer(Modifier.height(10.dp))
    Text(title, color = cs.meshText, fontSize = 25.sp, fontWeight = FontWeight.SemiBold, lineHeight = 30.sp, letterSpacing = (-0.4).sp)
    Spacer(Modifier.height(8.dp))
    Text(body, color = cs.meshText2, fontSize = 14.5.sp, lineHeight = 22.sp)
    Spacer(Modifier.height(14.dp))
}

@Composable
private fun HowItWorks() {
    Column(Modifier.fillMaxSize().padding(start = 26.dp, end = 26.dp, top = 12.dp)) {
        OnbHeader(
            eyebrow = stringResource(R.string.mesh_onb_how_eyebrow),
            title = stringResource(R.string.mesh_onb_how_title),
            body = stringResource(R.string.mesh_onb_how_body),
        )
        OnbBullet(stringResource(R.string.mesh_onb_how_b1))
        OnbBullet(stringResource(R.string.mesh_onb_how_b2))
        OnbBullet(stringResource(R.string.mesh_onb_how_b3))
        OnbBullet(stringResource(R.string.mesh_onb_how_b4))
    }
}

@Composable
private fun Consent() {
    val cs = MaterialTheme.colorScheme
    Column(Modifier.fillMaxSize().padding(start = 26.dp, end = 26.dp, top = 12.dp)) {
        OnbHeader(
            eyebrow = stringResource(R.string.mesh_onb_consent_eyebrow),
            title = stringResource(R.string.mesh_onb_consent_title),
            body = stringResource(R.string.mesh_onb_consent_body),
        )
        OnbBullet(stringResource(R.string.mesh_onb_consent_b1))
        OnbBullet(stringResource(R.string.mesh_onb_consent_b2))
        OnbBullet(stringResource(R.string.mesh_onb_consent_b3))
        Spacer(Modifier.height(16.dp))
        Row(
            Modifier
                .fillMaxWidth()
                .background(cs.accentDim(0.10f), RoundedCornerShape(12.dp))
                .border(1.dp, cs.accentDim(0.30f), RoundedCornerShape(12.dp))
                .padding(14.dp),
            horizontalArrangement = Arrangement.spacedBy(11.dp),
            verticalAlignment = Alignment.Top,
        ) {
            MeshIcon("lock", size = 18.dp, color = cs.accent, modifier = Modifier.padding(top = 1.dp))
            Text(
                stringResource(R.string.mesh_onb_consent_note),
                color = cs.meshText2, fontSize = 13.sp, lineHeight = 20.sp,
            )
        }
    }
}
