package xyz.nkrypt.android.ui.localbuckets

import android.content.Context
import android.content.Intent
import android.net.Uri
import android.provider.DocumentsContract
import androidx.documentfile.provider.DocumentFile
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

    suspend fun downloadBucketToDirectory(bucket: LocalBucketEntity, destTreeUri: Uri, context: Context) {
        val masterPassword = masterPasswordStore.getMasterPassword() ?: return
        val bucketPassword = repository.decryptBucketPassword(bucket.cryptData, masterPassword)
        context.contentResolver.takePersistableUriPermission(
            destTreeUri,
            Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_GRANT_WRITE_URI_PERMISSION
        )
        withContext(Dispatchers.IO) {
            val root = DocumentFile.fromTreeUri(context, destTreeUri) ?: return@withContext
            val rootDirs = repository.getDirectories(bucket.id, null)
            val rootFiles = repository.getFiles(bucket.id, null)
            for (dir in rootDirs) {
                downloadDirectoryRecursive(repository, bucket, bucketPassword, dir, root, context)
            }
            for (file in rootFiles) {
                downloadFile(repository, bucket, bucketPassword, file, root, context)
            }
        }
    }

    private suspend fun downloadFile(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        file: LocalFileEntity,
        parentDir: DocumentFile,
        context: Context
    ) {
        val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: return
        val docFile = parentDir.createFile("application/octet-stream", file.name) ?: return
        context.contentResolver.openOutputStream(docFile.uri)?.use { it.write(content) }
    }

    private suspend fun downloadDirectoryRecursive(
        repo: LocalBucketRepository,
        bucket: LocalBucketEntity,
        bucketPassword: String,
        dir: LocalDirectoryEntity,
        parentDir: DocumentFile,
        context: Context
    ) {
        val targetDir = parentDir.createDirectory(dir.name) ?: return
        val files = repo.getFilesRecursive(bucket.id, dir.id)
        for (file in files) {
            val content = repo.readFileContent(bucket, file.id, bucketPassword) ?: continue
            val relPath = repo.getDirectoryPath(bucket.id, file.directoryId ?: dir.id, dir.id)
            val fileParentDir = if (relPath.isNullOrBlank()) {
                targetDir
            } else {
                findOrCreateDirectory(targetDir, relPath)
            }
            val docFile = fileParentDir?.createFile("application/octet-stream", file.name) ?: continue
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
