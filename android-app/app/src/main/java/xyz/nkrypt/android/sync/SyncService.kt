package xyz.nkrypt.android.sync

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.launch
import xyz.nkrypt.android.MainActivity
import xyz.nkrypt.android.data.local.ImportProgress
import xyz.nkrypt.android.data.local.RulesRepository
import xyz.nkrypt.android.data.local.SyncProgress
import javax.inject.Inject

private const val CHANNEL_ID = "sync_channel"
private const val NOTIFICATION_ID = 1001

@AndroidEntryPoint
class SyncService : Service() {

    companion object {
        const val EXTRA_RULE_ID = "rule_id"
        const val EXTRA_IS_IMPORT = "is_import"
    }

    @Inject
    lateinit var rulesRepository: RulesRepository

    @Inject
    lateinit var progressHolder: SyncProgressHolder

    private val serviceScope = CoroutineScope(SupervisorJob() + Dispatchers.Main)
    private var syncJob: Job? = null

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val ruleId = intent?.getStringExtra(EXTRA_RULE_ID) ?: return START_NOT_STICKY
        val isImport = intent.getBooleanExtra(EXTRA_IS_IMPORT, true)

        progressHolder.clearCancelRequest()
        syncJob?.cancel()
        syncJob = serviceScope.launch {
            if (isImport) {
                runImport(ruleId)
            } else {
                runSync(ruleId)
            }
            stopForeground(STOP_FOREGROUND_REMOVE)
            stopSelf(startId)
        }
        return START_NOT_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        syncJob?.cancel()
        serviceScope.cancel()
        // Don't clear progress here - let the user see the completed state on ProgressScreen.
        // Progress is cleared when user taps "Close" or when starting a new run.
    }

    override fun onBind(intent: Intent?): IBinder? = null

    private suspend fun runImport(ruleId: String) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE) {
            startForeground(NOTIFICATION_ID, buildNotification("Importing...", 0, 0, true), android.content.pm.ServiceInfo.FOREGROUND_SERVICE_TYPE_DATA_SYNC)
        } else {
            @Suppress("DEPRECATION")
            startForeground(NOTIFICATION_ID, buildNotification("Importing...", 0, 0, true))
        }
        rulesRepository.executeImport(ruleId)
            .catch { e ->
                progressHolder.update(
                    SyncProgressState(e.message ?: "Error", 0, 0, true, 0, 0, true)
                )
            }
            .collect { progress ->
                if (progressHolder.cancelRequested) throw CancellationException("Cancelled by user")
                when (progress) {
                    is ImportProgress.Scanning ->
                        progressHolder.update(SyncProgressState("Scanning...", 0, 0, false, 0, 0, true))
                    is ImportProgress.Progress -> {
                        progressHolder.update(
                            SyncProgressState(
                                "${progress.fileName}: ${progress.phase}",
                                progress.current, progress.total, false, 0, 0, true
                            )
                        )
                        updateNotification("Importing: ${progress.current}/${progress.total}", progress.current, progress.total, true)
                    }
                    is ImportProgress.Done -> {
                        progressHolder.update(
                            SyncProgressState("Done", 0, 0, true, progress.imported, progress.errors, true)
                        )
                        updateNotification("Import done: ${progress.imported} files", progress.imported, progress.imported, true)
                    }
                    is ImportProgress.Error ->
                        progressHolder.update(SyncProgressState(progress.message, 0, 0, true, 0, 0, true))
                }
            }
    }

    private suspend fun runSync(ruleId: String) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE) {
            startForeground(NOTIFICATION_ID, buildNotification("Syncing...", 0, 0, false), android.content.pm.ServiceInfo.FOREGROUND_SERVICE_TYPE_DATA_SYNC)
        } else {
            @Suppress("DEPRECATION")
            startForeground(NOTIFICATION_ID, buildNotification("Syncing...", 0, 0, false))
        }
        rulesRepository.executeSync(ruleId)
            .catch { e ->
                progressHolder.update(
                    SyncProgressState(e.message ?: "Error", 0, 0, true, 0, 0, false)
                )
            }
            .collect { progress ->
                if (progressHolder.cancelRequested) throw CancellationException("Cancelled by user")
                when (progress) {
                    is SyncProgress.Scanning ->
                        progressHolder.update(SyncProgressState("Scanning...", 0, 0, false, 0, 0, false))
                    is SyncProgress.Progress -> {
                        progressHolder.update(
                            SyncProgressState(
                                "${progress.fileName}: ${progress.phase}",
                                progress.current, progress.total, false, 0, 0, false
                            )
                        )
                        updateNotification("Syncing: ${progress.current}/${progress.total}", progress.current, progress.total, false)
                    }
                    is SyncProgress.Done -> {
                        progressHolder.update(
                            SyncProgressState("Done", 0, 0, true, progress.synced, progress.errors, false)
                        )
                        updateNotification("Sync done: ${progress.synced} files", progress.synced, progress.synced, false)
                    }
                    is SyncProgress.Error ->
                        progressHolder.update(SyncProgressState(progress.message, 0, 0, true, 0, 0, false))
                }
            }
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "Sync Progress",
                NotificationManager.IMPORTANCE_LOW
            ).apply { setShowBadge(false) }
            getSystemService(NotificationManager::class.java).createNotificationChannel(channel)
        }
    }

    private fun buildNotification(title: String, current: Int, total: Int, isImport: Boolean): android.app.Notification {
        val pendingIntent = PendingIntent.getActivity(
            this, 0,
            Intent(this, MainActivity::class.java).apply { flags = Intent.FLAG_ACTIVITY_SINGLE_TOP },
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )
        val progressText = if (total > 0) "$current / $total" else title
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle(if (isImport) "Importing" else "Syncing")
            .setContentText(progressText)
            .setSmallIcon(android.R.drawable.ic_menu_upload)
            .setContentIntent(pendingIntent)
            .setProgress(total.coerceAtLeast(1), current.coerceAtMost(total), total == 0)
            .setOngoing(true)
            .build()
    }

    private fun updateNotification(title: String, current: Int, total: Int, isImport: Boolean) {
        val notification = buildNotification(title, current, total, isImport)
        getSystemService(NotificationManager::class.java).notify(NOTIFICATION_ID, notification)
    }
}
