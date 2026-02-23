package xyz.nkrypt.android.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

/** Post-action after importing a file from source directory. */
enum class ImportPostAction {
    DELETE,   // Delete original file
    TRASH,   // Move to /storage/emulated/0/.trash/
    KEEP     // Keep original
}

@Entity(tableName = "auto_import_rules")
data class AutoImportRuleEntity(
    @PrimaryKey
    val id: String,
    val name: String,
    val sourceDirectoryPath: String,
    val targetBucketId: String,
    val postAction: String, // ImportPostAction.name()
    val createdAt: Long
)
