package xyz.nkrypt.android.ui.welcome

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Button
import androidx.compose.foundation.clickable
import androidx.compose.material3.Checkbox
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import xyz.nkrypt.android.ui.legal.LegalContent
import xyz.nkrypt.android.ui.legal.LegalDialog

data class WelcomeSlide(
    val title: String,
    val description: String
)

private val slides = listOf(
    WelcomeSlide(
        title = "Local Buckets",
        description = "Create fully offline encrypted storage. Your files stay on your device, encrypted with your password."
    ),
    WelcomeSlide(
        title = "Remote Buckets",
        description = "Connect to nkrypt.xyz servers. Browse and manage your encrypted files from anywhere."
    ),
    WelcomeSlide(
        title = "Auto-Import",
        description = "Automatically encrypt content from any folder into a local bucket. Keep, trash, or delete originals."
    ),
    WelcomeSlide(
        title = "Auto-Sync",
        description = "One-way sync from local buckets to remote. Backup your encrypted files to the cloud."
    ),
    WelcomeSlide(
        title = "End-to-End Encryption",
        description = "AES-256-GCM encryption. Your keys, your data. We never see your files."
    )
)

@OptIn(ExperimentalFoundationApi::class)
@Composable
fun WelcomeScreen(
    viewModel: WelcomeViewModel,
    onAgree: () -> Unit
) {
    var agreed by remember { mutableStateOf(false) }
    var showPrivacyPolicyDialog by remember { mutableStateOf(false) }
    var showTermsDialog by remember { mutableStateOf(false) }
    val pagerState = rememberPagerState(pageCount = { slides.size })
    val scope = rememberCoroutineScope()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp)
    ) {
        HorizontalPager(
            state = pagerState,
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
        ) { page ->
            val slide = slides[page]
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(horizontal = 16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                Text(
                    text = slide.title,
                    style = MaterialTheme.typography.headlineMedium,
                    color = MaterialTheme.colorScheme.onBackground,
                    textAlign = TextAlign.Center
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = slide.description,
                    style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.8f),
                    textAlign = TextAlign.Center
                )
            }
        }

        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 16.dp),
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically
        ) {
            slides.forEachIndexed { index, _ ->
                val isSelected = pagerState.currentPage == index
                val alpha by animateFloatAsState(
                    targetValue = if (isSelected) 1f else 0.4f,
                    animationSpec = tween(200)
                )
                val size = if (isSelected) 10.dp else 8.dp
                Box(
                    modifier = Modifier
                        .padding(4.dp)
                        .size(size)
                        .clip(CircleShape)
                        .background(
                            MaterialTheme.colorScheme.primary.copy(alpha = alpha)
                        )
                )
            }
        }

        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 8.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Checkbox(
                checked = agreed,
                onCheckedChange = { agreed = it },
                modifier = Modifier.padding(end = 8.dp)
            )
            Column {
                Row {
                    Text(
                        text = "I agree to the ",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onBackground
                    )
                    Text(
                        text = "Terms of Service",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.primary,
                        modifier = Modifier.clickable { showTermsDialog = true }
                    )
                    Text(
                        text = " and ",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onBackground
                    )
                    Text(
                        text = "Privacy Policy",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.primary,
                        modifier = Modifier.clickable { showPrivacyPolicyDialog = true }
                    )
                }
            }
        }

        Button(
            onClick = {
                scope.launch {
                    viewModel.agreeAndStart()
                    onAgree()
                }
            },
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
            enabled = agreed
        ) {
            Text("Agree and Start")
        }
    }

    if (showPrivacyPolicyDialog) {
        LegalDialog(
            title = "Privacy Policy",
            content = LegalContent.PRIVACY_POLICY,
            onDismiss = { showPrivacyPolicyDialog = false }
        )
    }
    if (showTermsDialog) {
        LegalDialog(
            title = "Terms of Service",
            content = LegalContent.TERMS_OF_SERVICE,
            onDismiss = { showTermsDialog = false }
        )
    }
}
