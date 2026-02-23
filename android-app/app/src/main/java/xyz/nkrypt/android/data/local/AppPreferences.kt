package xyz.nkrypt.android.data.local

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "app_preferences")

@Singleton
class AppPreferences @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val dataStore = context.dataStore

    private object Keys {
        val TO_S_AND_PRIVACY_AGREED = booleanPreferencesKey("tos_and_privacy_agreed")
    }

    val tosAndPrivacyAgreed: Flow<Boolean> = dataStore.data.map { prefs ->
        prefs[Keys.TO_S_AND_PRIVACY_AGREED] ?: false
    }

    suspend fun setTosAndPrivacyAgreed(agreed: Boolean) {
        dataStore.edit { prefs ->
            prefs[Keys.TO_S_AND_PRIVACY_AGREED] = agreed
        }
    }
}
