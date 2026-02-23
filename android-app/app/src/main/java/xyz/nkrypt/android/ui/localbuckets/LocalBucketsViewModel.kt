package xyz.nkrypt.android.ui.localbuckets

import android.content.Context
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
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
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity
import xyz.nkrypt.android.util.uniqueFileNameForFile
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

    fun onDirectorySelectedPath(path: String) {
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

    suspend fun downloadBucketToDirectory(bucket: LocalBucketEntity, destPath: String) {
        val masterPassword = masterPasswordStore.getMasterPassword() ?: return
        val bucketPassword = repository.decryptBucketPassword(bucket.cryptData, masterPassword)
        withContext(Dispatchers.IO) {
            val rootDir = java.io.File(destPath)
            rootDir.mkdirs()
            val rootDirs = repository.getDirectories(bucket.id, null)
            val rootFiles = repository.getFiles(bucket.id, null)
            for (dir in rootDirs) {
                downloadDirectoryRecursive(repository, bucket, bucketPassword, dir, rootDir)
            }
            for (file in rootFiles) {
                downloadFile(repository, bucket, bucketPassword, file, rootDir)
            }
        }
    }

    private suspend fun downloadFile(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        file: LocalFileEntity,
        parentDir: java.io.File
    ) {
        val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: return
        val name = uniqueFileNameForFile(parentDir, file.name)
        java.io.File(parentDir, name).writeBytes(content)
    }

    private suspend fun downloadDirectoryRecursive(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        dir: LocalDirectoryEntity,
        parentDir: java.io.File
    ) {
        val targetDir = java.io.File(parentDir, dir.name)
        targetDir.mkdirs()
        val files = repo.getFilesRecursive(bucket.id, dir.id)
        for (file in files) {
            val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: continue
            val relPath = repo.getDirectoryPath(bucket.id, file.directoryId ?: dir.id, dir.id)
            val fileParentDir = if (relPath.isNullOrBlank()) {
                targetDir
            } else {
                findOrCreateDirectory(targetDir, relPath)
            }
            val name = uniqueFileNameForFile(fileParentDir, file.name)
            java.io.File(fileParentDir, name).writeBytes(content)
        }
    }

    private fun findOrCreateDirectory(parent: java.io.File, path: String): java.io.File {
        val parts = path.split("/").filter { it.isNotBlank() }
        var current = parent
        for (name in parts) {
            current = java.io.File(current, name).also { it.mkdirs() }
        }
        return current
    }

}
