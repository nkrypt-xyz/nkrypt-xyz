package xyz.nkrypt.android.ui.rules

import android.content.Intent
import android.net.Uri
import android.provider.DocumentsContract
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
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
import androidx.compose.material.icons.filled.Cloud
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity
import xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity

@Composable
fun RulesScreen(
    viewModel: RulesViewModel = hiltViewModel(),
    onNavigateToProgress: () -> Unit = {},
    onNavigateToLocalBuckets: () -> Unit = {},
    onNavigateToRemoteBuckets: () -> Unit = {}
) {
    val context = LocalContext.current
    val importRules by viewModel.importRules.collectAsState()
    val syncRules by viewModel.syncRules.collectAsState()
    val localBuckets by viewModel.localBuckets.collectAsState()
    val remoteBuckets by viewModel.remoteBuckets.collectAsState()
    val showImportDialog by viewModel.showImportDialog.collectAsState()
    val showSyncDialog by viewModel.showSyncDialog.collectAsState()
    val importDialogSourcePath by viewModel.importDialogSourcePath.collectAsState()
    val editImportRule by viewModel.editImportRule.collectAsState()
    val editSyncRule by viewModel.editSyncRule.collectAsState()

    var showFabMenu by remember { mutableStateOf(false) }
    val progressState by viewModel.progressState.collectAsState()

    val directoryPickerLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.OpenDocumentTree()
    ) { uri: Uri? ->
        if (uri != null) {
            context.contentResolver.takePersistableUriPermission(
                uri,
                Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_GRANT_WRITE_URI_PERMISSION
            )
            val path = getPathFromUri(uri)
            viewModel.onImportSourcePathSelected(path)
        }
    }

    Scaffold(
        floatingActionButton = {
            if (localBuckets.isNotEmpty()) {
                Column(horizontalAlignment = Alignment.End) {
                    if (showFabMenu) {
                        DropdownMenu(
                            expanded = showFabMenu,
                            onDismissRequest = { showFabMenu = false }
                        ) {
                            DropdownMenuItem(
                                text = { Text("Auto-import rule") },
                                onClick = {
                                    showFabMenu = false
                                    viewModel.showCreateImportDialog()
                                }
                            )
                            DropdownMenuItem(
                                text = { Text("Auto-sync rule") },
                                onClick = {
                                    showFabMenu = false
                                    if (remoteBuckets.isNotEmpty()) {
                                        viewModel.showCreateSyncDialog()
                                    }
                                }
                            )
                        }
                    }
                    FloatingActionButton(
                        onClick = { showFabMenu = !showFabMenu }
                    ) {
                        Icon(Icons.Default.Add, contentDescription = "Add rule")
                    }
                }
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
                text = "Rules",
                style = MaterialTheme.typography.headlineMedium,
                color = MaterialTheme.colorScheme.onBackground
            )
            Spacer(modifier = Modifier.height(16.dp))

            Text(
                text = "Auto-import",
                style = MaterialTheme.typography.titleMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(modifier = Modifier.height(8.dp))
            if (localBuckets.isEmpty()) {
                RulesEmptyStateCard(
                    icon = Icons.Default.Folder,
                    message = "You need at least one local bucket to create auto-import rules.",
                    buttonText = "Go to Local Buckets",
                    onButtonClick = onNavigateToLocalBuckets
                )
            } else if (importRules.isEmpty()) {
                Text(
                    text = "No import rules. Tap + to add.",
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            } else {
                var ruleToDelete by remember { mutableStateOf<AutoImportRuleEntity?>(null) }
                var ruleMenuExpanded by remember { mutableStateOf<AutoImportRuleEntity?>(null) }
                LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    items(importRules) { rule ->
                        RuleCard(
                            title = rule.name,
                            subtitle = rule.sourceDirectoryPath,
                            onSyncClick = {
                                viewModel.runImport(rule.id) {}
                                onNavigateToProgress()
                            },
                            onEditClick = { viewModel.showEditImportDialog(rule) },
                            onDeleteClick = { ruleToDelete = rule },
                            menuExpanded = ruleMenuExpanded == rule,
                            onMenuClick = { ruleMenuExpanded = if (ruleMenuExpanded == rule) null else rule },
                            onDismissMenu = { ruleMenuExpanded = null }
                        )
                    }
                }
                ruleToDelete?.let { rule ->
                    AlertDialog(
                        onDismissRequest = { ruleToDelete = null },
                        title = { Text("Delete rule?") },
                        text = { Text("Delete \"${rule.name}\"?") },
                        confirmButton = {
                            TextButton(
                                onClick = {
                                    viewModel.deleteImportRule(rule.id)
                                    ruleToDelete = null
                                }
                            ) { Text("Delete", color = MaterialTheme.colorScheme.error) }
                        },
                        dismissButton = {
                            TextButton(onClick = { ruleToDelete = null }) { Text("Cancel") }
                        }
                    )
                }
            }

            Spacer(modifier = Modifier.height(24.dp))
            Text(
                text = "Auto-sync",
                style = MaterialTheme.typography.titleMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(modifier = Modifier.height(8.dp))
            if (localBuckets.isEmpty() || remoteBuckets.isEmpty()) {
                val message = when {
                    localBuckets.isEmpty() && remoteBuckets.isEmpty() ->
                        "You need at least one local bucket and one remote bucket to create sync rules."
                    localBuckets.isEmpty() ->
                        "You need at least one local bucket to create sync rules."
                    else ->
                        "You need at least one remote bucket to create sync rules."
                }
                RulesEmptyStateCard(
                    icon = Icons.Default.Cloud,
                    message = message,
                    buttonText = if (localBuckets.isEmpty()) "Go to Local Buckets" else "Go to Remote Buckets",
                    onButtonClick = if (localBuckets.isEmpty()) onNavigateToLocalBuckets else onNavigateToRemoteBuckets,
                    secondaryButtonText = if (localBuckets.isEmpty() && remoteBuckets.isEmpty()) "Go to Remote Buckets" else null,
                    onSecondaryButtonClick = if (localBuckets.isEmpty() && remoteBuckets.isEmpty()) onNavigateToRemoteBuckets else null
                )
            } else if (syncRules.isEmpty()) {
                Text(
                    text = "No sync rules. Tap + to add.",
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            } else {
                var syncRuleToDelete by remember { mutableStateOf<AutoSyncRuleEntity?>(null) }
                var syncRuleMenuExpanded by remember { mutableStateOf<AutoSyncRuleEntity?>(null) }
                LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    items(syncRules) { rule ->
                        RuleCard(
                            title = rule.name,
                            subtitle = "Local → Remote",
                            onSyncClick = {
                                viewModel.runSync(rule.id) {}
                                onNavigateToProgress()
                            },
                            onEditClick = { viewModel.showEditSyncDialog(rule) },
                            onDeleteClick = { syncRuleToDelete = rule },
                            menuExpanded = syncRuleMenuExpanded == rule,
                            onMenuClick = { syncRuleMenuExpanded = if (syncRuleMenuExpanded == rule) null else rule },
                            onDismissMenu = { syncRuleMenuExpanded = null }
                        )
                    }
                }
                syncRuleToDelete?.let { rule ->
                    AlertDialog(
                        onDismissRequest = { syncRuleToDelete = null },
                        title = { Text("Delete rule?") },
                        text = { Text("Delete \"${rule.name}\"?") },
                        confirmButton = {
                            TextButton(
                                onClick = {
                                    viewModel.deleteSyncRule(rule.id)
                                    syncRuleToDelete = null
                                }
                            ) { Text("Delete", color = MaterialTheme.colorScheme.error) }
                        },
                        dismissButton = {
                            TextButton(onClick = { syncRuleToDelete = null }) { Text("Cancel") }
                        }
                    )
                }
            }
        }
    }

    if (showImportDialog) {
        CreateImportRuleDialog(
            localBuckets = localBuckets,
            selectedSourcePath = importDialogSourcePath,
            editRule = editImportRule,
            onDismiss = { viewModel.dismissImportDialog() },
            onSelectSource = { directoryPickerLauncher.launch(null) },
            onCreate = { name, path, bucketId, postAction ->
                viewModel.createImportRule(name, path, bucketId, postAction)
            },
            onUpdate = { id, name, path, bucketId, postAction ->
                viewModel.updateImportRule(id, name, path, bucketId, postAction)
            }
        )
    }

    if (showSyncDialog) {
        CreateSyncRuleDialog(
            localBuckets = localBuckets,
            remoteBuckets = remoteBuckets,
            editRule = editSyncRule,
            onDismiss = { viewModel.dismissSyncDialog() },
            onCreate = { name, srcBucket, srcDir, tgtBucket, tgtDir, postAction ->
                viewModel.createSyncRule(name, srcBucket, srcDir, tgtBucket, tgtDir, postAction)
            },
            onUpdate = { id, name, srcBucket, srcDir, tgtBucket, tgtDir, postAction ->
                viewModel.updateSyncRule(id, name, srcBucket, srcDir, tgtBucket, tgtDir, postAction)
            }
        )
    }

}

@Composable
private fun RulesEmptyStateCard(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    message: String,
    buttonText: String,
    onButtonClick: () -> Unit,
    secondaryButtonText: String? = null,
    onSecondaryButtonClick: (() -> Unit)? = null
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.6f))
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(20.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = MaterialTheme.colorScheme.primary
            )
            if (message.isNotEmpty()) {
                Spacer(modifier = Modifier.height(12.dp))
                Text(
                    text = message,
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
            Spacer(modifier = Modifier.height(16.dp))
            Button(onClick = onButtonClick) {
                Text(buttonText)
            }
            if (secondaryButtonText != null && onSecondaryButtonClick != null) {
                Spacer(modifier = Modifier.height(8.dp))
                TextButton(onClick = onSecondaryButtonClick) {
                    Text(secondaryButtonText)
                }
            }
        }
    }
}

@Composable
private fun RuleCard(
    title: String,
    subtitle: String,
    onSyncClick: () -> Unit,
    onEditClick: () -> Unit,
    onDeleteClick: () -> Unit,
    menuExpanded: Boolean,
    onMenuClick: () -> Unit,
    onDismissMenu: () -> Unit
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = Icons.Default.Folder,
                contentDescription = "Rule",
                tint = MaterialTheme.colorScheme.primary
            )
            Spacer(modifier = Modifier.padding(8.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = title,
                    style = MaterialTheme.typography.titleMedium,
                    color = MaterialTheme.colorScheme.onSurface
                )
                Text(
                    text = subtitle,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
            IconButton(onClick = onSyncClick) {
                Icon(
                    Icons.Default.PlayArrow,
                    contentDescription = "Sync now",
                    tint = MaterialTheme.colorScheme.primary
                )
            }
            Box {
                IconButton(onClick = onMenuClick) {
                    Icon(
                        Icons.Default.MoreVert,
                        contentDescription = "More options",
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
                DropdownMenu(
                    expanded = menuExpanded,
                    onDismissRequest = onDismissMenu
                ) {
                    DropdownMenuItem(
                        text = { Text("Edit") },
                        onClick = {
                            onDismissMenu()
                            onEditClick()
                        },
                        leadingIcon = {
                            Icon(Icons.Default.Edit, contentDescription = null)
                        }
                    )
                    DropdownMenuItem(
                        text = { Text("Delete", color = MaterialTheme.colorScheme.error) },
                        onClick = {
                            onDismissMenu()
                            onDeleteClick()
                        },
                        leadingIcon = {
                            Icon(Icons.Default.Delete, contentDescription = null, tint = MaterialTheme.colorScheme.error)
                        }
                    )
                }
            }
        }
    }
}

private fun getPathFromUri(uri: Uri): String {
    val docId = DocumentsContract.getTreeDocumentId(uri)
    val split = docId.split(":")
    return when {
        split.size >= 2 -> {
            val type = split[0]
            val path = split[1]
            when (type) {
                "primary" -> "/storage/emulated/0/$path"
                else -> "/storage/$type/$path"
            }
        }
        else -> uri.path ?: ""
    }
}
