package xyz.nkrypt.android.data.remote

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import xyz.nkrypt.android.data.crypto.CryptoUtils
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.local.dao.RemoteBucketDao
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity
import xyz.nkrypt.android.data.remote.api.ApiError
import xyz.nkrypt.android.data.remote.api.BucketDto
import xyz.nkrypt.android.data.remote.api.DirectoryDto
import xyz.nkrypt.android.data.remote.api.FileDto
import xyz.nkrypt.android.util.generateId16
import java.io.InputStream
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class RemoteBucketRepository @Inject constructor(
    private val remoteBucketDao: RemoteBucketDao,
    private val masterPasswordStore: MasterPasswordStore
) {

    fun getAllBuckets() = remoteBucketDao.getAll()

    suspend fun getBucketById(id: String): RemoteBucketEntity? = remoteBucketDao.getById(id)

    suspend fun login(serverUrl: String, username: String, password: String): String {
        val api = NkryptApiClient.create(serverUrl)
        val response = api.login(
            xyz.nkrypt.android.data.remote.api.LoginRequest(username, password)
        )
        if (!response.isSuccessful) {
            val body = response.errorBody()?.string()
            throw ApiException("Login failed: ${response.code()} - $body")
        }
        val body = response.body() ?: throw ApiException("Empty response")
        if (body.hasError && body.error != null) {
            throw ApiException(body.error.message)
        }
        val apiKey = body.apiKey ?: throw ApiException("No API key in response")
        return apiKey
    }

    suspend fun listBuckets(serverUrl: String, apiKey: String): List<BucketDto> {
        val api = NkryptApiClient.create(serverUrl, apiKey)
        val response = api.listBuckets(
            xyz.nkrypt.android.data.remote.api.ApiKeyRequest(apiKey)
        )
        if (!response.isSuccessful) {
            throw ApiException("Failed to list buckets: ${response.code()}")
        }
        val body = response.body() ?: throw ApiException("Empty response")
        if (body.hasError && body.error != null) {
            throw ApiException(body.error.message)
        }
        return body.buckets ?: emptyList()
    }

    suspend fun addRemoteBucket(
        serverUrl: String,
        username: String,
        password: String,
        bucketId: String,
        bucketName: String,
        rootDirectoryId: String,
        encryptionPassword: String,
        masterPassword: String
    ): RemoteBucketEntity {
        val apiKey = login(serverUrl, username, password)
        val id = generateId16()
        val passwordEncrypted = encryptWithMaster(password, masterPassword)
        val encPasswordEncrypted = encryptWithMaster(encryptionPassword, masterPassword)
        val entity = RemoteBucketEntity(
            id = id,
            serverUrl = serverUrl.trim().removeSuffix("/"),
            username = username,
            passwordEncrypted = passwordEncrypted,
            bucketId = bucketId,
            bucketName = bucketName,
            rootDirectoryId = rootDirectoryId,
            encryptionPasswordEncrypted = encPasswordEncrypted,
            cachedApiKey = apiKey,
            apiKeyExpiresAt = System.currentTimeMillis() + 7 * 24 * 60 * 60 * 1000L,
            createdAt = System.currentTimeMillis()
        )
        remoteBucketDao.insert(entity)
        return entity
    }

    suspend fun getApiKey(bucket: RemoteBucketEntity): String {
        val masterPassword = masterPasswordStore.getMasterPassword() ?: throw ApiException("Master password required")
        val cached = bucket.cachedApiKey
        val expiresAt = bucket.apiKeyExpiresAt ?: 0L
        if (cached != null && expiresAt > System.currentTimeMillis()) {
            return cached
        }
        val apiKey = login(bucket.serverUrl, bucket.username, decryptWithMaster(bucket.passwordEncrypted, masterPassword))
        remoteBucketDao.updateCachedApiKey(
            bucket.id,
            apiKey,
            System.currentTimeMillis() + 7 * 24 * 60 * 60 * 1000L
        )
        return apiKey
    }

    suspend fun getDirectory(
        bucket: RemoteBucketEntity,
        directoryId: String?,
        apiKey: String
    ): Pair<List<DirectoryDto>, List<FileDto>> {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.getDirectory(
            xyz.nkrypt.android.data.remote.api.DirectoryGetRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                directoryId = directoryId
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to get directory: ${response.code()}")
        val body = response.body() ?: throw ApiException("Empty response")
        if (body.hasError && body.error != null) throw ApiException(body.error.message)
        return (body.subDirectories ?: emptyList()) to (body.files ?: emptyList())
    }

    suspend fun createDirectory(
        bucket: RemoteBucketEntity,
        parentId: String?,
        name: String,
        apiKey: String
    ): String {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.createDirectory(
            xyz.nkrypt.android.data.remote.api.DirectoryCreateRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                parentDirectoryId = parentId,
                name = name,
                metaData = emptyMap(),
                encryptedMetaData = ""
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to create directory: ${response.code()}")
        val body = response.body() ?: throw ApiException("Empty response")
        if (body.hasError && body.error != null) throw ApiException(body.error.message)
        return body.directoryId ?: throw ApiException("No directory ID in response")
    }

    /** Ensures the directory path exists on remote, creating as needed. Returns the leaf directory ID. */
    suspend fun ensureRemoteDirectoryPath(
        bucket: RemoteBucketEntity,
        path: String,
        targetRootDirId: String?,
        apiKey: String
    ): String? {
        if (path.isBlank()) return targetRootDirId
        val parts = path.split("/").filter { it.isNotEmpty() }
        var parentId = targetRootDirId
        for (name in parts) {
            val (subDirs, _) = getDirectory(bucket, parentId, apiKey)
            val existing = subDirs.find { it.name == name }
            parentId = if (existing != null) {
                existing._id
            } else {
                createDirectory(bucket, parentId, name, apiKey)
            }
        }
        return parentId
    }

    suspend fun createFileAndUpload(
        bucket: RemoteBucketEntity,
        directoryId: String?,
        name: String,
        content: ByteArray,
        apiKey: String
    ): FileDto {
        val masterPassword = masterPasswordStore.getMasterPassword() ?: throw ApiException("Master password required")
        val bucketPassword = decryptWithMaster(bucket.encryptionPasswordEncrypted, masterPassword)
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val salt = CryptoUtils.generateSalt()
        val key = CryptoUtils.createEncryptionKeyFromPassword(bucketPassword, salt)
        val iv = CryptoUtils.generateIv()
        val encrypted = CryptoUtils.encrypt(key, iv, content)
        val cryptoMeta = CryptoUtils.buildCryptoHeader(iv, salt)
        val fileResponse = api.createFile(
            xyz.nkrypt.android.data.remote.api.FileCreateRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                directoryId = directoryId,
                name = name,
                sizeInBytes = content.size.toLong(),
                metaData = emptyMap(),
                encryptedMetaData = ""
            )
        )
        if (!fileResponse.isSuccessful) throw ApiException("Failed to create file: ${fileResponse.code()}")
        val fileBody = fileResponse.body() ?: throw ApiException("Empty response")
        if (fileBody.hasError && fileBody.error != null) throw ApiException(fileBody.error.message)
        val file = fileBody.file ?: throw ApiException("No file in response")
        val blobResponse = api.writeBlob(
            bucketId = bucket.bucketId,
            fileId = file._id,
            authorization = "Bearer $apiKey",
            cryptoMeta = cryptoMeta,
            body = encrypted.toRequestBody("application/octet-stream".toMediaType())
        )
        if (!blobResponse.isSuccessful) throw ApiException("Failed to upload: ${blobResponse.code()}")
        return file
    }

    suspend fun downloadFile(
        bucket: RemoteBucketEntity,
        fileId: String,
        apiKey: String
    ): Pair<InputStream, String> {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.readBlob(bucket.bucketId, fileId, "Bearer $apiKey")
        if (!response.isSuccessful) throw ApiException("Failed to download: ${response.code()}")
        val body = response.body() ?: throw ApiException("Empty response")
        val cryptoMeta = response.headers()["nk-crypto-meta"] ?: throw ApiException("Missing crypto metadata")
        return body.byteStream() to cryptoMeta
    }

    suspend fun downloadAndDecryptFile(
        bucket: RemoteBucketEntity,
        fileId: String,
        apiKey: String
    ): ByteArray {
        val masterPassword = masterPasswordStore.getMasterPassword() ?: throw ApiException("Master password required")
        val bucketPassword = decryptWithMaster(bucket.encryptionPasswordEncrypted, masterPassword)
        val (stream, cryptoMeta) = downloadFile(bucket, fileId, apiKey)
        val encrypted = stream.readBytes()
        stream.close()
        val (iv, salt) = CryptoUtils.unbuildCryptoHeader(cryptoMeta)
        val key = CryptoUtils.createEncryptionKeyFromPassword(bucketPassword, salt)
        return CryptoUtils.decrypt(key, iv, encrypted)
    }

    suspend fun reEncryptAllWithNewMaster(oldMaster: String, newMaster: String) {
        val buckets = remoteBucketDao.getAll().first()
        for (bucket in buckets) {
            val password = decryptWithMaster(bucket.passwordEncrypted, oldMaster)
            val encPassword = decryptWithMaster(bucket.encryptionPasswordEncrypted, oldMaster)
            val newPasswordEnc = encryptWithMaster(password, newMaster)
            val newEncPasswordEnc = encryptWithMaster(encPassword, newMaster)
            remoteBucketDao.insert(
                bucket.copy(
                    passwordEncrypted = newPasswordEnc,
                    encryptionPasswordEncrypted = newEncPasswordEnc
                )
            )
        }
    }

    suspend fun deleteBucket(id: String) {
        remoteBucketDao.deleteById(id)
    }

    suspend fun renameDirectory(
        bucket: RemoteBucketEntity,
        directoryId: String,
        newName: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.renameDirectory(
            xyz.nkrypt.android.data.remote.api.DirectoryRenameRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                directoryId = directoryId,
                name = newName
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to rename: ${response.code()}")
    }

    suspend fun moveDirectory(
        bucket: RemoteBucketEntity,
        directoryId: String,
        newParentId: String,
        newName: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.moveDirectory(
            xyz.nkrypt.android.data.remote.api.DirectoryMoveRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                directoryId = directoryId,
                newParentDirectoryId = newParentId,
                newName = newName
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to move: ${response.code()}")
    }

    suspend fun deleteDirectory(
        bucket: RemoteBucketEntity,
        directoryId: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.deleteDirectory(
            xyz.nkrypt.android.data.remote.api.DirectoryDeleteRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                directoryId = directoryId
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to delete: ${response.code()}")
    }

    suspend fun renameFile(
        bucket: RemoteBucketEntity,
        fileId: String,
        newName: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.renameFile(
            xyz.nkrypt.android.data.remote.api.FileRenameRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                fileId = fileId,
                name = newName
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to rename: ${response.code()}")
    }

    suspend fun moveFile(
        bucket: RemoteBucketEntity,
        fileId: String,
        newDirectoryId: String,
        newName: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.moveFile(
            xyz.nkrypt.android.data.remote.api.FileMoveRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                fileId = fileId,
                newParentDirectoryId = newDirectoryId,
                newName = newName
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to move: ${response.code()}")
    }

    suspend fun deleteFile(
        bucket: RemoteBucketEntity,
        fileId: String,
        apiKey: String
    ) {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.deleteFile(
            xyz.nkrypt.android.data.remote.api.FileDeleteRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                fileId = fileId
            )
        )
        if (!response.isSuccessful) throw ApiException("Failed to delete: ${response.code()}")
    }

    suspend fun getFileMetadata(
        bucket: RemoteBucketEntity,
        fileId: String,
        apiKey: String
    ): xyz.nkrypt.android.data.remote.api.FileDto? {
        val api = NkryptApiClient.create(bucket.serverUrl, apiKey)
        val response = api.getFile(
            xyz.nkrypt.android.data.remote.api.FileGetRequest(
                apiKey = apiKey,
                bucketId = bucket.bucketId,
                fileId = fileId
            )
        )
        if (!response.isSuccessful) return null
        val body = response.body() ?: return null
        if (body.hasError || body.file == null) return null
        return body.file
    }

    fun decryptEncryptionPassword(encrypted: String, masterPassword: String): String {
        return decryptWithMaster(encrypted, masterPassword)
    }

    private fun encryptWithMaster(plain: String, master: String): String {
        val payload = CryptoUtils.encryptText(plain, master)
        return """{"cipher":"${payload.cipher}","iv":"${payload.iv}","salt":"${payload.salt}"}"""
    }

    private fun decryptWithMaster(json: String, master: String): String {
        val map = com.google.gson.Gson().fromJson(json, Map::class.java)
        val payload = CryptoUtils.EncryptedPayload(
            cipher = map["cipher"] as String,
            iv = map["iv"] as String,
            salt = map["salt"] as String
        )
        return CryptoUtils.decryptText(payload, master)
    }
}

class ApiException(message: String) : Exception(message)
