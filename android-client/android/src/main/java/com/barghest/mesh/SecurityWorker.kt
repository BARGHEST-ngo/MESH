package com.barghest.mesh

import android.Manifest
import android.os.Build
import android.content.Context
import android.content.pm.PackageManager
import androidx.core.app.ActivityCompat
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters
import androidx.work.ListenableWorker.Result
import com.barghest.mesh.util.SecurityChecker
import com.barghest.mesh.R


class SecWorker(appContext: Context, workerParams: WorkerParameters) :
    CoroutineWorker(appContext, workerParams) {
    override suspend fun doWork(): Result {
        if (SecurityChecker.isDevOptionsEnabled(applicationContext)) {
            showDevNotification(applicationContext)
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            if (SecurityChecker.isSPLstale(31)) {
                showSPLNotification(applicationContext)
            }
        }
	return Result.success()
    }
}


fun showDevNotification(context: Context) {
    val notification = NotificationCompat.Builder(context, UninitializedApp.CHANNEL_ID)
        .setSmallIcon(R.drawable.ic_notification)
        .setContentTitle("MESH: Developer Options enabled")
        .setContentText("If no forensics session is occuring, please disable for your own security")
        .setPriority(NotificationCompat.PRIORITY_DEFAULT)
        .setAutoCancel(true)
        .build()

    if (ActivityCompat.checkSelfPermission(
            context,
            Manifest.permission.POST_NOTIFICATIONS
        ) != PackageManager.PERMISSION_GRANTED
    ) {
        return
    }

    NotificationManagerCompat.from(context).notify(1001, notification)
}

fun showSPLNotification(context: Context) {
    val notification = NotificationCompat.Builder(context, UninitializedApp.CHANNEL_ID)
        .setSmallIcon(R.drawable.ic_notification)
        .setContentTitle("MESH: Security Patch Outdated")
        .setContentText("Your security patch is out of date. Please update your device.")
        .setPriority(NotificationCompat.PRIORITY_DEFAULT)
        .setAutoCancel(true)
        .build()

    if (ActivityCompat.checkSelfPermission(
            context,
            Manifest.permission.POST_NOTIFICATIONS
        ) != PackageManager.PERMISSION_GRANTED
    ) {
        return
    }

    NotificationManagerCompat.from(context).notify(1002, notification)
}
