package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "remote_buckets")
data class RemoteBucketEntity(
    @PrimaryKey
    val id: String,
    val serverUrl: String,
    val username: String,
    val passwordEncrypted: String,
    val bucketId: String,
    val bucketName: String,
    val rootDirectoryId: String,
    val encryptionPasswordEncrypted: String,
    val cachedApiKey: String?,
    val apiKeyExpiresAt: Long?,
    val createdAt: Long
)
