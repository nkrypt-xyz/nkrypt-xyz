package xyz.nkrypt.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import dagger.hilt.android.AndroidEntryPoint
import xyz.nkrypt.android.ui.navigation.AppNavigation
import xyz.nkrypt.android.ui.theme.NkryptTheme

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            NkryptTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    AppNavigation(
                        welcomeViewModel = hiltViewModel(),
                        masterPasswordViewModel = hiltViewModel(),
                        permissionsViewModel = hiltViewModel()
                    )
                }
            }
        }
    }
}
