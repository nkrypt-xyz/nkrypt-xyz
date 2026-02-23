package xyz.nkrypt.android.ui.rules

import android.content.Context
import android.content.Intent
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.LocalBucketRepository
import xyz.nkrypt.android.data.local.RulesRepository
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity
import xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity
import xyz.nkrypt.android.data.local.entity.ImportPostAction
import xyz.nkrypt.android.data.local.entity.SyncPostAction
import xyz.nkrypt.android.sync.SyncProgressHolder
import xyz.nkrypt.android.sync.SyncProgressState
import xyz.nkrypt.android.sync.SyncService
import javax.inject.Inject

@HiltViewModel
class RulesViewModel @Inject constructor(
    private val rulesRepository: RulesRepository,
    private val localBucketRepository: LocalBucketRepository,
    private val remoteBucketRepository: RemoteBucketRepository,
    private val progressHolder: SyncProgressHolder,
    @ApplicationContext private val context: Context
) : ViewModel() {

    private val _importRules = MutableStateFlow<List<AutoImportRuleEntity>>(emptyList())
    private val _syncRules = MutableStateFlow<List<AutoSyncRuleEntity>>(emptyList())
    private val _localBuckets = MutableStateFlow<List<xyz.nkrypt.android.data.local.entity.LocalBucketEntity>>(emptyList())
    private val _remoteBuckets = MutableStateFlow<List<xyz.nkrypt.android.data.local.entity.RemoteBucketEntity>>(emptyList())

    private val _showImportDialog = MutableStateFlow(false)
    private val _showSyncDialog = MutableStateFlow(false)
    private val _importDialogSourcePath = MutableStateFlow<String?>(null)
    private val _editImportRule = MutableStateFlow<AutoImportRuleEntity?>(null)
    private val _editSyncRule = MutableStateFlow<AutoSyncRuleEntity?>(null)

    val importRules: StateFlow<List<AutoImportRuleEntity>> = _importRules.asStateFlow()
    val syncRules: StateFlow<List<AutoSyncRuleEntity>> = _syncRules.asStateFlow()
    val localBuckets: StateFlow<List<xyz.nkrypt.android.data.local.entity.LocalBucketEntity>> = _localBuckets.asStateFlow()
    val remoteBuckets: StateFlow<List<xyz.nkrypt.android.data.local.entity.RemoteBucketEntity>> = _remoteBuckets.asStateFlow()
    val showImportDialog: StateFlow<Boolean> = _showImportDialog.asStateFlow()
    val showSyncDialog: StateFlow<Boolean> = _showSyncDialog.asStateFlow()
    val importDialogSourcePath: StateFlow<String?> = _importDialogSourcePath.asStateFlow()
    val editImportRule: StateFlow<AutoImportRuleEntity?> = _editImportRule.asStateFlow()
    val editSyncRule: StateFlow<AutoSyncRuleEntity?> = _editSyncRule.asStateFlow()

    private val _progressState = MutableStateFlow<ProgressState?>(null)
    val progressState: StateFlow<ProgressState?> = _progressState.asStateFlow()

    data class ProgressState(
        val message: String,
        val current: Int,
        val total: Int,
        val done: Boolean,
        val count: Int,
        val errors: Int,
        val isImport: Boolean
    )

    init {
        viewModelScope.launch {
            progressHolder.progress.collect { state ->
                _progressState.value = state?.let { s ->
                    ProgressState(s.message, s.current, s.total, s.done, s.count, s.errors, s.isImport)
                }
            }
        }
        viewModelScope.launch {
            rulesRepository.getAllImportRules().collect { _importRules.value = it }
        }
        viewModelScope.launch {
            rulesRepository.getAllSyncRules().collect { _syncRules.value = it }
        }
        viewModelScope.launch {
            localBucketRepository.getAllBuckets().collect { _localBuckets.value = it }
        }
        viewModelScope.launch {
            remoteBucketRepository.getAllBuckets().collect { _remoteBuckets.value = it }
        }
    }

    fun showCreateImportDialog() {
        _importDialogSourcePath.value = null
        _showImportDialog.value = true
    }

    fun showCreateSyncDialog() {
        _showSyncDialog.value = true
    }

    fun dismissImportDialog() {
        _showImportDialog.value = false
        _importDialogSourcePath.value = null
        _editImportRule.value = null
    }

    fun dismissSyncDialog() {
        _showSyncDialog.value = false
        _editSyncRule.value = null
    }

    fun showEditImportDialog(rule: AutoImportRuleEntity) {
        _editImportRule.value = rule
        _importDialogSourcePath.value = rule.sourceDirectoryPath
        _showImportDialog.value = true
    }

    fun showEditSyncDialog(rule: AutoSyncRuleEntity) {
        _editSyncRule.value = rule
        _showSyncDialog.value = true
    }

    fun onImportSourcePathSelected(path: String) {
        _importDialogSourcePath.value = path
    }

    fun createImportRule(
        name: String,
        sourcePath: String,
        targetBucketId: String,
        postAction: ImportPostAction
    ) {
        viewModelScope.launch {
            try {
                rulesRepository.createImportRule(name, sourcePath, targetBucketId, postAction)
                dismissImportDialog()
            } catch (_: Exception) {}
        }
    }

    fun createSyncRule(
        name: String,
        sourceBucketId: String,
        sourceDirectoryId: String?,
        targetRemoteBucketId: String,
        targetDirectoryId: String?,
        postAction: SyncPostAction
    ) {
        viewModelScope.launch {
            try {
                rulesRepository.createSyncRule(
                    name, sourceBucketId, sourceDirectoryId,
                    targetRemoteBucketId, targetDirectoryId, postAction
                )
                dismissSyncDialog()
            } catch (_: Exception) {}
        }
    }

    fun updateImportRule(
        id: String,
        name: String,
        sourcePath: String,
        targetBucketId: String,
        postAction: ImportPostAction
    ) {
        viewModelScope.launch {
            try {
                rulesRepository.updateImportRule(id, name, sourcePath, targetBucketId, postAction)
                dismissImportDialog()
            } catch (_: Exception) {}
        }
    }

    fun updateSyncRule(
        id: String,
        name: String,
        sourceBucketId: String,
        sourceDirectoryId: String?,
        targetRemoteBucketId: String,
        targetDirectoryId: String?,
        postAction: SyncPostAction
    ) {
        viewModelScope.launch {
            try {
                rulesRepository.updateSyncRule(
                    id, name, sourceBucketId, sourceDirectoryId,
                    targetRemoteBucketId, targetDirectoryId, postAction
                )
                dismissSyncDialog()
            } catch (_: Exception) {}
        }
    }

    fun deleteImportRule(id: String) {
        viewModelScope.launch {
            rulesRepository.deleteImportRule(id)
        }
    }

    fun deleteSyncRule(id: String) {
        viewModelScope.launch {
            rulesRepository.deleteSyncRule(id)
        }
    }

    fun runImport(ruleId: String, onComplete: () -> Unit) {
        _progressState.value = ProgressState("Scanning...", 0, 0, false, 0, 0, true)
        val intent = Intent(context, SyncService::class.java).apply {
            putExtra(SyncService.EXTRA_RULE_ID, ruleId)
            putExtra(SyncService.EXTRA_IS_IMPORT, true)
        }
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
            context.startForegroundService(intent)
        } else {
            context.startService(intent)
        }
    }

    fun runSync(ruleId: String, onComplete: () -> Unit) {
        _progressState.value = ProgressState("Scanning...", 0, 0, false, 0, 0, false)
        val intent = Intent(context, SyncService::class.java).apply {
            putExtra(SyncService.EXTRA_RULE_ID, ruleId)
            putExtra(SyncService.EXTRA_IS_IMPORT, false)
        }
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
            context.startForegroundService(intent)
        } else {
            context.startService(intent)
        }
    }

    fun clearProgress() {
        progressHolder.clear()
        _progressState.value = null
    }

    fun cancelSync() {
        progressHolder.requestCancel()
        _progressState.value = null
    }
}
