package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.ForeignKey
import androidx.room.Index
import androidx.room.PrimaryKey

@Entity(
    tableName = "local_blobs",
    foreignKeys = [
        ForeignKey(
            entity = LocalFileEntity::class,
            parentColumns = ["id"],
            childColumns = ["fileId"],
            onDelete = ForeignKey.CASCADE
        )
    ],
    indices = [Index("fileId")]
)
data class LocalBlobEntity(
    @PrimaryKey
    val id: String,
    val fileId: String,
    val sizeInBytes: Long,
    val blobPath: String,
    val ivBase64: String,
    val saltBase64: String,
    val contentHashHex: String?,
    val contentHashSaltBase64: String?,
    val createdAt: Long
)
