package xyz.nkrypt.android.ui.remotebuckets

import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.combinedClickable
import androidx.compose.foundation.clickable
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
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Description
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.Upload
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
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
import androidx.hilt.navigation.compose.hiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import xyz.nkrypt.android.data.remote.api.DirectoryDto
import xyz.nkrypt.android.data.remote.api.FileDto

private sealed class RemoteContextMenuTarget {
    data class Dir(val dir: DirectoryDto) : RemoteContextMenuTarget()
    data class File(val file: FileDto) : RemoteContextMenuTarget()
}

@OptIn(ExperimentalMaterial3Api::class, ExperimentalFoundationApi::class)
@Composable
fun BrowseRemoteBucketScreen(
    bucketId: String,
    onBack: () -> Unit,
    onFileClick: (fileId: String) -> Unit,
    viewModel: BrowseRemoteBucketViewModel = hiltViewModel()
) {
    val state by viewModel.state.collectAsState()
    val showNewFolderDialog by viewModel.showNewFolderDialog.collectAsState(initial = false)
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var showFabMenu by remember { mutableStateOf(false) }
    var folderName by remember { mutableStateOf("") }
    var contextMenuTarget by remember { mutableStateOf<RemoteContextMenuTarget?>(null) }
    var renameTarget by remember { mutableStateOf<RemoteContextMenuTarget?>(null) }
    var renameName by remember { mutableStateOf("") }
    var pendingDownloadFileId by remember { mutableStateOf<String?>(null) }

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

    val downloadLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.CreateDocument("*/*")
    ) { uri: Uri? ->
        uri?.let { saveUri ->
            val fileId = pendingDownloadFileId
            if (fileId != null) {
                scope.launch {
                    val bytes = viewModel.downloadFileDecrypted(fileId)
                    bytes?.let {
                        context.contentResolver.openOutputStream(saveUri)?.use { out -> out.write(bytes) }
                    }
                    pendingDownloadFileId = null
                }
            }
        }
    }

    viewModel.loadBucket(bucketId)

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        text = state.bucket?.bucketName ?: "Bucket",
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
                    FloatingActionButton(onClick = { showFabMenu = !showFabMenu }) {
                        Icon(Icons.Default.Add, contentDescription = "Add")
                    }
                    DropdownMenu(
                        expanded = showFabMenu,
                        onDismissRequest = { showFabMenu = false }
                    ) {
                        DropdownMenuItem(
                            text = { Text("New folder") },
                            onClick = {
                                showFabMenu = false
                                scope.launch {
                                    delay(150)
                                    viewModel.showNewFolderDialog()
                                }
                            },
                            leadingIcon = { Icon(Icons.Default.Folder, null) }
                        )
                        DropdownMenuItem(
                            text = { Text("Upload file") },
                            onClick = { showFabMenu = false; filePickerLauncher.launch("*/*") },
                            leadingIcon = { Icon(Icons.Default.Upload, null) }
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
                            RemoteDirectoryItem(
                                directory = dir,
                                onClick = { viewModel.navigateInto(dir) },
                                onLongClick = { contextMenuTarget = RemoteContextMenuTarget.Dir(dir) }
                            )
                        }
                        items(state.files) { file ->
                            RemoteFileItem(
                                file = file,
                                onClick = { onFileClick(file._id) },
                                onLongClick = { contextMenuTarget = RemoteContextMenuTarget.File(file) }
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
                        when (target) {
                            is RemoteContextMenuTarget.File -> {
                                TextButton(
                                    onClick = {
                                        pendingDownloadFileId = target.file._id
                                        contextMenuTarget = null
                                        downloadLauncher.launch(target.file.name)
                                    }
                                ) { Text("Download to device") }
                            }
                            else -> {}
                        }
                        TextButton(
                            onClick = {
                                renameName = when (target) {
                                    is RemoteContextMenuTarget.Dir -> target.dir.name
                                    is RemoteContextMenuTarget.File -> target.file.name
                                }
                                renameTarget = target
                                contextMenuTarget = null
                            }
                        ) { Text("Rename") }
                        TextButton(
                            onClick = {
                                scope.launch {
                                    when (target) {
                                        is RemoteContextMenuTarget.Dir -> viewModel.deleteDirectory(target.dir._id)
                                        is RemoteContextMenuTarget.File -> viewModel.deleteFile(target.file._id)
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
        if (renameTarget != null) {
            val target = renameTarget!!
            AlertDialog(
                onDismissRequest = { renameTarget = null },
                content = {
                    Column {
                        OutlinedTextField(
                            value = renameName,
                            onValueChange = { renameName = it },
                            label = { Text("Name") },
                            singleLine = true
                        )
                        Row(modifier = Modifier.padding(top = 16.dp)) {
                            TextButton(
                                onClick = {
                                    if (renameName.isNotBlank()) {
                                        scope.launch {
                                            when (target) {
                                                is RemoteContextMenuTarget.Dir -> viewModel.renameDirectory(target.dir._id, renameName)
                                                is RemoteContextMenuTarget.File -> viewModel.renameFile(target.file._id, renameName)
                                            }
                                            renameTarget = null
                                        }
                                    }
                                }
                            ) { Text("Rename") }
                            TextButton(onClick = { renameTarget = null }) { Text("Cancel") }
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
                        Row(modifier = Modifier.padding(top = 16.dp)) {
                            TextButton(
                                onClick = {
                                    if (folderName.isNotBlank()) {
                                        scope.launch {
                                            viewModel.createDirectory(folderName)
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
private fun RemoteDirectoryItem(
    directory: DirectoryDto,
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
private fun RemoteFileItem(
    file: FileDto,
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
