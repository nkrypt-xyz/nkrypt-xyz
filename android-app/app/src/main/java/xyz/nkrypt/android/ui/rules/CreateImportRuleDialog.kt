package xyz.nkrypt.android.ui.rules

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
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
import androidx.compose.material3.OutlinedButton
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
import xyz.nkrypt.android.data.local.entity.ImportPostAction
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity

@Composable
fun CreateImportRuleDialog(
    localBuckets: List<LocalBucketEntity>,
    selectedSourcePath: String?,
    editRule: xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity?,
    onDismiss: () -> Unit,
    onSelectSource: () -> Unit,
    onCreate: (name: String, sourcePath: String, targetBucketId: String, postAction: ImportPostAction) -> Unit,
    onUpdate: ((id: String, name: String, sourcePath: String, targetBucketId: String, postAction: ImportPostAction) -> Unit)? = null
) {
    var name by remember(editRule) { mutableStateOf(editRule?.name ?: "") }
    var error by remember { mutableStateOf<String?>(null) }
    var postAction by remember(editRule) { mutableStateOf(editRule?.let { ImportPostAction.valueOf(it.postAction) } ?: ImportPostAction.KEEP) }
    var selectedBucket by remember(editRule, localBuckets) {
        mutableStateOf(localBuckets.find { it.id == editRule?.targetBucketId })
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text(
                text = if (editRule != null) "Edit Auto-Import Rule" else "Create Auto-Import Rule",
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
                Row(modifier = Modifier.fillMaxWidth()) {
                    OutlinedTextField(
                        value = selectedSourcePath ?: "",
                        onValueChange = {},
                        label = { Text("Source directory") },
                        modifier = Modifier.weight(1f),
                        readOnly = true
                    )
                    Spacer(modifier = Modifier.padding(8.dp))
                    OutlinedButton(onClick = onSelectSource) {
                        Text(if (selectedSourcePath == null) "Select" else "Change")
                    }
                }
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "Target bucket",
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
                            val selected = selectedBucket?.id == bucket.id
                            Card(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(vertical = 2.dp)
                                    .clickable { selectedBucket = bucket },
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
                    text = "After import",
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
                        for (action in ImportPostAction.entries) {
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
                                        ImportPostAction.DELETE -> "Delete original"
                                        ImportPostAction.TRASH -> "Move to trash"
                                        ImportPostAction.KEEP -> "Keep original"
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
                        selectedSourcePath == null -> error = "Select source directory"
                        selectedBucket == null -> error = "Select target bucket"
                        editRule != null && onUpdate != null -> {
                            onUpdate(editRule.id, name, selectedSourcePath!!, selectedBucket!!.id, postAction)
                            onDismiss()
                        }
                        else -> {
                            onCreate(name, selectedSourcePath!!, selectedBucket!!.id, postAction)
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
