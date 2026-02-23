package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.ForeignKey
import androidx.room.Index
import androidx.room.PrimaryKey

@Entity(
    tableName = "local_files",
    foreignKeys = [
        ForeignKey(
            entity = LocalBucketEntity::class,
            parentColumns = ["id"],
            childColumns = ["bucketId"],
            onDelete = ForeignKey.CASCADE
        ),
        ForeignKey(
            entity = LocalDirectoryEntity::class,
            parentColumns = ["id"],
            childColumns = ["directoryId"],
            onDelete = ForeignKey.CASCADE
        )
    ],
    indices = [
        Index("bucketId"),
        Index("directoryId"),
        Index(value = ["bucketId", "directoryId", "name"], unique = true)
    ]
)
data class LocalFileEntity(
    @PrimaryKey
    val id: String,
    val bucketId: String,
    val directoryId: String?,
    val name: String,
    val sizeInBytes: Long,
    val metaData: String,
    val encryptedMetaData: String,
    val createdAt: Long
)
