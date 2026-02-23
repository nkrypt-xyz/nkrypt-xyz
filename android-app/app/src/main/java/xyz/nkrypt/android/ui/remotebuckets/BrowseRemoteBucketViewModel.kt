package xyz.nkrypt.android.ui.remotebuckets

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import xyz.nkrypt.android.data.remote.api.DirectoryDto
import xyz.nkrypt.android.data.remote.api.FileDto
import javax.inject.Inject

data class BrowseRemoteBucketState(
    val bucket: RemoteBucketEntity? = null,
    val directories: List<DirectoryDto> = emptyList(),
    val files: List<FileDto> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class BrowseRemoteBucketViewModel @Inject constructor(
    private val repository: RemoteBucketRepository,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val _state = MutableStateFlow(BrowseRemoteBucketState())
    val state: StateFlow<BrowseRemoteBucketState> = _state.asStateFlow()

    private val _showNewFolderDialog = MutableStateFlow(false)
    val showNewFolderDialog: StateFlow<Boolean> = _showNewFolderDialog.asStateFlow()

    private var bucketId: String? = savedStateHandle.get<String>("bucketId")
    private val directoryStack = mutableListOf<String?>()

    init {
        bucketId?.let { loadBucket(it) }
    }

    fun loadBucket(id: String) {
        if (bucketId == id && directoryStack.isNotEmpty()) return
        bucketId = id
        directoryStack.clear()
        directoryStack.add(null)
        loadContent()
    }

    private fun loadContent() {
        val id = bucketId ?: return
        viewModelScope.launch {
            _state.update { it.copy(isLoading = true, error = null) }
            try {
                val bucket = repository.getBucketById(id) ?: run {
                    _state.update { it.copy(isLoading = false, error = "Bucket not found") }
                    return@launch
                }
                val dirId = directoryStack.lastOrNull() ?: bucket.rootDirectoryId
                val apiKey = repository.getApiKey(bucket)
                val (dirs, files) = repository.getDirectory(bucket, dirId, apiKey)
                _state.update {
                    it.copy(
                        bucket = bucket,
                        directories = dirs,
                        files = files,
                        isLoading = false
                    )
                }
            } catch (e: Exception) {
                _state.update {
                    it.copy(
                        isLoading = false,
                        error = e.message ?: "Failed to load"
                    )
                }
            }
        }
    }

    fun navigateInto(directory: DirectoryDto) {
        directoryStack.add(directory._id)
        loadContent()
    }

    fun navigateUp() {
        if (directoryStack.size > 1) {
            directoryStack.removeLast()
            loadContent()
        }
    }

    fun canNavigateUp(): Boolean = directoryStack.size > 1

    fun showNewFolderDialog() {
        _showNewFolderDialog.value = true
    }

    fun dismissNewFolderDialog() {
        _showNewFolderDialog.value = false
    }

    suspend fun getCurrentDirectoryId(): String? {
        val id = bucketId ?: return null
        val bucket = repository.getBucketById(id) ?: return null
        return directoryStack.lastOrNull() ?: bucket.rootDirectoryId
    }

    suspend fun createDirectory(name: String) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        val parentId = getCurrentDirectoryId() ?: return
        repository.createDirectory(bucket, parentId, name, apiKey)
        dismissNewFolderDialog()
        loadContent()
    }

    suspend fun uploadFile(fileName: String, content: ByteArray) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        val dirId = getCurrentDirectoryId()
        repository.createFileAndUpload(bucket, dirId, fileName, content, apiKey)
        loadContent()
    }

    suspend fun deleteDirectory(directoryId: String) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        repository.deleteDirectory(bucket, directoryId, apiKey)
        loadContent()
    }

    suspend fun renameDirectory(directoryId: String, newName: String) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        repository.renameDirectory(bucket, directoryId, newName, apiKey)
        loadContent()
    }

    suspend fun deleteFile(fileId: String) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        repository.deleteFile(bucket, fileId, apiKey)
        loadContent()
    }

    suspend fun renameFile(fileId: String, newName: String) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val apiKey = repository.getApiKey(bucket)
        repository.renameFile(bucket, fileId, newName, apiKey)
        loadContent()
    }

    suspend fun downloadFileDecrypted(fileId: String): ByteArray? {
        val id = bucketId ?: return null
        val bucket = repository.getBucketById(id) ?: return null
        val apiKey = repository.getApiKey(bucket)
        return repository.downloadAndDecryptFile(bucket, fileId, apiKey)
    }
}
