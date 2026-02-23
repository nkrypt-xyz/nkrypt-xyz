package xyz.nkrypt.android.ui.navigation

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.platform.LocalContext
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import xyz.nkrypt.android.ui.masterpassword.MasterPasswordScreen
import xyz.nkrypt.android.ui.masterpassword.MasterPasswordViewModel
import xyz.nkrypt.android.ui.permissions.PermissionsScreen
import xyz.nkrypt.android.ui.permissions.PermissionsViewModel
import xyz.nkrypt.android.ui.main.MainScreen
import xyz.nkrypt.android.ui.welcome.WelcomeScreen
import xyz.nkrypt.android.ui.welcome.WelcomeViewModel

sealed class AppRoute(val route: String) {
    data object Welcome : AppRoute("welcome")
    data object MasterPassword : AppRoute("master_password")
    data object Permissions : AppRoute("permissions")
    data object Main : AppRoute("main")
}

@Composable
fun AppNavigation(
    navController: NavHostController = rememberNavController(),
    welcomeViewModel: WelcomeViewModel,
    masterPasswordViewModel: MasterPasswordViewModel,
    permissionsViewModel: PermissionsViewModel
) {
    val context = LocalContext.current
    val tosAgreed by welcomeViewModel.tosAndPrivacyAgreed.collectAsState(initial = false)
    val hasMasterPassword = masterPasswordViewModel.hasMasterPassword()
    val hasAllPermissions = permissionsViewModel.hasAllPermissions(context)

    val startDestination = when {
        !tosAgreed -> AppRoute.Welcome.route
        !hasMasterPassword -> AppRoute.MasterPassword.route
        !hasAllPermissions -> AppRoute.Permissions.route
        else -> AppRoute.Main.route
    }

    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable(AppRoute.Welcome.route) {
            WelcomeScreen(
                viewModel = welcomeViewModel,
                onAgree = { navController.navigate(AppRoute.MasterPassword.route) }
            )
        }

        composable(AppRoute.MasterPassword.route) {
            MasterPasswordScreen(
                isFirstTime = !hasMasterPassword,
                onSuccess = { navController.navigate(AppRoute.Permissions.route) },
                viewModel = masterPasswordViewModel
            )
        }

        composable(AppRoute.Permissions.route) {
            PermissionsScreen(
                viewModel = permissionsViewModel,
                onAllGranted = { navController.navigate(AppRoute.Main.route) }
            )
        }

        composable(AppRoute.Main.route) {
            MainScreen(
                onLogout = {
                    masterPasswordViewModel.clearMasterPassword()
                    navController.navigate(AppRoute.MasterPassword.route) {
                        popUpTo(0) { inclusive = true }
                    }
                }
            )
        }
    }
}
