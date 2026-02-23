package xyz.nkrypt.android.ui.welcome

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.AppPreferences
import javax.inject.Inject

@HiltViewModel
class WelcomeViewModel @Inject constructor(
    private val appPreferences: AppPreferences
) : ViewModel() {

    private val _tosAndPrivacyAgreed = MutableStateFlow(false)
    val tosAndPrivacyAgreed: StateFlow<Boolean> = _tosAndPrivacyAgreed.asStateFlow()

    init {
        viewModelScope.launch {
            _tosAndPrivacyAgreed.value = appPreferences.tosAndPrivacyAgreed.first()
        }
    }

    fun setTosAndPrivacyAgreed(agreed: Boolean) {
        _tosAndPrivacyAgreed.value = agreed
    }

    suspend fun agreeAndStart() {
        appPreferences.setTosAndPrivacyAgreed(true)
    }
}
