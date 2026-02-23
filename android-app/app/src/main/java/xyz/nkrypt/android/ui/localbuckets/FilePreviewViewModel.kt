package xyz.nkrypt.android.ui.localbuckets

import android.graphics.BitmapFactory
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.LocalBucketRepository
import xyz.nkrypt.android.data.local.MasterPasswordStore
import javax.inject.Inject

data class FilePreviewState(
    val fileName: String? = null,
    val textContent: String? = null,
    val htmlContent: String? = null,
    val imageBitmap: android.graphics.Bitmap? = null,
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class FilePreviewViewModel @Inject constructor(
    private val repository: LocalBucketRepository,
    private val masterPasswordStore: MasterPasswordStore,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val _state = MutableStateFlow(FilePreviewState())
    val state: StateFlow<FilePreviewState> = _state.asStateFlow()

    private var loadedFileKey: String? = null

    init {
        val bucketId = savedStateHandle.get<String>("bucketId")
        val fileId = savedStateHandle.get<String>("fileId")
        if (bucketId != null && fileId != null) {
            loadFile(bucketId, fileId)
        }
    }

    fun loadFile(bucketId: String, fileId: String) {
        val key = "$bucketId/$fileId"
        if (loadedFileKey == key) return
        loadedFileKey = key
        viewModelScope.launch {
            _state.update { it.copy(isLoading = true, error = null) }
            try {
                val bucket = repository.getBucketById(bucketId) ?: run {
                    _state.update { it.copy(isLoading = false, error = "Bucket not found") }
                    return@launch
                }
                val masterPassword = masterPasswordStore.getMasterPassword() ?: run {
                    _state.update { it.copy(isLoading = false, error = "Master password required") }
                    return@launch
                }
                val bucketPassword = repository.decryptBucketPassword(bucket.cryptData, masterPassword)
                val file = repository.getFileById(fileId) ?: run {
                    _state.update { it.copy(isLoading = false, error = "File not found") }
                    return@launch
                }
                val content = repository.readFileContent(bucket, fileId, bucketPassword) ?: run {
                    _state.update { it.copy(isLoading = false, error = "Failed to read file") }
                    return@launch
                }
                val ext = file.name.substringAfterLast('.', "").lowercase()
                val isImage = ext in listOf("jpg", "jpeg", "png", "gif", "webp", "bmp")
                val isHtml = ext in listOf("html", "htm")
                val isText = ext in listOf("txt", "md", "json", "xml", "css", "js", "log", "csv")
                when {
                    isImage -> {
                        val bitmap = BitmapFactory.decodeByteArray(content, 0, content.size)
                        _state.update {
                            it.copy(
                                fileName = file.name,
                                imageBitmap = bitmap,
                                isLoading = false
                            )
                        }
                    }
                    isHtml -> {
                        val html = content.toString(Charsets.UTF_8)
                        _state.update {
                            it.copy(
                                fileName = file.name,
                                htmlContent = html,
                                isLoading = false
                            )
                        }
                    }
                    isText -> {
                        val text = content.toString(Charsets.UTF_8)
                        _state.update {
                            it.copy(
                                fileName = file.name,
                                textContent = text,
                                isLoading = false
                            )
                        }
                    }
                    else -> {
                        _state.update {
                            it.copy(
                                fileName = file.name,
                                textContent = "Binary file (${content.size} bytes). Preview not supported.",
                                isLoading = false
                            )
                        }
                    }
                }
            } catch (e: Exception) {
                _state.update {
                    it.copy(
                        isLoading = false,
                        error = e.message ?: "Failed to load"
                    )
                }
            }
        }
    }
}
