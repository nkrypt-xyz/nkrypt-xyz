package xyz.nkrypt.android.data.local

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.withContext
import xyz.nkrypt.android.data.local.dao.AutoImportRuleDao
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import xyz.nkrypt.android.data.local.dao.AutoSyncRuleDao
import xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity
import xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity
import xyz.nkrypt.android.data.local.entity.ImportPostAction
import xyz.nkrypt.android.data.local.entity.SyncPostAction
import xyz.nkrypt.android.util.generateId16
import java.io.File
import javax.inject.Inject
import javax.inject.Singleton

/** Trash directory: /storage/emulated/0/.trash/ */
const val TRASH_ROOT = "/storage/emulated/0/.trash"

sealed class ImportProgress {
    data object Scanning : ImportProgress()
    data class Progress(val current: Int, val total: Int, val fileName: String, val phase: String) : ImportProgress()
    data class Done(val imported: Int, val errors: Int) : ImportProgress()
    data class Error(val message: String) : ImportProgress()
}

sealed class SyncProgress {
    data object Scanning : SyncProgress()
    data class Progress(val current: Int, val total: Int, val fileName: String, val phase: String) : SyncProgress()
    data class Done(val synced: Int, val errors: Int) : SyncProgress()
    data class Error(val message: String) : SyncProgress()
}

@Singleton
class RulesRepository @Inject constructor(
    private val autoImportRuleDao: AutoImportRuleDao,
    private val autoSyncRuleDao: AutoSyncRuleDao,
    private val localBucketRepository: LocalBucketRepository,
    private val remoteBucketRepository: RemoteBucketRepository,
    private val masterPasswordStore: MasterPasswordStore
) {

    fun getAllImportRules(): Flow<List<AutoImportRuleEntity>> = autoImportRuleDao.getAll()
    fun getAllSyncRules(): Flow<List<AutoSyncRuleEntity>> = autoSyncRuleDao.getAll()

    suspend fun getImportRuleById(id: String): AutoImportRuleEntity? = autoImportRuleDao.getById(id)
    suspend fun getSyncRuleById(id: String): AutoSyncRuleEntity? = autoSyncRuleDao.getById(id)

    suspend fun createImportRule(
        name: String,
        sourceDirectoryPath: String,
        targetBucketId: String,
        postAction: ImportPostAction
    ): AutoImportRuleEntity {
        val entity = AutoImportRuleEntity(
            id = generateId16(),
            name = name,
            sourceDirectoryPath = sourceDirectoryPath,
            targetBucketId = targetBucketId,
            postAction = postAction.name,
            createdAt = System.currentTimeMillis()
        )
        autoImportRuleDao.insert(entity)
        return entity
    }

    suspend fun createSyncRule(
        name: String,
        sourceBucketId: String,
        sourceDirectoryId: String?,
        targetRemoteBucketId: String,
        targetDirectoryId: String?,
        postAction: SyncPostAction
    ): AutoSyncRuleEntity {
        val entity = AutoSyncRuleEntity(
            id = generateId16(),
            name = name,
            sourceBucketId = sourceBucketId,
            sourceDirectoryId = sourceDirectoryId,
            targetRemoteBucketId = targetRemoteBucketId,
            targetDirectoryId = targetDirectoryId,
            postAction = postAction.name,
            createdAt = System.currentTimeMillis()
        )
        autoSyncRuleDao.insert(entity)
        return entity
    }

    suspend fun updateImportRule(
        id: String,
        name: String,
        sourceDirectoryPath: String,
        targetBucketId: String,
        postAction: ImportPostAction
    ) {
        val existing = autoImportRuleDao.getById(id) ?: return
        val updated = existing.copy(
            name = name,
            sourceDirectoryPath = sourceDirectoryPath,
            targetBucketId = targetBucketId,
            postAction = postAction.name
        )
        autoImportRuleDao.insert(updated)
    }

    suspend fun updateSyncRule(
        id: String,
        name: String,
        sourceBucketId: String,
        sourceDirectoryId: String?,
        targetRemoteBucketId: String,
        targetDirectoryId: String?,
        postAction: SyncPostAction
    ) {
        val existing = autoSyncRuleDao.getById(id) ?: return
        val updated = existing.copy(
            name = name,
            sourceBucketId = sourceBucketId,
            sourceDirectoryId = sourceDirectoryId,
            targetRemoteBucketId = targetRemoteBucketId,
            targetDirectoryId = targetDirectoryId,
            postAction = postAction.name
        )
        autoSyncRuleDao.insert(updated)
    }

    suspend fun deleteImportRule(id: String) = autoImportRuleDao.deleteById(id)
    suspend fun deleteSyncRule(id: String) = autoSyncRuleDao.deleteById(id)

    fun executeImport(ruleId: String): Flow<ImportProgress> = flow {
        val rule = autoImportRuleDao.getById(ruleId) ?: run {
            emit(ImportProgress.Error("Rule not found"))
            return@flow
        }
        val bucket = localBucketRepository.getBucketById(rule.targetBucketId) ?: run {
            emit(ImportProgress.Error("Target bucket not found"))
            return@flow
        }
        val masterPassword = masterPasswordStore.getMasterPassword() ?: run {
            emit(ImportProgress.Error("Master password required"))
            return@flow
        }
        val bucketPassword = localBucketRepository.decryptBucketPassword(bucket.cryptData, masterPassword)

        val sourceDir = File(rule.sourceDirectoryPath)
        if (!sourceDir.exists() || !sourceDir.isDirectory) {
            emit(ImportProgress.Error("Source directory not found"))
            return@flow
        }

        emit(ImportProgress.Scanning)
        val files = withContext(Dispatchers.IO) { collectFilesRecursive(sourceDir) }
        val total = files.size
        if (total == 0) {
            emit(ImportProgress.Done(0, 0))
            return@flow
        }

        var imported = 0
        var errors = 0
        val postAction = ImportPostAction.valueOf(rule.postAction)
        val trashRoot = File(TRASH_ROOT)

        for ((index, file) in files.withIndex()) {
            try {
                emit(ImportProgress.Progress(index + 1, total, file.name, "Encrypting"))
                val content = file.readBytes()
                val relativePath = file.relativeTo(sourceDir).parent?.toString() ?: "."
                val targetDirId = ensureDirectoryPath(bucket.id, relativePath, localBucketRepository)
                emit(ImportProgress.Progress(index + 1, total, file.name, "Importing"))
                localBucketRepository.createFileWithContent(
                    bucket = bucket,
                    directoryId = targetDirId.ifEmpty { null },
                    name = file.name,
                    content = content,
                    bucketPassword = bucketPassword
                )
                imported++

                when (postAction) {
                    ImportPostAction.DELETE -> file.delete()
                    ImportPostAction.TRASH -> {
                        val trashPath = File(trashRoot, file.relativeTo(sourceDir).path)
                        trashPath.parentFile?.mkdirs()
                        file.renameTo(trashPath)
                    }
                    ImportPostAction.KEEP -> {}
                }
            } catch (e: Exception) {
                errors++
                emit(ImportProgress.Progress(index + 1, total, file.name, "Error: ${e.message}"))
            }
        }
        emit(ImportProgress.Done(imported, errors))
    }

    fun executeSync(ruleId: String): Flow<SyncProgress> = flow {
        val rule = autoSyncRuleDao.getById(ruleId) ?: run {
            emit(SyncProgress.Error("Rule not found"))
            return@flow
        }
        val localBucket = localBucketRepository.getBucketById(rule.sourceBucketId) ?: run {
            emit(SyncProgress.Error("Source bucket not found"))
            return@flow
        }
        val remoteBucket = remoteBucketRepository.getBucketById(rule.targetRemoteBucketId) ?: run {
            emit(SyncProgress.Error("Target remote bucket not found"))
            return@flow
        }
        val masterPassword = masterPasswordStore.getMasterPassword() ?: run {
            emit(SyncProgress.Error("Master password required"))
            return@flow
        }
        val apiKey = try {
            remoteBucketRepository.getApiKey(remoteBucket)
        } catch (e: Exception) {
            emit(SyncProgress.Error("Failed to get API key: ${e.message}"))
            return@flow
        }
        val localBucketPassword = localBucketRepository.decryptBucketPassword(localBucket.cryptData, masterPassword)

        emit(SyncProgress.Scanning)
        val files = withContext(Dispatchers.IO) {
            localBucketRepository.getFilesRecursive(localBucket.id, rule.sourceDirectoryId)
        }
        val total = files.size
        if (total == 0) {
            emit(SyncProgress.Done(0, 0))
            return@flow
        }

        var synced = 0
        var errors = 0
        val postAction = SyncPostAction.valueOf(rule.postAction)

        for ((index, file) in files.withIndex()) {
            try {
                emit(SyncProgress.Progress(index + 1, total, file.name, "Reading"))
                val content = localBucketRepository.readFileContent(localBucket, file.id, localBucketPassword)
                    ?: throw Exception("Failed to read file")
                emit(SyncProgress.Progress(index + 1, total, file.name, "Uploading"))
                val targetRoot = rule.targetDirectoryId ?: remoteBucket.rootDirectoryId
                val remoteDirId = mapLocalDirToRemote(
                    localBucketRepository,
                    remoteBucketRepository,
                    localBucket,
                    remoteBucket,
                    file.directoryId,
                    rule.sourceDirectoryId,
                    targetRoot,
                    apiKey
                ) ?: targetRoot
                remoteBucketRepository.createFileAndUpload(
                    bucket = remoteBucket,
                    directoryId = remoteDirId,
                    name = file.name,
                    content = content,
                    apiKey = apiKey
                )
                synced++

                if (postAction == SyncPostAction.DELETE_LOCAL) {
                    localBucketRepository.deleteFile(localBucket.id, file.id)
                }
            } catch (e: Exception) {
                errors++
                emit(SyncProgress.Progress(index + 1, total, file.name, "Error: ${e.message}"))
            }
        }
        emit(SyncProgress.Done(synced, errors))
    }

    private fun collectFilesRecursive(dir: File): List<File> {
        val result = mutableListOf<File>()
        dir.listFiles()?.forEach { file ->
            if (file.isDirectory) {
                result.addAll(collectFilesRecursive(file))
            } else {
                result.add(file)
            }
        }
        return result
    }

    private suspend fun ensureDirectoryPath(
        bucketId: String,
        relativePath: String,
        repo: LocalBucketRepository
    ): String {
        if (relativePath == "." || relativePath.isEmpty()) return ""
        val parts = relativePath.split(File.separator).filter { it.isNotEmpty() }
        var parentId: String? = null
        for (part in parts) {
            val existing = repo.getDirectories(bucketId, parentId).find { it.name == part }
            parentId = if (existing != null) {
                existing.id
            } else {
                repo.createDirectory(bucketId, parentId, part).id
            }
        }
        return parentId ?: ""
    }

    private suspend fun mapLocalDirToRemote(
        localRepo: LocalBucketRepository,
        remoteRepo: RemoteBucketRepository,
        localBucket: xyz.nkrypt.android.data.local.entity.LocalBucketEntity,
        remoteBucket: xyz.nkrypt.android.data.local.entity.RemoteBucketEntity,
        localDirId: String?,
        sourceRootDirId: String?,
        targetRootDirId: String?,
        apiKey: String
    ): String? {
        if (localDirId == null) return targetRootDirId
        val localPath = localRepo.getDirectoryPath(localBucket.id, localDirId, sourceRootDirId) ?: return targetRootDirId
        return remoteRepo.ensureRemoteDirectoryPath(remoteBucket, localPath, targetRootDirId, apiKey)
    }
}
