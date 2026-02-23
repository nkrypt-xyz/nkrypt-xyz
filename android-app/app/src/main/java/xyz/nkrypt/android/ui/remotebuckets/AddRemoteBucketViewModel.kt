package xyz.nkrypt.android.ui.remotebuckets

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import xyz.nkrypt.android.data.remote.api.BucketDto
import javax.inject.Inject

enum class AddStep {
    CREDENTIALS,
    SELECT_BUCKETS
}

data class AddRemoteBucketState(
    val step: AddStep = AddStep.CREDENTIALS,
    val buckets: List<BucketDto> = emptyList(),
    val selectedBucketIds: Set<String> = emptySet(),
    val bucketPasswords: Map<String, String> = emptyMap(),
    val serverUrl: String = "",
    val username: String = "",
    val password: String = "",
    val apiKey: String = "",
    val isLoading: Boolean = false,
    val isAdding: Boolean = false,
    val completed: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class AddRemoteBucketViewModel @Inject constructor(
    private val repository: RemoteBucketRepository,
    private val masterPasswordStore: MasterPasswordStore
) : ViewModel() {

    private val _state = MutableStateFlow(AddRemoteBucketState())
    val state: StateFlow<AddRemoteBucketState> = _state.asStateFlow()

    fun fetchBuckets(serverUrl: String, username: String, password: String) {
        viewModelScope.launch {
            _state.update {
                it.copy(isLoading = true, error = null)
            }
            try {
                val apiKey = repository.login(serverUrl, username, password)
                val buckets = repository.listBuckets(serverUrl, apiKey)
                _state.update {
                    it.copy(
                        step = AddStep.SELECT_BUCKETS,
                        buckets = buckets,
                        serverUrl = serverUrl,
                        username = username,
                        password = password,
                        apiKey = apiKey,
                        isLoading = false,
                        error = null
                    )
                }
            } catch (e: Exception) {
                _state.update {
                    it.copy(
                        isLoading = false,
                        error = e.message ?: "Failed to connect"
                    )
                }
            }
        }
    }

    fun toggleBucketSelection(bucketId: String) {
        _state.update {
            val newSelected = if (it.selectedBucketIds.contains(bucketId)) {
                it.selectedBucketIds - bucketId
            } else {
                it.selectedBucketIds + bucketId
            }
            it.copy(selectedBucketIds = newSelected)
        }
    }

    fun setBucketPassword(bucketId: String, password: String) {
        _state.update {
            it.copy(
                bucketPasswords = it.bucketPasswords + (bucketId to password)
            )
        }
    }

    fun addSelectedBuckets() {
        viewModelScope.launch {
            val s = _state.value
            val masterPassword = masterPasswordStore.getMasterPassword() ?: return@launch
            _state.update { it.copy(isAdding = true, error = null) }
            try {
                for (bucketId in s.selectedBucketIds) {
                    val bucket = s.buckets.find { it._id == bucketId } ?: continue
                    val encPassword = s.bucketPasswords[bucketId]
                        ?: throw Exception("Encryption password required for ${bucket.name}")
                    if (encPassword.length < 8) {
                        throw Exception("Encryption password must be at least 8 characters for ${bucket.name}")
                    }
                    val rootDirId = bucket.rootDirectoryId
                        ?: throw Exception("Bucket ${bucket.name} has no root directory")
                    repository.addRemoteBucket(
                        serverUrl = s.serverUrl,
                        username = s.username,
                        password = s.password,
                        bucketId = bucket._id,
                        bucketName = bucket.name,
                        rootDirectoryId = rootDirId,
                        encryptionPassword = encPassword,
                        masterPassword = masterPassword
                    )
                }
                _state.update { it.copy(isAdding = false, completed = true) }
            } catch (e: Exception) {
                _state.update {
                    it.copy(
                        isAdding = false,
                        error = e.message ?: "Failed to add buckets"
                    )
                }
            }
        }
    }
}
