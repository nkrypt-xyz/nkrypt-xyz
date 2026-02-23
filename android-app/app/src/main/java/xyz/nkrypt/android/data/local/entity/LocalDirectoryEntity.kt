package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.ForeignKey
import androidx.room.Index
import androidx.room.PrimaryKey

@Entity(
    tableName = "local_directories",
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
            childColumns = ["parentDirectoryId"],
            onDelete = ForeignKey.CASCADE
        )
    ],
    indices = [
        Index("bucketId"),
        Index("parentDirectoryId"),
        Index(value = ["bucketId", "parentDirectoryId", "name"], unique = true)
    ]
)
data class LocalDirectoryEntity(
    @PrimaryKey
    val id: String,
    val bucketId: String,
    val parentDirectoryId: String?,
    val name: String,
    val metaData: String,
    val encryptedMetaData: String,
    val createdAt: Long
)
