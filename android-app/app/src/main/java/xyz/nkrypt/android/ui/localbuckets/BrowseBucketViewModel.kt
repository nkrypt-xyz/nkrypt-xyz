package xyz.nkrypt.android.ui.localbuckets

import android.content.Context
import android.content.Intent
import android.net.Uri
import androidx.documentfile.provider.DocumentFile
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import xyz.nkrypt.android.data.local.LocalBucketRepository
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity
import xyz.nkrypt.android.util.uniqueFileName
import javax.inject.Inject

sealed class DownloadTarget {
    data class File(val file: LocalFileEntity) : DownloadTarget()
    data class Directory(val dir: LocalDirectoryEntity) : DownloadTarget()
}

data class BrowseBucketState(
    val bucket: LocalBucketEntity? = null,
    val directories: List<LocalDirectoryEntity> = emptyList(),
    val files: List<LocalFileEntity> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class BrowseBucketViewModel @Inject constructor(
    private val repository: LocalBucketRepository,
    private val masterPasswordStore: MasterPasswordStore,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val _state = MutableStateFlow(BrowseBucketState())
    val state: StateFlow<BrowseBucketState> = _state.asStateFlow()

    private val _showNewFolderDialog = MutableStateFlow(false)
    val showNewFolderDialog: StateFlow<Boolean> = _showNewFolderDialog.asStateFlow()

    private val _createFolderError = MutableStateFlow<String?>(null)
    val createFolderError: StateFlow<String?> = _createFolderError.asStateFlow()

    private val _renameError = MutableStateFlow<String?>(null)
    val renameError: StateFlow<String?> = _renameError.asStateFlow()

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
        val dirId = directoryStack.lastOrNull()
        viewModelScope.launch {
            _state.update { it.copy(isLoading = true, error = null) }
            try {
                val bucket = repository.getBucketById(id)
                if (bucket == null) {
                    _state.update {
                        it.copy(isLoading = false, error = "Bucket not found")
                    }
                    return@launch
                }
                val masterPassword = masterPasswordStore.getMasterPassword()
                if (masterPassword == null) {
                    _state.update {
                        it.copy(isLoading = false, error = "Master password required")
                    }
                    return@launch
                }
                val dirs = repository.getDirectories(id, dirId)
                val files = repository.getFiles(id, dirId)
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

    fun navigateInto(directory: LocalDirectoryEntity) {
        directoryStack.add(directory.id)
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
        _createFolderError.value = null
    }

    fun clearRenameError() {
        _renameError.value = null
    }

    suspend fun createFolder(name: String) {
        val id = bucketId ?: return
        val dirId = directoryStack.lastOrNull()
        try {
            repository.createDirectory(id, dirId, name)
            dismissNewFolderDialog()
            loadContent()
        } catch (e: IllegalArgumentException) {
            _createFolderError.value = e.message
        }
    }

    suspend fun uploadFile(fileName: String, content: ByteArray) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val masterPassword = masterPasswordStore.getMasterPassword() ?: return
        val bucketPassword = repository.decryptBucketPassword(bucket.cryptData, masterPassword)
        val dirId = directoryStack.lastOrNull()
        try {
            repository.createFileWithContent(bucket, dirId, fileName, content, bucketPassword)
            loadContent()
        } catch (e: IllegalArgumentException) {
            _state.update { it.copy(error = e.message) }
        }
    }

    suspend fun deleteFile(fileId: String) {
        val id = bucketId ?: return
        repository.deleteFile(id, fileId)
        loadContent()
    }

    suspend fun renameFile(fileId: String, newName: String): Boolean {
        val id = bucketId ?: return false
        return try {
            repository.renameFile(id, fileId, newName)
            _renameError.value = null
            loadContent()
            true
        } catch (e: IllegalArgumentException) {
            _renameError.value = e.message
            false
        }
    }

    suspend fun moveFile(fileId: String, newDirectoryId: String?) {
        val id = bucketId ?: return
        repository.moveFile(id, fileId, newDirectoryId)
        loadContent()
    }

    suspend fun deleteDirectory(directoryId: String) {
        val id = bucketId ?: return
        repository.deleteDirectory(id, directoryId)
        loadContent()
    }

    suspend fun renameDirectory(directoryId: String, newName: String): Boolean {
        val id = bucketId ?: return false
        return try {
            repository.renameDirectory(id, directoryId, newName)
            _renameError.value = null
            loadContent()
            true
        } catch (e: IllegalArgumentException) {
            _renameError.value = e.message
            false
        }
    }

    suspend fun moveDirectory(directoryId: String, newParentId: String?) {
        val id = bucketId ?: return
        repository.moveDirectory(id, directoryId, newParentId)
        loadContent()
    }

    suspend fun downloadToDirectory(
        target: DownloadTarget,
        destTreeUri: Uri,
        context: Context
    ) {
        val id = bucketId ?: return
        val bucket = repository.getBucketById(id) ?: return
        val masterPassword = masterPasswordStore.getMasterPassword() ?: return
        val bucketPassword = repository.decryptBucketPassword(bucket.cryptData, masterPassword)
        context.contentResolver.takePersistableUriPermission(
            destTreeUri,
            Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_GRANT_WRITE_URI_PERMISSION
        )
        withContext(Dispatchers.IO) {
            when (target) {
                is DownloadTarget.File -> downloadFile(
                    repository, bucket, bucketPassword, target.file,
                    destTreeUri, context
                )
                is DownloadTarget.Directory -> downloadDirectoryRecursive(
                    repository, bucket, bucketPassword, target.dir,
                    destTreeUri, context
                )
            }
        }
        loadContent()
    }

    private suspend fun downloadFile(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        file: LocalFileEntity,
        destTreeUri: Uri,
        context: Context
    ) {
        val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: return
        val root = DocumentFile.fromTreeUri(context, destTreeUri) ?: return
        val name = uniqueFileName(root, file.name)
        val docFile = root.createFile("application/octet-stream", name) ?: return
        context.contentResolver.openOutputStream(docFile.uri)?.use { it.write(content) }
    }

    private suspend fun downloadDirectoryRecursive(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        dir: LocalDirectoryEntity,
        destTreeUri: Uri,
        context: Context
    ) {
        val root = DocumentFile.fromTreeUri(context, destTreeUri) ?: return
        val targetDir = root.createDirectory(dir.name) ?: return
        val files = repo.getFilesRecursive(bucket.id, dir.id)
        for (file in files) {
            val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: continue
            val relPath = repo.getDirectoryPath(bucket.id, file.directoryId ?: dir.id, dir.id)
            val parentDir = if (relPath.isNullOrBlank()) {
                targetDir
            } else {
                findOrCreateDirectory(targetDir, relPath)
            }
            val name = parentDir?.let { uniqueFileName(it, file.name) } ?: file.name
            val docFile = parentDir?.createFile("application/octet-stream", name) ?: continue
            context.contentResolver.openOutputStream(docFile.uri)?.use { it.write(content) }
        }
    }

    private fun findOrCreateDirectory(parent: DocumentFile, path: String): DocumentFile? {
        val parts = path.split("/").filter { it.isNotBlank() }
        var current = parent
        for (name in parts) {
            var child = current.findFile(name)
            if (child == null || !child.isDirectory) {
                child = current.createDirectory(name)
            }
            current = child ?: return null
        }
        return current
    }
}
