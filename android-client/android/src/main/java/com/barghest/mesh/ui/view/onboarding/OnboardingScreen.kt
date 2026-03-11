// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.view.onboarding

import android.app.Activity
import android.graphics.Color as AndroidColor
import android.graphics.drawable.ColorDrawable
import androidx.activity.compose.BackHandler
import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.Animatable
import androidx.compose.animation.core.spring
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.navigationBarsPadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Check
import androidx.compose.material.icons.outlined.Close
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalView
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import androidx.compose.ui.window.DialogWindowProvider
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

private const val PAGE_COUNT = 3

private enum class QuizDialogPhase {
    Question,
    Feedback,
}

@Composable
fun OnboardingScreen(onComplete: () -> Unit) {
    val pagerState = rememberPagerState(pageCount = { PAGE_COUNT })
    val scope = rememberCoroutineScope()
    val context = LocalContext.current

    var showQuiz by rememberSaveable { mutableStateOf(false) }
    var showWarning by rememberSaveable { mutableStateOf(false) }
    var warningMessage by rememberSaveable { mutableStateOf("") }
    var activeQuiz by rememberSaveable { mutableStateOf(0) }

    BackHandler {
        if (pagerState.currentPage == 0) {
            (context as? Activity)?.finish()
        } else {
            scope.launch { pagerState.animateScrollToPage(pagerState.currentPage - 1) }
        }
    }

    if (showQuiz) {
        val (title, question) = when (activeQuiz) {
            1 -> "Pop quiz!" to
                    "Consensual forensics may involve granting access " +
                    "to sensitive data on your device."
            2 -> "Pop quiz!" to
                    "Leaving ADB enabled when not in use can increase " +
                    "security risks if your device is accessed by an untrusted party."
            else -> "" to ""
        }

        QuizDialog(
            title = title,
            question = question,
            onAnswer = { answeredTrue ->
                showQuiz = false
                if (answeredTrue) {
                    if (activeQuiz == 2) {
                        onComplete()
                    } else {
                        scope.launch {
                            pagerState.animateScrollToPage(pagerState.currentPage + 1)
                        }
                    }
                } else {
                    warningMessage = when (activeQuiz) {
                        1 -> "Consensual forensics can involve access to sensitive " +
                                "information. Only proceed with analysts you trust, " +
                                "and agree in advance on the scope of access."
                        2 -> "If ADB remains enabled, it may be misused by anyone " +
                                "with physical or network access to your device. " +
                                "Disable ADB immediately after a session."
                        else -> ""
                    }
                    showWarning = true
                }
            }
        )
    }

    if (showWarning) {
        WarningDialog(
            message = warningMessage,
            onDismiss = {
                showWarning = false
                if (activeQuiz == 2) {
                    onComplete()
                } else {
                    scope.launch {
                        pagerState.animateScrollToPage(pagerState.currentPage + 1)
                    }
                }
            }
        )
    }

    Scaffold(
        containerColor = Color.Black
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            HorizontalPager(
                state = pagerState,
                modifier = Modifier.weight(1f)
            ) { page ->
                when (page) {
                    0 -> WelcomePage()
                    1 -> HowItWorksPage()
                    2 -> ConsentPage()
                }
            }

            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .navigationBarsPadding()
                    .padding(horizontal = 24.dp, vertical = 16.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    repeat(PAGE_COUNT) { index ->
                        PageDot(isActive = pagerState.currentPage == index)
                    }
                }

                val isLastPage = pagerState.currentPage == PAGE_COUNT - 1

                Button(
                    onClick = {
                        when (pagerState.currentPage) {
                            1 -> {
                                activeQuiz = 1
                                showQuiz = true
                            }
                            2 -> {
                                activeQuiz = 2
                                showQuiz = true
                            }
                            else -> scope.launch {
                                pagerState.animateScrollToPage(pagerState.currentPage + 1)
                            }
                        }
                    },
                    colors = ButtonDefaults.buttonColors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainer,
                        contentColor = MaterialTheme.colorScheme.onSurface
                    )
                ) {
                    Text(if (isLastPage) "Begin Session" else "Next")
                }
            }
        }
    }
}

@Composable
private fun QuizDialog(
    title: String,
    question: String,
    onAnswer: (Boolean) -> Unit
) {
    var phase by rememberSaveable { mutableStateOf(QuizDialogPhase.Question) }
    var answerIsTrue by rememberSaveable { mutableStateOf(true) }
    var hasDispatchedResult by rememberSaveable { mutableStateOf(false) }

    val isCorrect = answerIsTrue

    val quizBlue = MaterialTheme.colorScheme.surfaceContainer
    val contentColor = MaterialTheme.colorScheme.onPrimary
    val quizGreen = Color(0xFF2F9E44)
    val quizRed = Color(0xFFE03131)

    val containerColor by animateColorAsState(
        targetValue = when (phase) {
            QuizDialogPhase.Question -> quizBlue
            QuizDialogPhase.Feedback -> if (isCorrect) quizGreen else quizRed
        },
        animationSpec = tween(180),
        label = "quizContainerColor"
    )

    val containerShape = when (phase) {
        QuizDialogPhase.Question -> RoundedCornerShape(16.dp)
        QuizDialogPhase.Feedback -> RoundedCornerShape(8.dp)
    }

    val view = LocalView.current
    SideEffect {
        val window = (view.parent as? DialogWindowProvider)?.window
        window?.setBackgroundDrawable(ColorDrawable(AndroidColor.TRANSPARENT))
        window?.setDimAmount(0.90f)
    }

    val feedbackScale = remember { Animatable(0.75f) }
    val feedbackAlpha = remember { Animatable(0f) }

    LaunchedEffect(phase) {
        if (phase != QuizDialogPhase.Feedback) return@LaunchedEffect
        feedbackScale.snapTo(0.75f)
        feedbackAlpha.snapTo(0f)
        feedbackAlpha.animateTo(1f, animationSpec = tween(120))
        feedbackScale.animateTo(
            1f,
            animationSpec = spring(dampingRatio = 0.55f, stiffness = 500f)
        )
    }

    LaunchedEffect(phase, hasDispatchedResult, answerIsTrue) {
        if (phase != QuizDialogPhase.Feedback || hasDispatchedResult) return@LaunchedEffect
        delay(800)
        hasDispatchedResult = true
        onAnswer(answerIsTrue)
    }

    fun submitAnswer(isTrue: Boolean) {
        if (phase != QuizDialogPhase.Question) return
        answerIsTrue = isTrue
        phase = QuizDialogPhase.Feedback
    }

    Dialog(
        onDismissRequest = {},
        properties = DialogProperties(
            dismissOnBackPress = false,
            dismissOnClickOutside = false,
            usePlatformDefaultWidth = false
        )
    ) {
        Surface(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 24.dp)
                .height(256.dp),
            color = containerColor,
            shape = containerShape
        ) {
            Box(modifier = Modifier.fillMaxSize()) {
                when (phase) {
                    QuizDialogPhase.Question -> {
                        Column(modifier = Modifier.fillMaxSize()) {
                            Text(
                                text = title,
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(8.dp),
                                style = MaterialTheme.typography.headlineSmall.copy(
                                    fontFamily = JetBrainsMonoFamily
                                ),
                                color = contentColor,
                                textAlign = TextAlign.Center
                            )

                            Box(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(1.dp)
                                    .background(contentColor)
                            )

                            Column(
                                modifier = Modifier
                                    .weight(1f)
                                    .padding(16.dp),
                                horizontalAlignment = Alignment.CenterHorizontally
                            ) {
                                Text(
                                    text = question,
                                    style = MaterialTheme.typography.bodyLarge.copy(
                                        fontFamily = JetBrainsMonoFamily
                                    ),
                                    color = contentColor,
                                    textAlign = TextAlign.Center
                                )
                            }

                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(48.dp)
                            ) {
                                Button(
                                    onClick = { submitAnswer(true) },
                                    modifier = Modifier.weight(1f),
                                    shape = RoundedCornerShape(bottomStart = 16.dp),
                                    colors = ButtonDefaults.buttonColors(
                                        containerColor = quizGreen,
                                        contentColor = contentColor
                                    ),
                                    contentPadding = PaddingValues(0.dp)
                                ) {
                                    Text("True", fontWeight = FontWeight.SemiBold, fontFamily = JetBrainsMonoFamily)
                                }

                                Button(
                                    onClick = { submitAnswer(false) },
                                    modifier = Modifier.weight(1f),
                                    shape = RoundedCornerShape(bottomEnd = 16.dp),
                                    colors = ButtonDefaults.buttonColors(
                                        containerColor = quizRed,
                                        contentColor = contentColor
                                    ),
                                    contentPadding = PaddingValues(0.dp)
                                ) {
                                    Text("False", fontWeight = FontWeight.SemiBold, fontFamily = JetBrainsMonoFamily)
                                }
                            }
                        }
                    }

                    QuizDialogPhase.Feedback -> {
                        Icon(
                            imageVector = if (isCorrect) Icons.Outlined.Check else Icons.Outlined.Close,
                            contentDescription = if (isCorrect) "Correct" else "Incorrect",
                            modifier = Modifier
                                .align(Alignment.Center)
                                .size(96.dp)
                                .graphicsLayer(
                                    alpha = feedbackAlpha.value,
                                    scaleX = feedbackScale.value,
                                    scaleY = feedbackScale.value
                                ),
                            tint = contentColor
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun WarningDialog(
    message: String,
    onDismiss: () -> Unit
) {
    Dialog(
        onDismissRequest = onDismiss,
        properties = DialogProperties(dismissOnBackPress = true, dismissOnClickOutside = true)
    ) {
        Card(
            shape = RoundedCornerShape(24.dp),
            colors = CardDefaults.cardColors(
                containerColor = Color.Black
            )
        ) {
            Column(
                modifier = Modifier.padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Icon(
                    imageVector = Icons.Outlined.Check,
                    contentDescription = null,
                    modifier = Modifier.size(36.dp),
                    tint = Color.White
                )
                Spacer(Modifier.height(12.dp))
                Text(
                    text = "Important",
                    style = MaterialTheme.typography.titleLarge.copy(
                        fontWeight = FontWeight.Bold,
                        fontFamily = JetBrainsMonoFamily
                    ),
                    color = Color.White,
                    textAlign = TextAlign.Center
                )
                Spacer(Modifier.height(12.dp))
                Text(
                    text = message,
                    style = MaterialTheme.typography.bodyLarge.copy(
                        fontFamily = JetBrainsMonoFamily
                    ),
                    color = Color.White,
                    textAlign = TextAlign.Center
                )
                Spacer(Modifier.height(20.dp))
                Button(
                    onClick = onDismiss,
                    shape = RoundedCornerShape(12.dp),
                    modifier = Modifier.fillMaxWidth(),
                    colors = ButtonDefaults.buttonColors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainer,
                        contentColor = Color.White
                    )
                ) {
                    Text(
                        "Understood",
                        fontWeight = FontWeight.SemiBold,
                        fontFamily = JetBrainsMonoFamily
                    )
                }
            }
        }
    }
}

@Composable
private fun PageDot(isActive: Boolean) {
    val color by animateColorAsState(
        targetValue = if (isActive)
            MaterialTheme.colorScheme.surfaceContainer
        else
            MaterialTheme.colorScheme.onSurface.copy(alpha = 0.3f),
        label = "dot"
    )

    Box(
        modifier = Modifier
            .size(8.dp)
            .clip(CircleShape)
            .background(color)
    )
}