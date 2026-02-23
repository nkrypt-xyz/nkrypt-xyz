#!/bin/bash
# Wrapper to run unit tests with ANDROID_USER_HOME in project dir
# (avoids "Cannot create directory /root/.android" in CI/sandbox environments)
cd "$(dirname "$0")"
ANDROID_USER_HOME="$PWD/.android" ./gradlew testDebugUnitTest "$@"
