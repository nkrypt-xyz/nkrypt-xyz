package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "local_buckets")
data class LocalBucketEntity(
    @PrimaryKey
    val id: String,
    val name: String,
    val rootPath: String,
    val cryptSpec: String,
    val cryptData: String,
    val metaData: String,
    val createdAt: Long
)
