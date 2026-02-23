package xyz.nkrypt.android.ui.localbuckets

import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.clickable
import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.combinedClickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
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
import androidx.compose.material.icons.filled.Download
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.DriveFileMove
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Upload
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.TextButton
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import androidx.hilt.navigation.compose.hiltViewModel
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity

private sealed class ContextMenuTarget {
    data class Dir(val dir: LocalDirectoryEntity) : ContextMenuTarget()
    data class File(val file: LocalFileEntity) : ContextMenuTarget()
}

@OptIn(ExperimentalMaterial3Api::class, ExperimentalFoundationApi::class)
@Composable
fun BrowseBucketScreen(
    bucketId: String,
    onBack: () -> Unit,
    onFileClick: (fileId: String) -> Unit = {},
    viewModel: BrowseBucketViewModel = hiltViewModel()
) {
    val state by viewModel.state.collectAsState()
    val showNewFolderDialog by viewModel.showNewFolderDialog.collectAsState(initial = false)
    val createFolderError by viewModel.createFolderError.collectAsState(initial = null)
    val renameError by viewModel.renameError.collectAsState(initial = null)
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var showMenu by remember { mutableStateOf(false) }
    var folderName by remember { mutableStateOf("") }
    var contextMenuTarget by remember { mutableStateOf<ContextMenuTarget?>(null) }
    var renameTarget by remember { mutableStateOf<ContextMenuTarget?>(null) }
    var showRenameDialog by remember { mutableStateOf(false) }
    var renameName by remember { mutableStateOf("") }
    var pendingDownloadTarget by remember { mutableStateOf<DownloadTarget?>(null) }

    val filePickerLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.GetContent()
    ) { uri: Uri? ->
        uri?.let {
            scope.launch {
                val fileName = context.contentResolver.query(uri, null, null, null, null)?.use { cursor ->
                    val nameIndex = cursor.getColumnIndex(android.provider.OpenableColumns.DISPLAY_NAME)
                    if (cursor.moveToFirst() && nameIndex >= 0) cursor.getString(nameIndex) else "file"
                } ?: "file"
                val content = withContext(Dispatchers.IO) {
                    context.contentResolver.openInputStream(uri)?.readBytes() ?: ByteArray(0)
                }
                viewModel.uploadFile(fileName, content)
            }
        }
    }

    val downloadDirPickerLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.OpenDocumentTree()
    ) { uri: Uri? ->
        uri?.let { treeUri ->
            val target = pendingDownloadTarget
            if (target != null) {
                scope.launch {
                    viewModel.downloadToDirectory(target, treeUri, context)
                    pendingDownloadTarget = null
                }
            }
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = state.bucket?.name ?: "Bucket",
                        style = MaterialTheme.typography.titleLarge
                    )
                },
                navigationIcon = {
                    IconButton(
                        onClick = {
                            if (viewModel.canNavigateUp()) {
                                viewModel.navigateUp()
                            } else {
                                onBack()
                            }
                        }
                    ) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        },
        floatingActionButton = {
            if (state.bucket != null && state.error == null) {
                Box {
                    androidx.compose.material3.FloatingActionButton(
                        onClick = { showMenu = !showMenu }
                    ) {
                        Icon(Icons.Default.Add, contentDescription = "Add")
                    }
                    DropdownMenu(
                        expanded = showMenu,
                        onDismissRequest = { showMenu = false }
                    ) {
                        DropdownMenuItem(
                            text = { Text("New folder") },
                            onClick = {
                                showMenu = false
                                scope.launch {
                                    delay(150)
                                    viewModel.showNewFolderDialog()
                                }
                            },
                            leadingIcon = { Icon(Icons.Default.Folder, contentDescription = null) }
                        )
                        DropdownMenuItem(
                            text = { Text("Upload file") },
                            onClick = {
                                showMenu = false
                                filePickerLauncher.launch("*/*")
                            },
                            leadingIcon = { Icon(Icons.Default.Upload, contentDescription = null) }
                        )
                    }
                }
            }
        }
    ) { paddingValues ->
        when {
            state.isLoading -> {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    CircularProgressIndicator()
                }
            }
            state.error != null -> {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues)
                        .padding(16.dp),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Text(
                        text = state.error!!,
                        color = MaterialTheme.colorScheme.error
                    )
                }
            }
            else -> {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues)
                        .padding(16.dp)
                ) {
                    LazyColumn(
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        items(state.directories) { dir ->
                            DirectoryItem(
                                directory = dir,
                                onClick = { viewModel.navigateInto(dir) },
                                onLongClick = { contextMenuTarget = ContextMenuTarget.Dir(dir) }
                            )
                        }
                        items(state.files) { file ->
                            FileItem(
                                file = file,
                                onClick = { onFileClick(file.id) },
                                onLongClick = { contextMenuTarget = ContextMenuTarget.File(file) }
                            )
                        }
                    }
                }
            }
        }
        if (contextMenuTarget != null) {
            val target = contextMenuTarget!!
            AlertDialog(
                onDismissRequest = { contextMenuTarget = null },
                content = {
                    Column {
                        Text("Actions", style = MaterialTheme.typography.titleMedium)
                        Spacer(modifier = Modifier.height(16.dp))
                        TextButton(
                            onClick = {
                                pendingDownloadTarget = when (target) {
                                    is ContextMenuTarget.Dir -> DownloadTarget.Directory(target.dir)
                                    is ContextMenuTarget.File -> DownloadTarget.File(target.file)
                                }
                                contextMenuTarget = null
                                downloadDirPickerLauncher.launch(null)
                            }
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Icon(Icons.Default.Download, contentDescription = null, modifier = Modifier.padding(end = 8.dp))
                                Text("Download")
                            }
                        }
                        TextButton(
                            onClick = {
                                renameName = when (target) {
                                    is ContextMenuTarget.Dir -> target.dir.name
                                    is ContextMenuTarget.File -> target.file.name
                                }
                                renameTarget = target
                                contextMenuTarget = null
                                showRenameDialog = true
                            }
                        ) { Text("Rename") }
                        TextButton(
                            onClick = {
                                scope.launch {
                                    when (target) {
                                        is ContextMenuTarget.Dir -> viewModel.deleteDirectory(target.dir.id)
                                        is ContextMenuTarget.File -> viewModel.deleteFile(target.file.id)
                                    }
                                    contextMenuTarget = null
                                }
                            }
                        ) { Text("Delete", color = MaterialTheme.colorScheme.error) }
                        TextButton(onClick = { contextMenuTarget = null }) { Text("Cancel") }
                    }
                }
            )
        }
        if (showRenameDialog && renameTarget != null) {
            val target = renameTarget!!
            AlertDialog(
                onDismissRequest = {
                    showRenameDialog = false
                    renameTarget = null
                    viewModel.clearRenameError()
                },
                content = {
                    Column {
                        OutlinedTextField(
                            value = renameName,
                            onValueChange = { renameName = it },
                            label = { Text("Name") },
                            singleLine = true
                        )
                        if (renameError != null) {
                            Text(
                                text = renameError!!,
                                color = MaterialTheme.colorScheme.error,
                                style = MaterialTheme.typography.bodySmall,
                                modifier = Modifier.padding(top = 8.dp)
                            )
                        }
                        Row(modifier = Modifier.padding(top = 16.dp)) {
                            TextButton(
                                onClick = {
                                    if (renameName.isNotBlank()) {
                                        scope.launch {
                                            val success = when (target) {
                                                is ContextMenuTarget.Dir -> viewModel.renameDirectory(target.dir.id, renameName)
                                                is ContextMenuTarget.File -> viewModel.renameFile(target.file.id, renameName)
                                            }
                                            if (success) {
                                                showRenameDialog = false
                                                renameTarget = null
                                            }
                                        }
                                    }
                                }
                            ) { Text("Rename") }
                            TextButton(onClick = {
                                showRenameDialog = false
                                renameTarget = null
                                viewModel.clearRenameError()
                            }) { Text("Cancel") }
                        }
                    }
                }
            )
        }
        if (showNewFolderDialog) {
            AlertDialog(
                onDismissRequest = { viewModel.dismissNewFolderDialog() },
                content = {
                    Column {
                        OutlinedTextField(
                            value = folderName,
                            onValueChange = { folderName = it },
                            label = { Text("Folder name") },
                            singleLine = true
                        )
                        if (createFolderError != null) {
                            Text(
                                text = createFolderError!!,
                                color = MaterialTheme.colorScheme.error,
                                style = MaterialTheme.typography.bodySmall,
                                modifier = Modifier.padding(top = 8.dp)
                            )
                        }
                        Row(modifier = Modifier.padding(top = 16.dp)) {
                            TextButton(
                                onClick = {
                                    if (folderName.isNotBlank()) {
                                        scope.launch {
                                            viewModel.createFolder(folderName)
                                            folderName = ""
                                        }
                                    }
                                }
                            ) { Text("Create") }
                            TextButton(onClick = { viewModel.dismissNewFolderDialog() }) { Text("Cancel") }
                        }
                    }
                }
            )
        }
    }
}

@Composable
private fun DirectoryItem(
    directory: LocalDirectoryEntity,
    onClick: () -> Unit,
    onLongClick: () -> Unit = {}
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .combinedClickable(onClick = onClick, onLongClick = onLongClick),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                Icons.Default.Folder,
                contentDescription = null,
                tint = MaterialTheme.colorScheme.primary
            )
            Spacer(modifier = Modifier.padding(8.dp))
            Text(
                text = directory.name,
                style = MaterialTheme.typography.titleMedium,
                color = MaterialTheme.colorScheme.onSurface
            )
        }
    }
}

@Composable
private fun FileItem(
    file: LocalFileEntity,
    onClick: () -> Unit,
    onLongClick: () -> Unit = {}
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .combinedClickable(onClick = onClick, onLongClick = onLongClick),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                Icons.Default.Description,
                contentDescription = null,
                tint = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(modifier = Modifier.padding(8.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = file.name,
                    style = MaterialTheme.typography.titleMedium,
                    color = MaterialTheme.colorScheme.onSurface
                )
                Text(
                    text = formatFileSize(file.sizeInBytes),
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        }
    }
}

private fun formatFileSize(bytes: Long): String {
    return when {
        bytes < 1024 -> "$bytes B"
        bytes < 1024 * 1024 -> "${bytes / 1024} KB"
        bytes < 1024 * 1024 * 1024 -> "${bytes / (1024 * 1024)} MB"
        else -> "${bytes / (1024 * 1024 * 1024)} GB"
    }
}
