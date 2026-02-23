package xyz.nkrypt.android.ui.remotebuckets

import android.graphics.BitmapFactory
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.crypto.CryptoUtils
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import javax.inject.Inject

data class RemoteFilePreviewState(
    val fileName: String? = null,
    val textContent: String? = null,
    val htmlContent: String? = null,
    val imageBitmap: android.graphics.Bitmap? = null,
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class RemoteFilePreviewViewModel @Inject constructor(
    private val repository: RemoteBucketRepository,
    private val masterPasswordStore: MasterPasswordStore
) : ViewModel() {

    private val _state = MutableStateFlow(RemoteFilePreviewState())
    val state: StateFlow<RemoteFilePreviewState> = _state.asStateFlow()

    private var loadedFileKey: String? = null

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
                val bucketPassword = repository.decryptEncryptionPassword(bucket.encryptionPasswordEncrypted, masterPassword)
                val apiKey = repository.getApiKey(bucket)
                val file = repository.getFileMetadata(bucket, fileId, apiKey) ?: run {
                    _state.update { it.copy(isLoading = false, error = "File not found") }
                    return@launch
                }
                val (inputStream, cryptoMeta) = repository.downloadFile(bucket, fileId, apiKey)
                val encryptedBytes = inputStream.readBytes()
                inputStream.close()
                val (iv, salt) = CryptoUtils.unbuildCryptoHeader(cryptoMeta)
                val encKey = CryptoUtils.createEncryptionKeyFromPassword(bucketPassword, salt)
                val content = CryptoUtils.decrypt(encKey, iv, encryptedBytes)
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
