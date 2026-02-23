package xyz.nkrypt.android.ui.masterpassword

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import xyz.nkrypt.android.data.local.MasterPasswordStore
import javax.inject.Inject

@HiltViewModel
class MasterPasswordViewModel @Inject constructor(
    private val masterPasswordStore: MasterPasswordStore
) : ViewModel() {

    fun hasMasterPassword(): Boolean = masterPasswordStore.hasMasterPassword()

    /**
     * Verify password (returning user) or set password (first time).
     * Returns true if success.
     */
    fun verifyOrSet(password: String, isFirstTime: Boolean): Boolean {
        return if (isFirstTime) {
            masterPasswordStore.setMasterPassword(password)
            true
        } else {
            val stored = masterPasswordStore.getMasterPassword()
            stored != null && stored == password
        }
    }

    fun clearMasterPassword() {
        masterPasswordStore.clearMasterPassword()
    }
}
