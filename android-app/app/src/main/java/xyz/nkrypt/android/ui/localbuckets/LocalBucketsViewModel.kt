package xyz.nkrypt.android.ui.localbuckets

import android.content.Context
import android.content.Intent
import android.net.Uri
import android.provider.DocumentsContract
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.LocalBucketRepository
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import javax.inject.Inject

@HiltViewModel
class LocalBucketsViewModel @Inject constructor(
    private val repository: LocalBucketRepository,
    private val masterPasswordStore: MasterPasswordStore
) : ViewModel() {

    private val _buckets = MutableStateFlow<List<LocalBucketEntity>>(emptyList())
    val buckets: StateFlow<List<LocalBucketEntity>> = _buckets.asStateFlow()

    private val _createDialogState = MutableStateFlow<CreateBucketDialogState?>(null)
    val createDialogState: StateFlow<CreateBucketDialogState?> = _createDialogState.asStateFlow()

    private var pendingPath: String? = null

    init {
        viewModelScope.launch {
            repository.getAllBuckets().collect { list ->
                _buckets.value = list
            }
        }
    }

    fun showCreateDialog() {
        val masterPassword = masterPasswordStore.getMasterPassword()
        if (masterPassword == null) {
            return
        }
        _createDialogState.value = CreateBucketDialogState(selectedPath = null)
        pendingPath = null
    }

    fun onDirectorySelected(context: Context, uri: Uri) {
        context.contentResolver.takePersistableUriPermission(
            uri,
            Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_GRANT_WRITE_URI_PERMISSION
        )
        val path = getPathFromUri(uri)
        if (path.isNotBlank()) {
            pendingPath = path
            _createDialogState.value = CreateBucketDialogState(selectedPath = path)
        }
    }

    fun createBucket(name: String, encryptionPassword: String) {
        val path = pendingPath ?: return
        val masterPassword = masterPasswordStore.getMasterPassword() ?: return
        viewModelScope.launch {
            try {
                repository.createBucket(name, path, encryptionPassword, masterPassword)
                dismissCreateDialog()
            } catch (e: Exception) {
            }
        }
    }

    fun dismissCreateDialog() {
        _createDialogState.value = null
        pendingPath = null
    }

    fun deleteBucket(id: String, deleteFiles: Boolean) {
        viewModelScope.launch {
            try {
                repository.deleteBucket(id, deleteFiles)
            } catch (_: Exception) {}
        }
    }

    private fun getPathFromUri(uri: Uri): String {
        val docId = DocumentsContract.getTreeDocumentId(uri)
        val split = docId.split(":")
        return when {
            split.size >= 2 -> {
                val type = split[0]
                val path = split[1]
                when (type) {
                    "primary" -> "/storage/emulated/0/$path"
                    else -> "/storage/$type/$path"
                }
            }
            else -> uri.path ?: ""
        }
    }
}
