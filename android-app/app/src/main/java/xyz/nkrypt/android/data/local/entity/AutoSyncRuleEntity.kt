package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

/** Post-action after syncing a file to remote. */
enum class SyncPostAction {
    DELETE_LOCAL,  // Delete from local bucket
    KEEP          // Keep in local bucket
}

@Entity(tableName = "auto_sync_rules")
data class AutoSyncRuleEntity(
    @PrimaryKey
    val id: String,
    val name: String,
    val sourceBucketId: String,
    val sourceDirectoryId: String?,  // null = root of bucket
    val targetRemoteBucketId: String,
    val targetDirectoryId: String?, // null = root of remote bucket
    val postAction: String,         // SyncPostAction.name()
    val createdAt: Long
)
