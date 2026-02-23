package xyz.nkrypt.android.data.local

import android.content.Context
import dagger.hilt.android.qualifiers.ApplicationContext
import java.io.File
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Stores master password in a local file for testing.
 * TODO: Replace with EncryptedSharedPreferences or Android Keystore.
 */
@Singleton
class MasterPasswordStore @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val file: File
        get() = File(context.filesDir, "master_password_test.txt")

    fun hasMasterPassword(): Boolean = file.exists() && file.readText().isNotBlank()

    fun getMasterPassword(): String? {
        return if (file.exists()) {
            file.readText().takeIf { it.isNotBlank() }
        } else null
    }

    fun setMasterPassword(password: String) {
        file.writeText(password)
    }

    fun clearMasterPassword() {
        if (file.exists()) {
            file.delete()
        }
    }
}
