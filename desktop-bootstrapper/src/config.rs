//! Configuration management for the desktop instance.
//! Reads/writes .env for docker-compose.

#![allow(dead_code)]

use std::path::PathBuf;

/// Path to the repo root (parent of desktop-bootstrapper)
pub fn repo_root() -> PathBuf {
    // CARGO_MANIFEST_DIR is desktop-bootstrapper/ at compile time
    let manifest_dir = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    manifest_dir
        .parent()
        .expect("desktop-bootstrapper should have parent")
        .to_path_buf()
}

/// Path to the docker-compose file
pub fn compose_path() -> PathBuf {
    repo_root().join("desktop-bootstrapper").join("docker-compose.yml")
}

/// Path to the .env file for docker-compose
pub fn env_path() -> PathBuf {
    repo_root().join("desktop-bootstrapper").join(".env")
}
