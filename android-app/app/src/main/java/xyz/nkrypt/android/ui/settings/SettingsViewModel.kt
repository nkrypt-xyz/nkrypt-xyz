package xyz.nkrypt.android.ui.settings

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.LocalBucketRepository
import xyz.nkrypt.android.data.local.MasterPasswordStore
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import javax.inject.Inject

@HiltViewModel
class SettingsViewModel @Inject constructor(
    private val masterPasswordStore: MasterPasswordStore,
    private val localBucketRepository: LocalBucketRepository,
    private val remoteBucketRepository: RemoteBucketRepository
) : ViewModel() {

    private val _changePasswordState = MutableStateFlow<ChangePasswordState>(ChangePasswordState.Idle)
    val changePasswordState: StateFlow<ChangePasswordState> = _changePasswordState.asStateFlow()

    fun changeMasterPassword(currentPassword: String, newPassword: String) {
        viewModelScope.launch {
            _changePasswordState.value = ChangePasswordState.Loading
            try {
                val stored = masterPasswordStore.getMasterPassword()
                if (stored == null || stored != currentPassword) {
                    _changePasswordState.value = ChangePasswordState.Error("Current password is incorrect")
                    return@launch
                }
                if (newPassword.isBlank()) {
                    _changePasswordState.value = ChangePasswordState.Error("New password cannot be empty")
                    return@launch
                }
                localBucketRepository.reEncryptAllWithNewMaster(currentPassword, newPassword)
                remoteBucketRepository.reEncryptAllWithNewMaster(currentPassword, newPassword)
                masterPasswordStore.setMasterPassword(newPassword)
                _changePasswordState.value = ChangePasswordState.Success
            } catch (e: Exception) {
                _changePasswordState.value = ChangePasswordState.Error(e.message ?: "Failed to change password")
            }
        }
    }

    fun clearChangePasswordState() {
        _changePasswordState.value = ChangePasswordState.Idle
    }
}

sealed class ChangePasswordState {
    data object Idle : ChangePasswordState()
    data object Loading : ChangePasswordState()
    data object Success : ChangePasswordState()
    data class Error(val message: String) : ChangePasswordState()
}
