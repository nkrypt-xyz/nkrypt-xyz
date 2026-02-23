package xyz.nkrypt.android.ui.filepicker

import android.os.Build
import android.os.Environment
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Description
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.FolderOpen
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import java.io.File

enum class FilePickerMode {
    DIRECTORY,
    FILE
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun FilePickerScreen(
    mode: FilePickerMode,
    onSelect: (String) -> Unit,
    onDismiss: () -> Unit,
    message: String = ""
) {
    val pathStack = remember { mutableStateListOf<String>().apply { add(getInitialPath()) } }
    val currentPath = pathStack.last()
    var refreshTrigger by remember { mutableStateOf(0) }
    var showCreateDirDialog by remember { mutableStateOf(false) }
    var newFolderName by remember { mutableStateOf("") }
    var createDirError by remember { mutableStateOf<String?>(null) }

    val (entries, loadError) = remember(currentPath, refreshTrigger) {
        loadDirectoryEntries(currentPath)
    }

    if (showCreateDirDialog) {
        AlertDialog(
            onDismissRequest = {
                showCreateDirDialog = false
                newFolderName = ""
                createDirError = null
            },
            title = { Text("New folder") },
            text = {
                Column {
                    OutlinedTextField(
                        value = newFolderName,
                        onValueChange = {
                            newFolderName = it
                            createDirError = null
                        },
                        label = { Text("Folder name") },
                        singleLine = true,
                        isError = createDirError != null
                    )
                    createDirError?.let { err ->
                        Text(
                            text = err,
                            color = MaterialTheme.colorScheme.error,
                            style = MaterialTheme.typography.bodySmall,
                            modifier = Modifier.padding(top = 8.dp)
                        )
                    }
                }
            },
            confirmButton = {
                TextButton(
                    onClick = {
                        val name = newFolderName.trim()
                        when {
                            name.isEmpty() -> createDirError = "Name cannot be empty"
                            name.contains("/") || name == ".." || name == "." ->
                                createDirError = "Invalid folder name"
                            else -> {
                                val dir = File(currentPath, name)
                                when {
                                    dir.exists() -> createDirError = "Folder already exists"
                                    dir.mkdirs() -> {
                                        refreshTrigger++
                                        showCreateDirDialog = false
                                        newFolderName = ""
                                        createDirError = null
                                        pathStack.add(dir.absolutePath)
                                    }
                                    else -> createDirError = "Could not create folder"
                                }
                            }
                        }
                    }
                ) { Text("Create") }
            },
            dismissButton = {
                TextButton(
                    onClick = {
                        showCreateDirDialog = false
                        newFolderName = ""
                        createDirError = null
                    }
                ) { Text("Cancel") }
            }
        )
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = when (mode) {
                            FilePickerMode.DIRECTORY -> "Select folder"
                            FilePickerMode.FILE -> "Select file"
                        },
                        style = MaterialTheme.typography.titleLarge
                    )
                },
                navigationIcon = {
                    IconButton(
                        onClick = {
                            if (pathStack.size > 1) {
                                pathStack.removeAt(pathStack.lastIndex)
                            } else {
                                onDismiss()
                            }
                        }
                    ) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                },
                actions = {
                    IconButton(onClick = { showCreateDirDialog = true }) {
                        Icon(Icons.Default.Add, contentDescription = "New folder")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp)
        ) {
            if (message.isNotBlank()) {
                Text(
                    text = message,
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurface,
                    modifier = Modifier.padding(bottom = 8.dp)
                )
            }
            Text(
                text = currentPath,
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.padding(bottom = 8.dp)
            )
            if (loadError != null) {
                Text(
                    text = loadError,
                    color = MaterialTheme.colorScheme.error,
                    style = MaterialTheme.typography.bodyMedium
                )
            } else {
                LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    if (pathStack.size > 1) {
                        item {
                            Card(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable { pathStack.removeAt(pathStack.lastIndex) },
                                colors = CardDefaults.cardColors(
                                    containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.6f)
                                )
                            ) {
                                Row(
                                    modifier = Modifier.padding(16.dp),
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        Icons.Default.FolderOpen,
                                        contentDescription = null,
                                        tint = MaterialTheme.colorScheme.primary
                                    )
                                    Spacer(modifier = Modifier.padding(8.dp))
                                    Text("..", style = MaterialTheme.typography.titleMedium)
                                }
                            }
                        }
                    }
                    items(entries) { entry ->
                        Card(
                            modifier = Modifier
                                .fillMaxWidth()
                                .clickable {
                                    when {
                                        entry.isDirectory -> pathStack.add(entry.file.absolutePath)
                                        mode == FilePickerMode.FILE -> onSelect(entry.file.absolutePath)
                                    }
                                },
                            colors = CardDefaults.cardColors(
                                containerColor = MaterialTheme.colorScheme.surfaceVariant
                            )
                        ) {
                            Row(
                                modifier = Modifier.padding(16.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    imageVector = if (entry.isDirectory) Icons.Default.Folder else Icons.Default.Description,
                                    contentDescription = null,
                                    tint = if (entry.isDirectory) MaterialTheme.colorScheme.primary
                                    else MaterialTheme.colorScheme.onSurfaceVariant
                                )
                                Spacer(modifier = Modifier.padding(8.dp))
                                Column(modifier = Modifier.weight(1f)) {
                                    Text(
                                        text = entry.file.name,
                                        style = MaterialTheme.typography.titleMedium,
                                        color = MaterialTheme.colorScheme.onSurface
                                    )
                                    if (!entry.isDirectory && entry.file.length() > 0) {
                                        Text(
                                            text = formatFileSize(entry.file.length()),
                                            style = MaterialTheme.typography.bodySmall,
                                            color = MaterialTheme.colorScheme.onSurfaceVariant
                                        )
                                    }
                                }
                            }
                        }
                    }
                }
                if (mode == FilePickerMode.DIRECTORY) {
                    Spacer(modifier = Modifier.height(16.dp))
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable { onSelect(currentPath) },
                        colors = CardDefaults.cardColors(
                            containerColor = MaterialTheme.colorScheme.primaryContainer
                        )
                    ) {
                        Row(
                            modifier = Modifier.padding(16.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(
                                Icons.Default.Folder,
                                contentDescription = null,
                                tint = MaterialTheme.colorScheme.onPrimaryContainer
                            )
                            Spacer(modifier = Modifier.padding(8.dp))
                            Text(
                                "Select this folder",
                                style = MaterialTheme.typography.titleMedium,
                                color = MaterialTheme.colorScheme.onPrimaryContainer
                            )
                        }
                    }
                }
            }
        }
    }
}

private data class DirEntry(val file: File, val isDirectory: Boolean)

private fun loadDirectoryEntries(path: String): Pair<List<DirEntry>, String?> {
    return try {
        val dir = File(path)
        if (!dir.exists() || !dir.isDirectory) {
            emptyList<DirEntry>() to "Not a directory"
        } else {
            val files = dir.listFiles()?.filter { it.canRead() } ?: emptyList()
            val dirs = files.filter { it.isDirectory }.sortedBy { it.name.lowercase() }
            val regularFiles = files.filter { it.isFile }.sortedBy { it.name.lowercase() }
            val entries = dirs.map { DirEntry(it, true) } + regularFiles.map { DirEntry(it, false) }
            entries to null
        }
    } catch (e: SecurityException) {
        emptyList<DirEntry>() to "Permission denied"
    } catch (e: Exception) {
        emptyList<DirEntry>() to (e.message ?: "Error loading directory")
    }
}

private fun getInitialPath(): String {
    return if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
        Environment.getExternalStorageDirectory().absolutePath
    } else {
        @Suppress("DEPRECATION")
        Environment.getExternalStorageDirectory().absolutePath
    }.ifBlank { "/storage/emulated/0" }
}

private fun formatFileSize(bytes: Long): String = when {
    bytes < 1024 -> "$bytes B"
    bytes < 1024 * 1024 -> "${bytes / 1024} KB"
    bytes < 1024 * 1024 * 1024 -> "${bytes / (1024 * 1024)} MB"
    else -> "${bytes / (1024 * 1024 * 1024)} GB"
}
