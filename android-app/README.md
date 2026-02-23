# nkrypt.xyz Android App

Native Android app for nkrypt.xyz encrypted file storage. Kotlin + Jetpack Compose, Material 3.

## Quick Start

### Prerequisites

- **JDK 17**
- **Android SDK 35** (API 35)
- **Android Studio Ladybug (2024.2.1)** or newer (optional, for IDE)

Min SDK: 31 (Android 12)

### Build Debug APK

```bash
./build-debug.sh
```

Or with Gradle directly:

```bash
ANDROID_USER_HOME="$PWD/.android" ./gradlew assembleDebug
```

Output: `app/build/outputs/apk/debug/app-debug.apk`

### Run on Device/Emulator

```bash
./gradlew installDebug
```

## Commands

| Command | Description |
|---------|-------------|
| `./build-debug.sh` | Build debug APK (uses project `.android` dir) |
| `./gradlew assembleDebug` | Build debug APK |
| `./gradlew installDebug` | Install debug APK on connected device |
| `./gradlew clean` | Clean build artifacts |

## Documentation

See [dev-docs/](../dev-docs/) for architecture and contribution guidelines.
