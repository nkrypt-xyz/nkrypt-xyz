package xyz.nkrypt.android.sync

import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Shared progress state for sync/import operations running in foreground service.
 * Service updates; ViewModel observes for UI.
 */
@Singleton
class SyncProgressHolder @Inject constructor() {

    private val _progress = MutableStateFlow<SyncProgressState?>(null)
    val progress: StateFlow<SyncProgressState?> = _progress.asStateFlow()

    @Volatile
    var cancelRequested: Boolean = false
        private set

    fun update(state: SyncProgressState?) {
        _progress.value = state
    }

    fun requestCancel() {
        cancelRequested = true
    }

    fun clearCancelRequest() {
        cancelRequested = false
    }

    fun clear() {
        _progress.value = null
        cancelRequested = false
    }
}

data class SyncProgressState(
    val message: String,
    val current: Int,
    val total: Int,
    val done: Boolean,
    val count: Int,
    val errors: Int,
    val isImport: Boolean
)
