//! Docker and Docker Compose availability checks

use std::process::Command;

/// Check if Docker and Docker Compose are installed and usable.
/// Returns (ready, status_message).
pub fn check_docker_ready() -> (bool, String) {
    // Check for docker (or podman)
    let docker_cmd = if command_exists("docker") {
        "docker"
    } else if command_exists("podman") {
        "podman"
    } else {
        return (
            false,
            "Docker or Podman not found. Install Docker to continue.".to_string(),
        );
    };

    // Check docker/podman works (e.g. daemon running, permissions)
    let check = Command::new(docker_cmd)
        .arg("info")
        .output();

    match check {
        Ok(output) if output.status.success() => {}
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            if stderr.contains("permission denied") || stderr.contains("permissions") {
                return (
                    false,
                    format!(
                        "{} found but permission denied. Add your user to the 'docker' group:\n  sudo usermod -aG docker $USER\nThen log out and back in.",
                        docker_cmd
                    ),
                );
            }
            if stderr.contains("Cannot connect") || stderr.contains("Is the docker daemon running") {
                return (
                    false,
                    format!("{} found but daemon not running. Start the {} service.", docker_cmd, docker_cmd),
                );
            }
            return (
                false,
                format!("{} check failed: {}", docker_cmd, stderr.trim()),
            );
        }
        Err(e) => {
            return (
                false,
                format!("Failed to run {}: {}", docker_cmd, e),
            );
        }
    }

    // Check docker compose (plugin or standalone)
    let compose_ok = Command::new(docker_cmd)
        .args(["compose", "version"])
        .output()
        .map(|o| o.status.success())
        .unwrap_or(false);

    if !compose_ok {
        return (
            false,
            format!(
                "{} found but Docker Compose not available.\nInstall: sudo apt-get install docker-compose-plugin",
                docker_cmd
            ),
        );
    }

    (
        true,
        format!("{} and Docker Compose are ready.", docker_cmd),
    )
}

pub fn command_exists(cmd: &str) -> bool {
    Command::new("which")
        .arg(cmd)
        .output()
        .map(|o| o.status.success())
        .unwrap_or(false)
}

