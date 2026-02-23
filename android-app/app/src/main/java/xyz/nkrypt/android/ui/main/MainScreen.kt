package xyz.nkrypt.android.ui.main

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Rule
import androidx.compose.material.icons.filled.Folder
import androidx.compose.material.icons.filled.Cloud
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.ui.localbuckets.BrowseBucketScreen
import xyz.nkrypt.android.ui.localbuckets.LocalBucketsScreen
import xyz.nkrypt.android.ui.remotebuckets.AddRemoteBucketScreen
import xyz.nkrypt.android.ui.remotebuckets.BrowseRemoteBucketScreen
import xyz.nkrypt.android.ui.remotebuckets.RemoteBucketsScreen
import xyz.nkrypt.android.ui.rules.RulesScreen
import xyz.nkrypt.android.ui.settings.SettingsScreen

sealed class MainTab(val route: String, val title: String) {
    data object LocalBuckets : MainTab("local_buckets", "Local Buckets")
    data object RemoteBuckets : MainTab("remote_buckets", "Remote Buckets")
    data object Rules : MainTab("rules", "Rules")
    data object Settings : MainTab("settings", "Settings")
}

const val BROWSE_BUCKET_ROUTE = "local_bucket/{bucketId}"
const val FILE_PREVIEW_ROUTE = "local_file/{bucketId}/{fileId}"
const val ADD_REMOTE_BUCKET_ROUTE = "add_remote_bucket"
const val BROWSE_REMOTE_BUCKET_ROUTE = "remote_bucket/{bucketId}"
const val REMOTE_FILE_PREVIEW_ROUTE = "remote_file/{bucketId}/{fileId}"
const val PROGRESS_ROUTE = "progress"

private val tabs = listOf(
    MainTab.LocalBuckets,
    MainTab.RemoteBuckets,
    MainTab.Rules,
    MainTab.Settings
)

@Composable
fun MainScreen(
    onLogout: () -> Unit
) {
    val navController = rememberNavController()
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentDestination = navBackStackEntry?.destination

    Scaffold(
        bottomBar = {
            NavigationBar {
                tabs.forEach { tab ->
                    val selected = currentDestination?.hierarchy?.any { it.route == tab.route } == true
                    NavigationBarItem(
                        icon = {
                            Icon(
                                imageVector = when (tab) {
                                    MainTab.LocalBuckets -> Icons.Default.Folder
                                    MainTab.RemoteBuckets -> Icons.Default.Cloud
                                    MainTab.Rules -> Icons.AutoMirrored.Filled.Rule
                                    MainTab.Settings -> Icons.Default.Settings
                                },
                                contentDescription = tab.title
                            )
                        },
                        label = { Text(tab.title) },
                        selected = selected,
                        onClick = {
                            navController.navigate(tab.route) {
                                popUpTo(navController.graph.findStartDestination().id) {
                                    saveState = true
                                }
                                launchSingleTop = true
                                restoreState = true
                            }
                        },
                        colors = NavigationBarItemDefaults.colors(
                            selectedIconColor = MaterialTheme.colorScheme.onPrimaryContainer,
                            selectedTextColor = MaterialTheme.colorScheme.onPrimaryContainer,
                            indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                            unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                            unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    )
                }
            }
        }
    ) { paddingValues ->
        NavHost(
            navController = navController,
            startDestination = MainTab.LocalBuckets.route,
            modifier = Modifier.padding(paddingValues)
        ) {
            composable(MainTab.LocalBuckets.route) {
                val viewModel: xyz.nkrypt.android.ui.localbuckets.LocalBucketsViewModel = hiltViewModel()
                LocalBucketsScreen(
                    viewModel = viewModel,
                    onBucketClick = { bucket ->
                        navController.navigate("local_bucket/${bucket.id}")
                    }
                )
            }
            composable(
                route = BROWSE_BUCKET_ROUTE,
                arguments = listOf(navArgument("bucketId") { type = NavType.StringType })
            ) { backStackEntry ->
                val bucketId = backStackEntry.arguments?.getString("bucketId") ?: return@composable
                BrowseBucketScreen(
                    bucketId = bucketId,
                    onBack = { navController.popBackStack() },
                    onFileClick = { fileId ->
                        navController.navigate("local_file/$bucketId/$fileId")
                    }
                )
            }
            composable(
                route = FILE_PREVIEW_ROUTE,
                arguments = listOf(
                    navArgument("bucketId") { type = NavType.StringType },
                    navArgument("fileId") { type = NavType.StringType }
                )
            ) { backStackEntry ->
                val bucketId = backStackEntry.arguments?.getString("bucketId") ?: return@composable
                val fileId = backStackEntry.arguments?.getString("fileId") ?: return@composable
                xyz.nkrypt.android.ui.localbuckets.FilePreviewScreen(
                    bucketId = bucketId,
                    fileId = fileId,
                    onBack = { navController.popBackStack() }
                )
            }
            composable(MainTab.RemoteBuckets.route) {
                val viewModel: xyz.nkrypt.android.ui.remotebuckets.RemoteBucketsViewModel = hiltViewModel()
                RemoteBucketsScreen(
                    viewModel = viewModel,
                    onBucketClick = { bucket ->
                        navController.navigate("remote_bucket/${bucket.id}")
                    },
                    onAddBucket = { navController.navigate(ADD_REMOTE_BUCKET_ROUTE) }
                )
            }
            composable(ADD_REMOTE_BUCKET_ROUTE) {
                AddRemoteBucketScreen(
                    onBack = { navController.popBackStack() },
                    onBucketsAdded = { navController.popBackStack() }
                )
            }
            composable(
                route = BROWSE_REMOTE_BUCKET_ROUTE,
                arguments = listOf(navArgument("bucketId") { type = NavType.StringType })
            ) { backStackEntry ->
                val bucketId = backStackEntry.arguments?.getString("bucketId") ?: return@composable
                BrowseRemoteBucketScreen(
                    bucketId = bucketId,
                    onBack = { navController.popBackStack() },
                    onFileClick = { fileId ->
                        navController.navigate("remote_file/$bucketId/$fileId")
                    }
                )
            }
            composable(
                route = REMOTE_FILE_PREVIEW_ROUTE,
                arguments = listOf(
                    navArgument("bucketId") { type = NavType.StringType },
                    navArgument("fileId") { type = NavType.StringType }
                )
            ) { backStackEntry ->
                val bucketId = backStackEntry.arguments?.getString("bucketId") ?: return@composable
                val fileId = backStackEntry.arguments?.getString("fileId") ?: return@composable
                xyz.nkrypt.android.ui.remotebuckets.RemoteFilePreviewScreen(
                    bucketId = bucketId,
                    fileId = fileId,
                    onBack = { navController.popBackStack() }
                )
            }
            composable(MainTab.Rules.route) {
                RulesScreen(
                    onNavigateToProgress = { navController.navigate(PROGRESS_ROUTE) },
                    onNavigateToLocalBuckets = {
                        navController.navigate(MainTab.LocalBuckets.route) {
                            popUpTo(navController.graph.findStartDestination().id) { saveState = true }
                            launchSingleTop = true
                            restoreState = true
                        }
                    },
                    onNavigateToRemoteBuckets = {
                        navController.navigate(MainTab.RemoteBuckets.route) {
                            popUpTo(navController.graph.findStartDestination().id) { saveState = true }
                            launchSingleTop = true
                            restoreState = true
                        }
                    }
                )
            }
            composable(PROGRESS_ROUTE) {
                xyz.nkrypt.android.ui.rules.ProgressScreen(
                    onBack = { navController.popBackStack() }
                )
            }
            composable(MainTab.Settings.route) {
                SettingsScreen(onLogout = onLogout)
            }
        }
    }
}
