//! nkrypt.xyz Desktop Bootstrapper
//! Manages a standalone Docker-based instance of nkrypt.xyz

mod docker;
mod config;

use slint::{ComponentHandle, ModelRc, VecModel};
use std::io::{BufRead, BufReader};
use std::path::PathBuf;
use std::process::Stdio;

slint::include_modules!();

fn main() {
    let ui = MainWindow::new().expect("Failed to create UI");

    // Check Docker availability on startup
    let ui_handle = ui.as_weak();
    std::thread::spawn(move || {
        let (ready, status) = docker::check_docker_ready();
        let _ = slint::invoke_from_event_loop(move || {
            if let Some(ui) = ui_handle.upgrade() {
                ui.set_docker_ready(ready);
                ui.set_docker_status(status.into());
                if ready {
                    load_config(&ui);
                    refresh_status(&ui);
                    // Sync auto-start checkbox with real systemd state
                    ui.set_autostart_enabled(autostart_is_enabled());
                    // Ensure data directory is world-writable on startup
                    let data_dir = ui.get_data_dir().to_string();
                    if !data_dir.is_empty() {
                        apply_data_dir_permissions(&ui, &data_dir);
                    }
                }
            }
        });
    });

    // Wire up callbacks
    let ui_handle = ui.as_weak();
    ui.on_browse_data_dir(move || {
        if let Some(ui) = ui_handle.upgrade() {
            browse_data_dir(&ui);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_start_services(move || {
        if let Some(ui) = ui_handle.upgrade() {
            run_compose_async(&ui, &["compose", "up", "-d"]);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_stop_services(move || {
        if let Some(ui) = ui_handle.upgrade() {
            run_compose_async(&ui, &["compose", "down"]);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_refresh_status(move || {
        if let Some(ui) = ui_handle.upgrade() {
            refresh_status(&ui);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_toggle_autostart(move |enabled| {
        if let Some(ui) = ui_handle.upgrade() {
            toggle_autostart(&ui, enabled);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_remove_all(move || {
        if let Some(ui) = ui_handle.upgrade() {
            remove_all(&ui);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_save_logs(move || {
        if let Some(ui) = ui_handle.upgrade() {
            save_logs(&ui);
        }
    });

    let ui_handle = ui.as_weak();
    ui.on_copy_logs(move || {
        if let Some(ui) = ui_handle.upgrade() {
            copy_logs(&ui);
        }
    });

    ui.on_open_url(|url| {
        let url = url.to_string();
        std::thread::spawn(move || {
            let _ = std::process::Command::new("xdg-open").arg(&url).spawn();
        });
    });

    // Periodic status refresh every 10 seconds
    let ui_handle = ui.as_weak();
    std::thread::spawn(move || {
        loop {
            std::thread::sleep(std::time::Duration::from_secs(10));
            let ui_handle = ui_handle.clone();
            let _ = slint::invoke_from_event_loop(move || {
                if let Some(ui) = ui_handle.upgrade() {
                    if ui.get_docker_ready() && !ui.get_operation_in_progress() {
                        refresh_status(&ui);
                    }
                }
            });
        }
    });

    ui.run().expect("UI event loop failed");
}

/// Make the data directory and its subdirectories world-writable so that
/// rootless Podman containers (running as different UIDs inside the user namespace)
/// can write to the bind-mounted volumes.
///
/// Subdirectories may be owned by container uids (e.g. uid 524357 for postgres uid 70).
/// Those can't be chmod'd directly from the host user, so we fall back to
/// `podman unshare chmod` which runs inside the user namespace where we're root.
///
/// Returns Ok(()) on success, Err(message) on failure.
fn ensure_data_dir_writable(data_dir: &str) -> Result<(), String> {
    use std::os::unix::fs::PermissionsExt;
    let base = PathBuf::from(data_dir);
    for sub in ["", "postgres", "redis", "minio"] {
        let p = if sub.is_empty() { base.clone() } else { base.join(sub) };
        if let Err(e) = std::fs::create_dir_all(&p) {
            return Err(format!("Failed to create {}: {}", p.display(), e));
        }
        // Try direct chmod first; fall back to `podman unshare chmod` for dirs
        // owned by container uids that we don't own on the host.
        let chmod_result = std::fs::set_permissions(&p, std::fs::Permissions::from_mode(0o777));
        if chmod_result.is_err() {
            let unshare = std::process::Command::new("podman")
                .args(["unshare", "chmod", "777", &p.to_string_lossy()])
                .output();
            match unshare {
                Ok(o) if o.status.success() => {}
                Ok(o) => {
                    return Err(format!(
                        "chmod 777 {} failed (direct and via podman unshare): {}",
                        p.display(),
                        String::from_utf8_lossy(&o.stderr).trim()
                    ));
                }
                Err(e) => {
                    return Err(format!(
                        "chmod 777 {} failed and podman unshare not available: {}",
                        p.display(), e
                    ));
                }
            }
        }
    }
    Ok(())
}

fn browse_data_dir(ui: &MainWindow) {
    let start = if ui.get_data_dir().is_empty() {
        std::env::home_dir().unwrap_or_else(|| PathBuf::from("."))
    } else {
        PathBuf::from(&ui.get_data_dir())
    };

    if let Some(path) = rfd::FileDialog::new()
        .set_title("Select data directory")
        .set_directory(&start)
        .pick_folder()
    {
        let path_str = path.to_string_lossy().to_string();
        ui.set_data_dir(path_str.clone().into());
        save_config(ui);
        apply_data_dir_permissions(ui, &path_str);
    }
}

fn apply_data_dir_permissions(ui: &MainWindow, data_dir: &str) {
    match ensure_data_dir_writable(data_dir) {
        Ok(()) => {
            append_log(ui, &format!("Data directory permissions set (world-writable): {}", data_dir));
        }
        Err(e) => {
            let msg = format!(
                "WARNING: Could not make data directory world-writable: {}\n\
                Containers may fail to write to their volumes. \
                Run manually: chmod -R 777 '{}'",
                e, data_dir
            );
            append_log(ui, &msg);
            rfd::MessageDialog::new()
                .set_title("Data Directory Permission Warning")
                .set_description(&msg)
                .set_level(rfd::MessageLevel::Warning)
                .show();
        }
    }
}

fn load_config(ui: &MainWindow) {
    let env_path = config::env_path();
    if env_path.exists() {
        if let Ok(content) = std::fs::read_to_string(&env_path) {
            for line in content.lines() {
                let line = line.trim();
                if line.is_empty() || line.starts_with('#') {
                    continue;
                }
                if let Some((key, value)) = line.split_once('=') {
                    let key = key.trim();
                    let value = value.trim().trim_matches('"');
                    match key {
                        "DATA_DIR" => ui.set_data_dir(value.into()),
                        "POSTGRES_PORT" => ui.set_postgres_port(value.parse().unwrap_or(9200)),
                        "REDIS_PORT" => ui.set_redis_port(value.parse().unwrap_or(9201)),
                        "MINIO_PORT" => ui.set_minio_port(value.parse().unwrap_or(9202)),
                        "MINIO_CONSOLE_PORT" => ui.set_minio_console_port(value.parse().unwrap_or(9203)),
                        "WEB_SERVER_PORT" => ui.set_web_server_port(value.parse().unwrap_or(9204)),
                        "WEB_CLIENT_PORT" => ui.set_web_client_port(value.parse().unwrap_or(9205)),
                        "NK_IAM_DEFAULT_ADMIN_PASSWORD" => ui.set_admin_password(value.into()),
                        "NK_LOG_LEVEL" => ui.set_log_level(value.into()),
                        _ => {}
                    }
                }
            }
        }
    }
    if ui.get_data_dir().is_empty() {
        let default = config::repo_root()
            .join("nkrypt-xyz-data")
            .to_string_lossy()
            .into_owned();
        ui.set_data_dir(default.into());
    }
}

fn save_config(ui: &MainWindow) {
    let env_path = config::env_path();
    let data_dir = ui.get_data_dir();
    if data_dir.is_empty() {
        return;
    }

    let server_url = format!("http://localhost:{}", ui.get_web_server_port());

    let db_url = format!(
        "postgres://nkrypt:nkrypt_password@127.0.0.1:{}/nkrypt?sslmode=disable",
        ui.get_postgres_port()
    );

    let content = format!(
        r#"# nkrypt.xyz Desktop - managed by bootstrapper
DATA_DIR={}
POSTGRES_PORT={}
REDIS_PORT={}
MINIO_PORT={}
MINIO_CONSOLE_PORT={}
WEB_SERVER_PORT={}
WEB_CLIENT_PORT={}
POSTGRES_PASSWORD=nkrypt_password
NK_DATABASE_URL={}
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
NK_IAM_DEFAULT_ADMIN_PASSWORD={}
VITE_DEFAULT_SERVER_URL={}
VITE_CLIENT_APPLICATION_NAME=nkrypt-web-client
NK_LOG_LEVEL={}
NK_LOG_FORMAT=json
"#,
        data_dir,
        ui.get_postgres_port(),
        ui.get_redis_port(),
        ui.get_minio_port(),
        ui.get_minio_console_port(),
        ui.get_web_server_port(),
        ui.get_web_client_port(),
        db_url,
        ui.get_admin_password(),
        server_url,
        ui.get_log_level(),
    );

    if let Err(e) = std::fs::write(&env_path, content) {
        append_log(ui, &format!("Failed to save config: {}", e));
    }
}

/// Apply all pending migrations by running psql inside the postgres container.
/// Uses the unix socket (no TCP/network involved) — works on any Podman setup.
/// Maintains a `schema_migrations` table compatible with golang-migrate.
fn run_migrations_via_psql(
    docker_cmd: &str,
    migrations_dir: &std::path::Path,
    ui_handle: &slint::Weak<MainWindow>,
) -> Result<(), String> {
    // Helper: run a SQL string via `exec psql -c`
    let psql_exec = |sql: &str| -> Result<String, String> {
        let out = std::process::Command::new(docker_cmd)
            .args(["exec", "nkrypt-desktop-postgres",
                   "psql", "-U", "nkrypt", "-d", "nkrypt", "-tAq", "-c", sql])
            .output()
            .map_err(|e| format!("exec psql: {}", e))?;
        if out.status.success() {
            Ok(String::from_utf8_lossy(&out.stdout).to_string())
        } else {
            Err(String::from_utf8_lossy(&out.stderr).trim().to_string())
        }
    };

    // Helper: pipe a SQL file to psql via stdin
    let psql_file = |path: &std::path::Path| -> Result<(), String> {
        let sql = std::fs::read_to_string(path)
            .map_err(|e| format!("read {}: {}", path.display(), e))?;
        let mut child = std::process::Command::new(docker_cmd)
            .args(["exec", "-i", "nkrypt-desktop-postgres",
                   "psql", "-U", "nkrypt", "-d", "nkrypt", "-v", "ON_ERROR_STOP=1"])
            .stdin(std::process::Stdio::piped())
            .stdout(std::process::Stdio::piped())
            .stderr(std::process::Stdio::piped())
            .spawn()
            .map_err(|e| format!("spawn psql: {}", e))?;
        use std::io::Write;
        if let Some(mut stdin) = child.stdin.take() {
            stdin.write_all(sql.as_bytes()).ok();
        }
        let out = child.wait_with_output().map_err(|e| e.to_string())?;
        if out.status.success() {
            Ok(())
        } else {
            Err(String::from_utf8_lossy(&out.stderr).trim().to_string())
        }
    };

    // Ensure schema_migrations table exists (golang-migrate compatible)
    psql_exec(
        "CREATE TABLE IF NOT EXISTS schema_migrations \
         (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL);"
    )?;

        // Get already-applied (clean) versions
    let applied_raw = psql_exec(
        "SELECT version FROM schema_migrations WHERE dirty = false ORDER BY version;"
    )?;
    let applied: std::collections::HashSet<i64> = applied_raw
        .lines()
        .filter_map(|l| l.trim().parse::<i64>().ok())
        .collect();

    // Log any dirty (previously failed) migrations as warnings
    let dirty_raw = psql_exec(
        "SELECT version FROM schema_migrations WHERE dirty = true ORDER BY version;"
    )?;
    for v in dirty_raw.lines().filter(|l| !l.trim().is_empty()) {
        let ui_h = ui_handle.clone();
        let msg = format!("  WARNING: migration {} was previously left dirty — will retry", v.trim());
        let _ = slint::invoke_from_event_loop(move || {
            if let Some(ui) = ui_h.upgrade() { append_log(&ui, &msg); }
        });
    }

    // Collect and sort all .up.sql migration files
    let mut files: Vec<_> = std::fs::read_dir(migrations_dir)
        .map_err(|e| format!("read migrations dir: {}", e))?
        .filter_map(|e| e.ok())
        .filter(|e| {
            let name = e.file_name();
            let s = name.to_string_lossy();
            s.ends_with(".up.sql")
        })
        .collect();
    files.sort_by_key(|e| e.file_name());

    for entry in &files {
        let path = entry.path();
        let filename = entry.file_name();
        let name = filename.to_string_lossy();
        // Extract version number from filename (e.g. "000001_create_users_table.up.sql")
        let version: i64 = name.split('_').next()
            .and_then(|v| v.parse().ok())
            .ok_or_else(|| format!("Invalid migration filename: {}", name))?;

        if applied.contains(&version) {
            let ui_h = ui_handle.clone();
            let vstr = format!("  skip {} (already applied)", name);
            let _ = slint::invoke_from_event_loop(move || {
                if let Some(ui) = ui_h.upgrade() { append_log(&ui, &vstr); }
            });
            continue;
        }

        let ui_h = ui_handle.clone();
        let vstr = format!("  apply {}", name);
        let _ = slint::invoke_from_event_loop(move || {
            if let Some(ui) = ui_h.upgrade() { append_log(&ui, &vstr); }
        });

        // Mark dirty
        psql_exec(&format!(
            "INSERT INTO schema_migrations (version, dirty) VALUES ({}, true) \
             ON CONFLICT (version) DO UPDATE SET dirty = true;",
            version
        ))?;

        // Apply migration. If objects already exist (e.g. migration was applied
        // manually without being recorded), treat it as already applied.
        match psql_file(&path) {
            Ok(()) => {}
            Err(ref e) if e.contains("already exists") => {
                let ui_h = ui_handle.clone();
                let msg = format!(
                    "  WARNING: {} — objects already exist, marking as applied",
                    name
                );
                let _ = slint::invoke_from_event_loop(move || {
                    if let Some(ui) = ui_h.upgrade() { append_log(&ui, &msg); }
                });
            }
            Err(e) => return Err(format!("failed to apply {}: {}", name, e)),
        }

        // Mark clean
        psql_exec(&format!(
            "UPDATE schema_migrations SET dirty = false WHERE version = {};",
            version
        ))?;
    }

    Ok(())
}

fn run_compose_async(ui: &MainWindow, args: &[&str]) {
    save_config(ui);

    let data_dir = ui.get_data_dir().to_string();
    if data_dir.is_empty() {
        append_log(ui, "Error: Data directory is required");
        return;
    }

    // Create data subdirs and ensure world-writable for rootless Podman containers
    if let Err(e) = ensure_data_dir_writable(&data_dir) {
        let msg = format!(
            "WARNING: Could not make data directory world-writable: {}\n\
            Containers may fail to write to their volumes. \
            Run manually: chmod -R 777 '{}'",
            e, data_dir
        );
        append_log(ui, &msg);
        rfd::MessageDialog::new()
            .set_title("Data Directory Permission Warning")
            .set_description(&msg)
            .set_level(rfd::MessageLevel::Warning)
            .show();
    }

    let compose_dir = config::repo_root().join("desktop-bootstrapper");
    let compose_file = compose_dir.join("docker-compose.yml");
    let compose_file_str = compose_file.to_string_lossy().into_owned();

    let docker_cmd = if docker::command_exists("docker") {
        "docker".to_string()
    } else {
        "podman".to_string()
    };

    // Build args: compose -f <file> <subcommand> <subargs...>
    let mut all_args: Vec<String> = vec!["compose".into(), "-f".into(), compose_file_str.clone()];
    if let Some((_, rest)) = args.split_first() {
        all_args.extend(rest.iter().map(|s| (*s).to_string()));
    }

    let postgres_port = ui.get_postgres_port();
    let redis_port = ui.get_redis_port();
    let minio_port = ui.get_minio_port();
    let minio_console_port = ui.get_minio_console_port();
    let web_server_port = ui.get_web_server_port();
    let web_client_port = ui.get_web_client_port();
    let admin_password = ui.get_admin_password().to_string();
    let log_level = ui.get_log_level().to_string();

    let run_msg = format!("Running: {} {}", docker_cmd, all_args.join(" "));
    append_log(ui, &run_msg);

    let is_up_d = args.len() >= 3 && args[1] == "up" && args[2] == "-d";

    ui.set_operation_in_progress(true);
    let ui_handle = ui.as_weak();
    std::thread::spawn(move || {
        // For "up -d": stop deps, up deps, migrate via psql inside container, up app.
        // Migrations are applied with `podman exec psql` inside the postgres container.
        // This uses the unix socket — no network involved, always works.
        let phases: Vec<(String, Vec<String>, Option<std::path::PathBuf>, Vec<(&str, String)>)> = if is_up_d {
            vec![
                (docker_cmd.clone(), vec!["compose".into(), "-f".into(), compose_file_str.clone(), "stop".into(), "postgres".into(), "redis".into(), "minio".into()], None, vec![]),
                (docker_cmd.clone(), vec!["compose".into(), "-f".into(), compose_file_str.clone(), "up".into(), "-d".into(), "postgres".into(), "redis".into(), "minio".into()], Some(compose_dir.clone()), vec![]),
                // Phase index 2 is handled specially below (psql migration)
                ("__migrate__".into(), vec![], None, vec![]),
                (docker_cmd.clone(), vec!["compose".into(), "-f".into(), compose_file_str.clone(), "up".into(), "-d".into(), "web-server".into(), "web-client".into()], Some(compose_dir.clone()), vec![]),
            ]
        } else {
            vec![(docker_cmd.clone(), all_args.clone(), Some(compose_dir.clone()), vec![])]
        };

        let env_vars = [
            ("DATA_DIR", data_dir.clone()),
            ("POSTGRES_PASSWORD", "nkrypt_password".to_string()),
            ("POSTGRES_PORT", postgres_port.to_string()),
            ("REDIS_PORT", redis_port.to_string()),
            ("MINIO_PORT", minio_port.to_string()),
            ("MINIO_CONSOLE_PORT", minio_console_port.to_string()),
            ("WEB_SERVER_PORT", web_server_port.to_string()),
            ("WEB_CLIENT_PORT", web_client_port.to_string()),
            ("NK_IAM_DEFAULT_ADMIN_PASSWORD", admin_password.clone()),
            ("VITE_DEFAULT_SERVER_URL", format!("http://localhost:{}", web_server_port)),
            ("NK_LOG_LEVEL", log_level.clone()),
        ];

        for (i, (phase_cmd, phase_args, phase_cwd, phase_extra_env)) in phases.iter().enumerate() {
            if i > 0 && phases.len() == 4 && i == 1 {
                std::thread::sleep(std::time::Duration::from_secs(2));
            }

            // Phase 3 (index 2) for "up -d": run migrations via psql inside the container.
            // This is the only approach that works reliably on rootless Podman (no network issues).
            if phases.len() == 4 && i == 2 {
                // Wait for postgres pg_isready (uses unix socket inside container — reliable).
                let ui_handle_wait = ui_handle.clone();
                let _ = slint::invoke_from_event_loop(move || {
                    if let Some(ui) = ui_handle_wait.upgrade() {
                        append_log(&ui, "Waiting for PostgreSQL to be ready...");
                    }
                });
                for _ in 1..=60 {
                    let pg_ready = std::process::Command::new(&docker_cmd)
                        .args(["exec", "nkrypt-desktop-postgres", "pg_isready", "-U", "nkrypt", "-d", "nkrypt"])
                        .output();
                    if pg_ready.as_ref().map(|o| o.status.success()).unwrap_or(false) {
                        break;
                    }
                    std::thread::sleep(std::time::Duration::from_secs(2));
                }
                let ui_handle_mig = ui_handle.clone();
                let _ = slint::invoke_from_event_loop(move || {
                    if let Some(ui) = ui_handle_mig.upgrade() {
                        append_log(&ui, "Phase 3: migrate (via psql inside container)");
                    }
                });
                let migrations_dir = config::repo_root().join("web-server").join("migrations");
                let mig_result = run_migrations_via_psql(&docker_cmd, &migrations_dir, &ui_handle);
                let ui_handle_done = ui_handle.clone();
                match mig_result {
                    Ok(()) => {
                        let _ = slint::invoke_from_event_loop(move || {
                            if let Some(ui) = ui_handle_done.upgrade() {
                                append_log(&ui, "Phase 3 completed successfully");
                            }
                        });
                    }
                    Err(e) => {
                        let _ = slint::invoke_from_event_loop(move || {
                            if let Some(ui) = ui_handle_done.upgrade() {
                                append_log(&ui, &format!("Phase 3 failed: {}", e));
                                refresh_status(&ui);
                                ui.set_operation_in_progress(false);
                            }
                        });
                        return;
                    }
                }
                continue; // Skip the normal phase runner for this index
            }

            let phase_msg = if phase_args.len() > 3 {
                format!("Phase {}: {} {}", i + 1, phase_cmd, phase_args[3..].join(" "))
            } else {
                format!("Phase {}: {} {}", i + 1, phase_cmd, phase_args.join(" "))
            };
            let ui_handle_phase = ui_handle.clone();
            let _ = slint::invoke_from_event_loop(move || {
                if let Some(ui) = ui_handle_phase.upgrade() {
                    append_log(&ui, &phase_msg);
                }
            });

            {
                let phase_succeeded;

            let mut cmd = std::process::Command::new(phase_cmd);
            cmd.args(phase_args);
            if let Some(ref cwd) = phase_cwd {
                cmd.current_dir(cwd);
            } else {
                cmd.current_dir(&compose_dir);
            }
            if phase_cmd.as_str() == docker_cmd.as_str() {
                for (k, v) in &env_vars {
                    cmd.env(k, v);
                }
            }
            for (k, v) in phase_extra_env {
                cmd.env(k, v);
            }
            cmd.stdout(Stdio::piped()).stderr(Stdio::piped());

            let mut child = match cmd.spawn() {
                Err(e) => {
                    let _ = slint::invoke_from_event_loop(move || {
                        if let Some(ui) = ui_handle.upgrade() {
                            append_log(&ui, &format!("Failed to run phase: {}", e));
                            ui.set_operation_in_progress(false);
                        }
                    });
                    return;
                }
                Ok(c) => c,
            };

            let stdout = child.stdout.take().unwrap();
            let stderr = child.stderr.take().unwrap();

            let ui_handle_stdout = ui_handle.clone();
            let stdout_handle = std::thread::spawn(move || {
                let reader = BufReader::new(stdout);
                for line in reader.lines() {
                    if let Ok(line) = line {
                        let ui_handle = ui_handle_stdout.clone();
                        let _ = slint::invoke_from_event_loop(move || {
                            if let Some(ui) = ui_handle.upgrade() {
                                append_log(&ui, &line);
                            }
                        });
                    }
                }
            });

            let ui_handle_stderr = ui_handle.clone();
            let stderr_handle = std::thread::spawn(move || {
                let reader = BufReader::new(stderr);
                for line in reader.lines() {
                    if let Ok(line) = line {
                        let ui_handle = ui_handle_stderr.clone();
                        let _ = slint::invoke_from_event_loop(move || {
                            if let Some(ui) = ui_handle.upgrade() {
                                append_log(&ui, &line);
                            }
                        });
                    }
                }
            });

            let _ = stdout_handle.join();
            let _ = stderr_handle.join();

            let status_result = child.wait();
            phase_succeeded = matches!(&status_result, Ok(s) if s.success()); #[allow(unused_assignments)]
            let final_msg = match &status_result {
                Ok(s) if s.success() => "Phase completed successfully".to_string(),
                Ok(s) => format!("Phase failed with code: {:?}", s.code()),
                Err(e) => format!("Phase failed: {}", e),
            };
            let ui_handle_phase_final = ui_handle.clone();
            let _ = slint::invoke_from_event_loop(move || {
                if let Some(ui) = ui_handle_phase_final.upgrade() {
                    append_log(&ui, &final_msg);
                }
            });

            if !phase_succeeded {
                let ui_handle_fail = ui_handle.clone();
                let _ = slint::invoke_from_event_loop(move || {
                    if let Some(ui) = ui_handle_fail.upgrade() {
                        refresh_status(&ui);
                        ui.set_operation_in_progress(false);
                    }
                });
                return;
            }
            } // end single-attempt block
        }

        let ui_handle_final = ui_handle.clone();
        let _ = slint::invoke_from_event_loop(move || {
            if let Some(ui) = ui_handle_final.upgrade() {
                append_log(&ui, "All phases completed successfully");
                refresh_status(&ui);
                ui.set_operation_in_progress(false);
            }
        });
    });
}

fn refresh_status(ui: &MainWindow) {
    // Use `docker/podman ps` directly — podman-compose ps doesn't support
    // -a or --format with Go templates, so we query the daemon directly
    // and filter by our container name prefix.
    let docker_cmd = if docker::command_exists("docker") {
        "docker"
    } else {
        "podman"
    };

    let format_output = std::process::Command::new(docker_cmd)
        .args([
            "ps", "-a",
            "--filter", "name=nkrypt-desktop",
            "--format", "{{.Names}}\t{{.Status}}",
        ])
        .output();

    match &format_output {
        Ok(o) if o.status.success() => {
            let format_out = String::from_utf8_lossy(&o.stdout).to_string();
            let statuses = parse_service_statuses(&format_out, ui);
            let model = VecModel::from(
                statuses
                    .into_iter()
                    .map(|s| ServiceStatus {
                        name: s.name.into(),
                        status: s.status.into(),
                        is_healthy: s.is_healthy,
                        port: s.port.into(),
                        url: s.url.into(),
                    })
                    .collect::<Vec<_>>(),
            );
            ui.set_service_statuses(ModelRc::new(model));
        }
        Ok(o) => {
            // Command succeeded but non-zero exit — show stopped state for all services
            let _ = String::from_utf8_lossy(&o.stderr);
            let statuses = parse_service_statuses("", ui);
            let model = VecModel::from(
                statuses
                    .into_iter()
                    .map(|s| ServiceStatus {
                        name: s.name.into(),
                        status: s.status.into(),
                        is_healthy: s.is_healthy,
                        port: s.port.into(),
                        url: s.url.into(),
                    })
                    .collect::<Vec<_>>(),
            );
            ui.set_service_statuses(ModelRc::new(model));
        }
        Err(_) => {
            ui.set_service_statuses(ModelRc::new(VecModel::default()));
        }
    }
}

struct ParsedServiceStatus {
    name: String,
    status: String,
    is_healthy: bool,
    port: String,
    url: String,
}

fn service_meta(ui: &MainWindow) -> Vec<(&'static str, &'static str, String, String)> {
    let wp = ui.get_web_server_port();
    let wcp = ui.get_web_client_port();
    let mp = ui.get_minio_port();
    let mcp = ui.get_minio_console_port();
    vec![
        ("nkrypt-desktop-postgres", "PostgreSQL",
            format!("{}", ui.get_postgres_port()), String::new()),
        ("nkrypt-desktop-redis",    "Redis",
            format!("{}", ui.get_redis_port()), String::new()),
        ("nkrypt-desktop-minio",    "MinIO",
            format!("{}, {} (console)", mp, mcp),
            format!("http://localhost:{}", mcp)),
        ("nkrypt-desktop-web-server", "Web Server",
            format!("{}", wp), String::new()),
        ("nkrypt-desktop-web-client", "Web Client",
            format!("{}", wcp),
            format!("http://localhost:{}", wcp)),
    ]
}

fn parse_service_statuses(format_output: &str, ui: &MainWindow) -> Vec<ParsedServiceStatus> {
    let meta = service_meta(ui);

    let mut result = Vec::new();
    for line in format_output.lines() {
        let (name, status) = match line.split_once('\t') {
            Some((n, s)) => (n.trim(), s.trim()),
            None => continue,
        };
        for (container_name, display_name, port, url) in &meta {
            if name == *container_name {
                let is_healthy = status.contains("healthy")
                    || (status.starts_with("Up") && !status.contains("unhealthy"));
                result.push(ParsedServiceStatus {
                    name: display_name.to_string(),
                    status: status.to_string(),
                    is_healthy,
                    port: port.clone(),
                    url: url.clone(),
                });
                break;
            }
        }
    }

    // Ensure all expected services appear (with stopped status if missing)
    for (_, display_name, port, url) in &meta {
        if !result.iter().any(|s| s.name == *display_name) {
            result.push(ParsedServiceStatus {
                name: display_name.to_string(),
                status: "not found".to_string(),
                is_healthy: false,
                port: port.clone(),
                url: url.clone(),
            });
        }
    }
    result.sort_by_key(|a| {
        meta.iter().position(|(_, n, _, _)| *n == a.name).unwrap_or(99)
    });
    result
}

fn remove_all(ui: &MainWindow) {
    save_config(ui);

    let data_dir = ui.get_data_dir();
    if data_dir.is_empty() {
        append_log(ui, "Error: Data directory is required");
        return;
    }

    let compose_dir = config::repo_root().join("desktop-bootstrapper");
    let compose_file = compose_dir.join("docker-compose.yml");

    let docker_cmd = if docker::command_exists("docker") {
        "docker"
    } else {
        "podman"
    };

    append_log(ui, "Removing all containers, volumes, networks, and images...");

    let mut cmd = std::process::Command::new(docker_cmd);
    cmd.args([
        "compose",
        "-f",
        &compose_file.to_string_lossy(),
        "down",
        "-v",
        "--rmi",
        "all",
    ])
    .current_dir(&compose_dir)
    .env("DATA_DIR", &data_dir)
    .env("POSTGRES_PORT", ui.get_postgres_port().to_string())
    .env("REDIS_PORT", ui.get_redis_port().to_string())
    .env("MINIO_PORT", ui.get_minio_port().to_string())
    .env("MINIO_CONSOLE_PORT", ui.get_minio_console_port().to_string())
    .env("WEB_SERVER_PORT", ui.get_web_server_port().to_string())
    .env("WEB_CLIENT_PORT", ui.get_web_client_port().to_string());

    match cmd.output() {
        Ok(o) => {
            let stdout = String::from_utf8_lossy(&o.stdout);
            let stderr = String::from_utf8_lossy(&o.stderr);
            if !stdout.is_empty() {
                append_log(ui, &stdout);
            }
            if !stderr.is_empty() {
                append_log(ui, &stderr);
            }
            if o.status.success() {
                append_log(ui, "Remove all completed. Physical data directory was NOT touched.");
            } else {
                append_log(ui, &format!("Command exited with: {:?}", o.status.code()));
            }
        }
        Err(e) => {
            append_log(ui, &format!("Failed to run compose down: {}", e));
        }
    }

    refresh_status(ui);
}

fn append_log(ui: &MainWindow, msg: &str) {
    let mut current = ui.get_log_output().to_string();
    if !current.is_empty() {
        current.push('\n');
    }
    current.push_str(msg);
    ui.set_log_output(current.into());
}

fn copy_logs(ui: &MainWindow) {
    let content = ui.get_log_output().to_string();
    if content.is_empty() {
        append_log(ui, "No logs to copy.");
        return;
    }
    if let Err(e) = arboard::Clipboard::new().and_then(|mut c| c.set_text(content)) {
        append_log(ui, &format!("Failed to copy to clipboard: {}", e));
    } else {
        append_log(ui, "Logs copied to clipboard.");
    }
}

fn save_logs(ui: &MainWindow) {
    let content = ui.get_log_output().to_string();
    if content.is_empty() {
        append_log(ui, "No logs to save.");
        return;
    }

    if let Some(path) = rfd::FileDialog::new()
        .set_title("Save logs to file")
        .set_file_name("nkrypt-logs.txt")
        .save_file()
    {
        if let Err(e) = std::fs::write(&path, &content) {
            append_log(ui, &format!("Failed to save logs: {}", e));
        } else {
            append_log(ui, &format!("Logs saved to {}", path.display()));
        }
    }
}

fn user_service_dir() -> std::path::PathBuf {
    dirs::config_dir()
        .unwrap_or_else(|| std::path::PathBuf::from(
            std::env::var("HOME").unwrap_or_default()
        ).join(".config"))
        .join("systemd")
        .join("user")
}

fn user_service_path() -> std::path::PathBuf {
    user_service_dir().join("nkrypt-desktop.service")
}

/// Returns true if the user systemd unit is currently enabled.
fn autostart_is_enabled() -> bool {
    std::process::Command::new("systemctl")
        .args(["--user", "is-enabled", "nkrypt-desktop.service"])
        .output()
        .map(|o| {
            let stdout = String::from_utf8_lossy(&o.stdout);
            let s = stdout.trim();
            s == "enabled" || s == "enabled-runtime"
        })
        .unwrap_or(false)
}

fn toggle_autostart(ui: &MainWindow, enabled: bool) {
    let compose_dir = config::repo_root().join("desktop-bootstrapper");
    let compose_file = compose_dir.join("docker-compose.yml");
    let env_file = config::env_path();

    // Resolve the full path of the container runtime so the systemd unit
    // works even when the user's PATH is not fully set at login time.
    let docker_cmd = {
        let candidate = if docker::command_exists("docker") { "docker" } else { "podman" };
        std::process::Command::new("which")
            .arg(candidate)
            .output()
            .ok()
            .filter(|o| o.status.success())
            .and_then(|o| String::from_utf8(o.stdout).ok())
            .map(|s| s.trim().to_string())
            .unwrap_or_else(|| format!("/usr/bin/{}", candidate))
    };

    if enabled {
        // Build and write the unit file
        let service_content = format!(
            r#"[Unit]
Description=nkrypt.xyz Desktop Stack
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory={workdir}
EnvironmentFile={env_file}
ExecStart={cmd} compose -f {compose_file} up -d
ExecStop={cmd} compose -f {compose_file} down

[Install]
WantedBy=default.target
"#,
            workdir = compose_dir.display(),
            env_file = env_file.display(),
            cmd = docker_cmd,
            compose_file = compose_file.display(),
        );

        let service_dir = user_service_dir();
        if let Err(e) = std::fs::create_dir_all(&service_dir) {
            append_log(ui, &format!("Auto-start: failed to create systemd user dir: {}", e));
            ui.set_autostart_enabled(false);
            return;
        }

        let service_path = user_service_path();
        if let Err(e) = std::fs::write(&service_path, &service_content) {
            append_log(ui, &format!("Auto-start: failed to write unit file: {}", e));
            ui.set_autostart_enabled(false);
            return;
        }
        append_log(ui, &format!("Auto-start: wrote {}", service_path.display()));

        // daemon-reload
        let reload = std::process::Command::new("systemctl")
            .args(["--user", "daemon-reload"])
            .output();
        match reload {
            Ok(o) if o.status.success() => {}
            Ok(o) => {
                append_log(ui, &format!(
                    "Auto-start: daemon-reload failed: {}",
                    String::from_utf8_lossy(&o.stderr).trim()
                ));
                ui.set_autostart_enabled(false);
                return;
            }
            Err(e) => {
                append_log(ui, &format!("Auto-start: daemon-reload error: {}", e));
                ui.set_autostart_enabled(false);
                return;
            }
        }

        // enable
        let enable = std::process::Command::new("systemctl")
            .args(["--user", "enable", "nkrypt-desktop.service"])
            .output();
        match enable {
            Ok(o) if o.status.success() => {
                append_log(ui, "Auto-start: enabled. Services will start automatically on login.");
            }
            Ok(o) => {
                append_log(ui, &format!(
                    "Auto-start: enable failed: {}",
                    String::from_utf8_lossy(&o.stderr).trim()
                ));
                ui.set_autostart_enabled(false);
            }
            Err(e) => {
                append_log(ui, &format!("Auto-start: enable error: {}", e));
                ui.set_autostart_enabled(false);
            }
        }
    } else {
        // disable
        let disable = std::process::Command::new("systemctl")
            .args(["--user", "disable", "nkrypt-desktop.service"])
            .output();
        match disable {
            Ok(o) if o.status.success() => {
                append_log(ui, "Auto-start: disabled.");
            }
            Ok(o) => {
                append_log(ui, &format!(
                    "Auto-start: disable failed: {}",
                    String::from_utf8_lossy(&o.stderr).trim()
                ));
            }
            Err(e) => {
                append_log(ui, &format!("Auto-start: disable error: {}", e));
            }
        }

        // Remove the unit file so there's no stale service
        let service_path = user_service_path();
        if service_path.exists() {
            if let Err(e) = std::fs::remove_file(&service_path) {
                append_log(ui, &format!("Auto-start: could not remove unit file: {}", e));
            } else {
                append_log(ui, &format!("Auto-start: removed {}", service_path.display()));
            }
            let _ = std::process::Command::new("systemctl")
                .args(["--user", "daemon-reload"])
                .output();
        }

        ui.set_autostart_enabled(false);
    }
}
