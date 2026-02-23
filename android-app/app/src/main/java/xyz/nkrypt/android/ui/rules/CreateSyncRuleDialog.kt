package xyz.nkrypt.android.ui.rules

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity
import xyz.nkrypt.android.data.local.entity.SyncPostAction

@Composable
fun CreateSyncRuleDialog(
    localBuckets: List<LocalBucketEntity>,
    remoteBuckets: List<RemoteBucketEntity>,
    editRule: xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity?,
    onDismiss: () -> Unit,
    onCreate: (
        name: String,
        sourceBucketId: String,
        sourceDirectoryId: String?,
        targetRemoteBucketId: String,
        targetDirectoryId: String?,
        postAction: SyncPostAction
    ) -> Unit,
    onUpdate: ((id: String, name: String, sourceBucketId: String, sourceDirectoryId: String?, targetRemoteBucketId: String, targetDirectoryId: String?, postAction: SyncPostAction) -> Unit)? = null
) {
    var name by remember(editRule) { mutableStateOf(editRule?.name ?: "") }
    var error by remember { mutableStateOf<String?>(null) }
    var postAction by remember(editRule) { mutableStateOf(editRule?.let { SyncPostAction.valueOf(it.postAction) } ?: SyncPostAction.KEEP) }
    var selectedLocalBucket by remember(editRule, localBuckets) {
        mutableStateOf(localBuckets.find { it.id == editRule?.sourceBucketId })
    }
    var selectedRemoteBucket by remember(editRule, remoteBuckets) {
        mutableStateOf(remoteBuckets.find { it.id == editRule?.targetRemoteBucketId })
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text(
                text = if (editRule != null) "Edit Auto-Sync Rule" else "Create Auto-Sync Rule",
                style = MaterialTheme.typography.titleLarge,
                fontWeight = FontWeight.SemiBold
            )
        },
        text = {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .verticalScroll(rememberScrollState())
            ) {
                OutlinedTextField(
                    value = name,
                    onValueChange = { name = it; error = null },
                    label = { Text("Rule name") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "Source (local bucket)",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Spacer(modifier = Modifier.height(8.dp))
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f)
                    ),
                    shape = MaterialTheme.shapes.medium
                ) {
                    Column(modifier = Modifier.padding(4.dp)) {
                        for (bucket in localBuckets) {
                            val selected = selectedLocalBucket?.id == bucket.id
                            Card(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(vertical = 2.dp)
                                    .clickable { selectedLocalBucket = bucket },
                                colors = CardDefaults.cardColors(
                                    containerColor = if (selected) MaterialTheme.colorScheme.primaryContainer
                                    else MaterialTheme.colorScheme.surface
                                ),
                                shape = MaterialTheme.shapes.small
                            ) {
                                Text(
                                    text = bucket.name,
                                    modifier = Modifier.padding(14.dp),
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = if (selected) MaterialTheme.colorScheme.onPrimaryContainer
                                    else MaterialTheme.colorScheme.onSurface
                                )
                            }
                        }
                    }
                }
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "Target (remote bucket)",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Spacer(modifier = Modifier.height(8.dp))
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f)
                    ),
                    shape = MaterialTheme.shapes.medium
                ) {
                    Column(modifier = Modifier.padding(4.dp)) {
                        for (bucket in remoteBuckets) {
                            val selected = selectedRemoteBucket?.id == bucket.id
                            Card(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(vertical = 2.dp)
                                    .clickable { selectedRemoteBucket = bucket },
                                colors = CardDefaults.cardColors(
                                    containerColor = if (selected) MaterialTheme.colorScheme.primaryContainer
                                    else MaterialTheme.colorScheme.surface
                                ),
                                shape = MaterialTheme.shapes.small
                            ) {
                                Text(
                                    text = bucket.bucketName,
                                    modifier = Modifier.padding(14.dp),
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = if (selected) MaterialTheme.colorScheme.onPrimaryContainer
                                    else MaterialTheme.colorScheme.onSurface
                                )
                            }
                        }
                    }
                }
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "After sync",
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
                Spacer(modifier = Modifier.height(8.dp))
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f)
                    ),
                    shape = MaterialTheme.shapes.medium
                ) {
                    Column(modifier = Modifier.padding(4.dp)) {
                        for (action in SyncPostAction.entries) {
                            val isSelected = postAction == action
                            Card(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(vertical = 2.dp)
                                    .clickable { postAction = action },
                                colors = CardDefaults.cardColors(
                                    containerColor = if (isSelected) MaterialTheme.colorScheme.primaryContainer
                                    else MaterialTheme.colorScheme.surface
                                ),
                                shape = MaterialTheme.shapes.small
                            ) {
                                Text(
                                    text = when (action) {
                                        SyncPostAction.DELETE_LOCAL -> "Delete from local"
                                        SyncPostAction.KEEP -> "Keep in local"
                                    },
                                    modifier = Modifier.padding(14.dp),
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = if (isSelected) MaterialTheme.colorScheme.onPrimaryContainer
                                    else MaterialTheme.colorScheme.onSurface
                                )
                            }
                        }
                    }
                }
                if (error != null) {
                    Spacer(modifier = Modifier.height(12.dp))
                    Text(
                        text = error!!,
                        color = MaterialTheme.colorScheme.error,
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            }
        },
        confirmButton = {
            Button(
                onClick = {
                    when {
                        name.isBlank() -> error = "Enter a rule name"
                        selectedLocalBucket == null -> error = "Select source bucket"
                        selectedRemoteBucket == null -> error = "Select target bucket"
                        editRule != null && onUpdate != null -> {
                            onUpdate(
                                editRule.id,
                                name,
                                selectedLocalBucket!!.id,
                                editRule.sourceDirectoryId,
                                selectedRemoteBucket!!.id,
                                editRule.targetDirectoryId,
                                postAction
                            )
                            onDismiss()
                        }
                        else -> {
                            onCreate(
                                name,
                                selectedLocalBucket!!.id,
                                null,
                                selectedRemoteBucket!!.id,
                                null,
                                postAction
                            )
                            onDismiss()
                        }
                    }
                }
            ) {
                Text(if (editRule != null) "Save" else "Create")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Cancel")
            }
        }
    )
}
