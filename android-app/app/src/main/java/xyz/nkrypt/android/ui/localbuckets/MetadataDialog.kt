package xyz.nkrypt.android.ui.localbuckets

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity

@Composable
fun BucketMetadataDialog(
    bucket: LocalBucketEntity,
    onDismiss: () -> Unit
) {
    val scrollState = rememberScrollState()
    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text("Bucket metadata", style = MaterialTheme.typography.titleLarge)
        },
        text = {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .heightIn(max = 400.dp)
                    .verticalScroll(scrollState)
            ) {
                MetadataContent(
                    items = listOf(
                        "Name" to bucket.name,
                        "ID" to bucket.id,
                        "Root path" to bucket.rootPath,
                        "Crypt spec" to bucket.cryptSpec,
                        "Created" to formatTimestamp(bucket.createdAt),
                        "Meta data" to formatJson(bucket.metaData),
                        "Crypt data" to truncateIfLong(bucket.cryptData)
                    )
                )
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) {
                Text("Close")
            }
        }
    )
}

@Composable
fun MetadataDialog(
    directory: LocalDirectoryEntity?,
    file: LocalFileEntity?,
    onDismiss: () -> Unit
) {
    val scrollState = rememberScrollState()
    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text(
                text = when {
                    directory != null -> "Folder metadata"
                    file != null -> "File metadata"
                    else -> "Metadata"
                },
                style = MaterialTheme.typography.titleLarge
            )
        },
        text = {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .heightIn(max = 400.dp)
                    .verticalScroll(scrollState)
            ) {
                when {
                    directory != null -> DirectoryMetadataContent(directory)
                    file != null -> FileMetadataContent(file)
                }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) {
                Text("Close")
            }
        }
    )
}

@Composable
private fun DirectoryMetadataContent(dir: LocalDirectoryEntity) {
    MetadataContent(
        items = listOf(
            "Type" to "Folder",
            "Name" to dir.name,
            "ID" to dir.id,
            "Bucket ID" to dir.bucketId,
            "Parent ID" to (dir.parentDirectoryId ?: "(root)"),
            "Created" to formatTimestamp(dir.createdAt),
            "Meta data" to formatJson(dir.metaData),
            "Encrypted meta" to truncateIfLong(dir.encryptedMetaData)
        )
    )
}

@Composable
private fun FileMetadataContent(file: LocalFileEntity) {
    MetadataContent(
        items = listOf(
            "Type" to "File",
            "Name" to file.name,
            "ID" to file.id,
            "Bucket ID" to file.bucketId,
            "Directory ID" to (file.directoryId ?: "(root)"),
            "Size" to formatFileSize(file.sizeInBytes),
            "Created" to formatTimestamp(file.createdAt),
            "Meta data" to formatJson(file.metaData),
            "Encrypted meta" to truncateIfLong(file.encryptedMetaData)
        )
    )
}

@Composable
private fun MetadataContent(items: List<Pair<String, String>>) {
    Column(modifier = Modifier.padding(vertical = 8.dp)) {
        items.forEach { (label, value) ->
            Text(
                text = label,
                style = MaterialTheme.typography.labelMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Text(
                text = value,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface,
                modifier = Modifier.padding(bottom = 12.dp)
            )
        }
    }
}

private fun formatTimestamp(ms: Long): String {
    val s = java.util.concurrent.TimeUnit.MILLISECONDS.toSeconds(ms)
    val date = java.util.Date(s * 1000)
    return java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss", java.util.Locale.US).format(date)
}

private fun formatFileSize(bytes: Long): String = when {
    bytes < 1024 -> "$bytes B"
    bytes < 1024 * 1024 -> "${bytes / 1024} KB"
    bytes < 1024 * 1024 * 1024 -> "${bytes / (1024 * 1024)} MB"
    else -> "${bytes / (1024 * 1024 * 1024)} GB"
}

private fun formatJson(json: String): String {
    if (json.isBlank()) return "(empty)"
    return try {
        val gson = com.google.gson.GsonBuilder().setPrettyPrinting().create()
        val element = com.google.gson.JsonParser.parseString(json)
        gson.toJson(element)
    } catch (_: Exception) {
        json
    }
}

private fun truncateIfLong(s: String, maxLen: Int = 80): String {
    if (s.isBlank()) return "(empty)"
    return if (s.length <= maxLen) s else "${s.take(maxLen)}…"
}
