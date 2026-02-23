package xyz.nkrypt.android.ui.localbuckets

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
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Download
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.Info
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import kotlinx.coroutines.launch
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.ui.filepicker.FilePickerMode
import xyz.nkrypt.android.util.showSuccessToast
import xyz.nkrypt.android.ui.filepicker.FilePickerScreen

@OptIn(ExperimentalFoundationApi::class)
@Composable
fun LocalBucketsScreen(
    viewModel: LocalBucketsViewModel,
    onBucketClick: (LocalBucketEntity) -> Unit
) {
    val buckets by viewModel.buckets.collectAsState(initial = emptyList())
    val createDialogState by viewModel.createDialogState.collectAsState(initial = null)
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var contextMenuBucket by remember { mutableStateOf<LocalBucketEntity?>(null) }
    var metadataBucket by remember { mutableStateOf<LocalBucketEntity?>(null) }
    var bucketToDelete by remember { mutableStateOf<LocalBucketEntity?>(null) }
    var bucketToDownload by remember { mutableStateOf<LocalBucketEntity?>(null) }
    var showCreateDirPicker by remember { mutableStateOf(false) }
    var showDownloadDirPicker by remember { mutableStateOf(false) }

    when {
        showCreateDirPicker -> {
            FilePickerScreen(
                mode = FilePickerMode.DIRECTORY,
                onSelect = { path ->
                    viewModel.onDirectorySelectedPath(path)
                    showCreateDirPicker = false
                },
                onDismiss = { showCreateDirPicker = false },
                message = "Select where to store the local bucket"
            )
        }
        showDownloadDirPicker && bucketToDownload != null -> {
            val bucket = bucketToDownload!!
            FilePickerScreen(
                mode = FilePickerMode.DIRECTORY,
                onSelect = { path ->
                    scope.launch {
                        viewModel.downloadBucketToDirectory(bucket, path)
                        showSuccessToast(context, "Download completed.")
                        bucketToDownload = null
                        showDownloadDirPicker = false
                    }
                },
                onDismiss = {
                    bucketToDownload = null
                    showDownloadDirPicker = false
                },
                message = "Select where to download the bucket"
            )
        }
        else -> {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(
                onClick = { viewModel.showCreateDialog() }
            ) {
                Icon(Icons.Default.Add, contentDescription = "Add bucket")
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp)
        ) {
            Text(
                text = "Local Buckets",
                style = MaterialTheme.typography.headlineMedium,
                color = MaterialTheme.colorScheme.onBackground
            )
            Spacer(modifier = Modifier.height(16.dp))

            if (buckets.isEmpty()) {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .weight(1f),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(
                            Icons.Default.Folder,
                            contentDescription = null,
                            modifier = Modifier.padding(16.dp),
                            tint = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                        Text(
                            text = "No local buckets yet",
                            style = MaterialTheme.typography.bodyLarge,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                        Text(
                            text = "Tap + to create one",
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.onSurfaceVariant.copy(alpha = 0.7f)
                        )
                    }
                }
            } else {
                LazyColumn(
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(buckets) { bucket ->
                        Card(
                            modifier = Modifier
                                .fillMaxWidth()
                                .combinedClickable(
                                    onClick = { onBucketClick(bucket) },
                                    onLongClick = { contextMenuBucket = bucket }
                                ),
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
                                Column(modifier = Modifier.weight(1f)) {
                                    Text(
                                        text = bucket.name,
                                        style = MaterialTheme.typography.titleMedium,
                                        color = MaterialTheme.colorScheme.onSurface
                                    )
                                    Text(
                                        text = bucket.rootPath,
                                        style = MaterialTheme.typography.bodySmall,
                                        color = MaterialTheme.colorScheme.onSurfaceVariant
                                    )
                                }
                            }
                        }
                    }
                }
            }
        }
        if (contextMenuBucket != null) {
            val bucket = contextMenuBucket!!
            AlertDialog(
                onDismissRequest = { contextMenuBucket = null },
                content = {
                    Column {
                        Text("${bucket.name}", style = MaterialTheme.typography.titleMedium)
                        Spacer(modifier = Modifier.height(16.dp))
                        TextButton(
                            onClick = {
                                metadataBucket = bucket
                                contextMenuBucket = null
                            }
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Icon(Icons.Default.Info, contentDescription = null, modifier = Modifier.padding(end = 8.dp))
                                Text("View metadata")
                            }
                        }
                        TextButton(
                            onClick = {
                                bucketToDownload = bucket
                                contextMenuBucket = null
                                showDownloadDirPicker = true
                            }
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Icon(Icons.Default.Download, contentDescription = null, modifier = Modifier.padding(end = 8.dp))
                                Text("Download entire bucket")
                            }
                        }
                        TextButton(
                            onClick = {
                                bucketToDelete = bucket
                                contextMenuBucket = null
                            }
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Icon(Icons.Default.Delete, contentDescription = null, modifier = Modifier.padding(end = 8.dp))
                                Text("Delete", color = MaterialTheme.colorScheme.error)
                            }
                        }
                        TextButton(onClick = { contextMenuBucket = null }) { Text("Cancel") }
                    }
                }
            )
        }
        if (metadataBucket != null) {
            BucketMetadataDialog(
                bucket = metadataBucket!!,
                onDismiss = { metadataBucket = null }
            )
        }
        if (bucketToDelete != null) {
            val bucket = bucketToDelete!!
            AlertDialog(
                onDismissRequest = { bucketToDelete = null },
                content = {
                    Column {
                        Text("Delete \"${bucket.name}\"?")
                        Spacer(modifier = Modifier.height(16.dp))
                        TextButton(
                            onClick = {
                                viewModel.deleteBucket(bucket.id, deleteFiles = false)
                                bucketToDelete = null
                            }
                        ) { Text("Metadata only") }
                        TextButton(
                            onClick = {
                                viewModel.deleteBucket(bucket.id, deleteFiles = true)
                                bucketToDelete = null
                            }
                        ) { Text("Delete with files", color = MaterialTheme.colorScheme.error) }
                        TextButton(onClick = { bucketToDelete = null }) { Text("Cancel") }
                    }
                }
            )
        }
    }

    createDialogState?.let { state ->
        CreateBucketDialog(
            state = state,
            onDismiss = { viewModel.dismissCreateDialog() },
            onSelectLocation = { showCreateDirPicker = true },
            onCreate = { name, password ->
                viewModel.createBucket(name, password)
            }
        )
    }
        }
    }
}
