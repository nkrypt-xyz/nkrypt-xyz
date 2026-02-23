package xyz.nkrypt.android.data.local

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import xyz.nkrypt.android.data.crypto.CryptoUtils
import xyz.nkrypt.android.data.local.dao.LocalBlobDao
import xyz.nkrypt.android.data.local.dao.LocalBucketDao
import xyz.nkrypt.android.data.local.dao.LocalDirectoryDao
import xyz.nkrypt.android.data.local.dao.LocalFileDao
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity
import xyz.nkrypt.android.util.generateId16
import java.io.File
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class LocalBucketRepository @Inject constructor(
    private val bucketDao: LocalBucketDao,
    private val directoryDao: LocalDirectoryDao,
    private val fileDao: LocalFileDao,
    private val blobDao: LocalBlobDao
) {

    fun getAllBuckets(): Flow<List<LocalBucketEntity>> = bucketDao.getAll()

    suspend fun getBucketById(id: String): LocalBucketEntity? = bucketDao.getById(id)

    suspend fun createBucket(
        name: String,
        rootPath: String,
        encryptionPassword: String,
        masterPassword: String
    ): LocalBucketEntity {
        val id = generateId16()
        val cryptData = encryptBucketPassword(encryptionPassword, masterPassword)
        val entity = LocalBucketEntity(
            id = id,
            name = name,
            rootPath = rootPath,
            cryptSpec = "AES-256-GCM",
            cryptData = cryptData,
            metaData = "{}",
            createdAt = System.currentTimeMillis()
        )
        bucketDao.insert(entity)
        val rootDir = File(rootPath)
        if (!rootDir.exists()) rootDir.mkdirs()
        val blobsDir = File(rootDir, "blobs")
        if (!blobsDir.exists()) blobsDir.mkdirs()
        return entity
    }

    suspend fun reEncryptAllWithNewMaster(oldMaster: String, newMaster: String) {
        val buckets = bucketDao.getAll().first()
        for (bucket in buckets) {
            val password = decryptBucketPassword(bucket.cryptData, oldMaster)
            val newCryptData = encryptBucketPassword(password, newMaster)
            bucketDao.insert(bucket.copy(cryptData = newCryptData))
        }
    }

    suspend fun deleteBucket(id: String, deleteFiles: Boolean = false) {
        val bucket = bucketDao.getById(id) ?: return
        if (deleteFiles) {
            val rootDir = File(bucket.rootPath)
            if (rootDir.exists()) rootDir.deleteRecursively()
        }
        bucketDao.deleteById(id)
    }

    fun decryptBucketPassword(cryptData: String, masterPassword: String): String {
        val payload = parseCryptData(cryptData)
        return CryptoUtils.decryptText(payload, masterPassword)
    }

    private fun encryptBucketPassword(password: String, masterPassword: String): String {
        val payload = CryptoUtils.encryptText(password, masterPassword)
        return """{"cipher":"${payload.cipher}","iv":"${payload.iv}","salt":"${payload.salt}"}"""
    }

    private fun parseCryptData(json: String): CryptoUtils.EncryptedPayload {
        val map = com.google.gson.Gson().fromJson(json, Map::class.java)
        return CryptoUtils.EncryptedPayload(
            cipher = map["cipher"] as String,
            iv = map["iv"] as String,
            salt = map["salt"] as String
        )
    }

    suspend fun getDirectories(bucketId: String, parentId: String?): List<LocalDirectoryEntity> =
        directoryDao.getByBucketAndParent(bucketId, parentId)

    suspend fun getFiles(bucketId: String, directoryId: String?): List<LocalFileEntity> =
        fileDao.getByBucketAndDirectory(bucketId, directoryId)

    suspend fun createDirectory(
        bucketId: String,
        parentId: String?,
        name: String,
        encryptedMetaData: String = ""
    ): LocalDirectoryEntity {
        val existing = directoryDao.getByNameAndParent(bucketId, parentId, name)
        if (existing != null) {
            throw IllegalArgumentException("A directory with this name already exists in the parent.")
        }
        val id = generateId16()
        val entity = LocalDirectoryEntity(
            id = id,
            bucketId = bucketId,
            parentDirectoryId = parentId,
            name = name,
            metaData = "{}",
            encryptedMetaData = encryptedMetaData,
            createdAt = System.currentTimeMillis()
        )
        directoryDao.insert(entity)
        return entity
    }

    suspend fun getDirectoryById(id: String): LocalDirectoryEntity? = directoryDao.getById(id)

    suspend fun getFileById(id: String): LocalFileEntity? = fileDao.getById(id)

    suspend fun getBlobPath(bucket: LocalBucketEntity, fileId: String): File? {
        val blob = blobDao.getLatestByFileId(fileId) ?: return null
        return File(bucket.rootPath, blob.blobPath)
    }

    suspend fun readFileContent(
        bucket: LocalBucketEntity,
        fileId: String,
        bucketPassword: String
    ): ByteArray? {
        val blob = blobDao.getLatestByFileId(fileId) ?: return null
        val blobFile = File(bucket.rootPath, blob.blobPath)
        if (!blobFile.exists()) return null
        val encrypted = blobFile.readBytes()
        val iv = android.util.Base64.decode(blob.ivBase64, android.util.Base64.NO_WRAP)
        val salt = android.util.Base64.decode(blob.saltBase64, android.util.Base64.NO_WRAP)
        val key = CryptoUtils.createEncryptionKeyFromPassword(bucketPassword, salt)
        return CryptoUtils.decrypt(key, iv, encrypted)
    }

    suspend fun getFilesRecursive(bucketId: String, directoryId: String?): List<LocalFileEntity> {
        val result = mutableListOf<LocalFileEntity>()
        collectFilesRecursive(bucketId, directoryId, result)
        return result
    }

    private suspend fun collectFilesRecursive(
        bucketId: String,
        directoryId: String?,
        result: MutableList<LocalFileEntity>
    ) {
        result.addAll(fileDao.getByBucketAndDirectory(bucketId, directoryId))
        directoryDao.getByBucketAndParent(bucketId, directoryId).forEach { dir ->
            collectFilesRecursive(bucketId, dir.id, result)
        }
    }

    /** Returns path from root (or sourceRootDirId) to the given directory, e.g. "a/b/c". */
    suspend fun getDirectoryPath(
        bucketId: String,
        directoryId: String,
        stopAtDirectoryId: String?
    ): String? {
        val parts = mutableListOf<String>()
        var currentId: String? = directoryId
        while (currentId != null && currentId != stopAtDirectoryId) {
            val dir = directoryDao.getById(currentId) ?: return null
            parts.add(0, dir.name)
            currentId = dir.parentDirectoryId
        }
        return if (parts.isEmpty()) null else parts.joinToString("/")
    }

    suspend fun deleteFile(bucketId: String, fileId: String) {
        val bucket = bucketDao.getById(bucketId) ?: return
        val blob = blobDao.getLatestByFileId(fileId)
        if (blob != null) {
            val blobFile = File(bucket.rootPath, blob.blobPath)
            if (blobFile.exists()) blobFile.delete()
            blobDao.deleteByFileId(fileId)
        }
        fileDao.deleteById(fileId)
    }

    suspend fun renameFile(bucketId: String, fileId: String, newName: String) {
        val file = fileDao.getById(fileId) ?: return
        val existing = fileDao.getByNameAndDirectory(bucketId, file.directoryId, newName)
        if (existing != null && existing.id != fileId) {
            throw IllegalArgumentException("A file with this name already exists in the directory.")
        }
        fileDao.rename(fileId, newName)
    }

    suspend fun moveFile(bucketId: String, fileId: String, newDirectoryId: String?) {
        val file = fileDao.getById(fileId) ?: return
        val existing = fileDao.getByNameAndDirectory(bucketId, newDirectoryId, file.name)
        if (existing != null && existing.id != fileId) {
            throw IllegalArgumentException("A file with this name already exists in the target directory.")
        }
        fileDao.move(fileId, newDirectoryId)
    }

    suspend fun deleteDirectory(bucketId: String, directoryId: String) {
        val bucket = bucketDao.getById(bucketId) ?: return
        val files = getFilesRecursive(bucketId, directoryId)
        for (file in files) {
            val blob = blobDao.getLatestByFileId(file.id)
            if (blob != null) {
                val blobFile = File(bucket.rootPath, blob.blobPath)
                if (blobFile.exists()) blobFile.delete()
                blobDao.deleteByFileId(file.id)
            }
            fileDao.deleteById(file.id)
        }
        deleteDirectoryRecursive(bucketId, directoryId)
    }

    private suspend fun deleteDirectoryRecursive(bucketId: String, directoryId: String) {
        val subdirs = directoryDao.getByBucketAndParent(bucketId, directoryId)
        for (dir in subdirs) {
            deleteDirectoryRecursive(bucketId, dir.id)
        }
        directoryDao.deleteById(directoryId)
    }

    suspend fun renameDirectory(bucketId: String, directoryId: String, newName: String) {
        val dir = directoryDao.getById(directoryId) ?: return
        val existing = directoryDao.getByNameAndParent(bucketId, dir.parentDirectoryId, newName)
        if (existing != null && existing.id != directoryId) {
            throw IllegalArgumentException("A directory with this name already exists in the parent.")
        }
        directoryDao.rename(directoryId, newName)
    }

    suspend fun moveDirectory(bucketId: String, directoryId: String, newParentId: String?) {
        val dir = directoryDao.getById(directoryId) ?: return
        val existing = directoryDao.getByNameAndParent(bucketId, newParentId, dir.name)
        if (existing != null && existing.id != directoryId) {
            throw IllegalArgumentException("A directory with this name already exists in the target.")
        }
        directoryDao.move(directoryId, newParentId)
    }

    suspend fun createFileWithContent(
        bucket: LocalBucketEntity,
        directoryId: String?,
        name: String,
        content: ByteArray,
        bucketPassword: String
    ): LocalFileEntity {
        val existing = fileDao.getByNameAndDirectory(bucket.id, directoryId, name)
        if (existing != null) {
            throw IllegalArgumentException("A file with this name already exists in the directory.")
        }
        val fileId = generateId16()
        val blobId = generateId16()
        val salt = CryptoUtils.generateSalt()
        val key = CryptoUtils.createEncryptionKeyFromPassword(bucketPassword, salt)
        val iv = CryptoUtils.generateIv()
        val encrypted = CryptoUtils.encrypt(key, iv, content)
        val blobPath = "blobs/$fileId/$blobId.enc"
        val blobFile = File(bucket.rootPath, blobPath)
        blobFile.parentFile?.mkdirs()
        blobFile.writeBytes(encrypted)
        val fileEntity = LocalFileEntity(
            id = fileId,
            bucketId = bucket.id,
            directoryId = directoryId,
            name = name,
            sizeInBytes = content.size.toLong(),
            metaData = "{}",
            encryptedMetaData = "",
            createdAt = System.currentTimeMillis()
        )
        val blobEntity = xyz.nkrypt.android.data.local.entity.LocalBlobEntity(
            id = blobId,
            fileId = fileId,
            sizeInBytes = encrypted.size.toLong(),
            blobPath = blobPath,
            ivBase64 = android.util.Base64.encodeToString(iv, android.util.Base64.NO_WRAP),
            saltBase64 = android.util.Base64.encodeToString(salt, android.util.Base64.NO_WRAP),
            createdAt = System.currentTimeMillis()
        )
        fileDao.insert(fileEntity)
        blobDao.insert(blobEntity)
        return fileEntity
    }
}
