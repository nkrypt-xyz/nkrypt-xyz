# nkrypt.xyz Desktop Bootstrapper

Linux desktop app that manages a self-hosted nkrypt.xyz stack via a GUI — no terminal required.

## Prerequisites

- **Docker** or **Podman** with Compose plugin
- **Rust** (build only)
- **fontconfig-devel** — `sudo dnf install fontconfig-devel` (Fedora) / `libfontconfig-dev` (Debian/Ubuntu)

## Quick Start

```bash
# Build
cd desktop-bootstrapper
make build

# Run (release)
make run

# Run (debug / development)
make dev
```

Binary: `target/release/nkrypt-desktop-bootstrapper`

## What it manages

| Service | Default port |
|---------|-------------|
| PostgreSQL | 9200 |
| Redis | 9201 |
| MinIO | 9202 |
| MinIO Console | 9203 |
| Web Server | 9204 |
| Web Client | 9205 |

All ports are configurable from the Config tab. Configuration is written to `.env` and consumed by `docker-compose.yml`.

## Features

- **Docker/Podman check** — detects the runtime on startup; shows install instructions if missing
- **Data directory** — browse or type a path; subdirs (`postgres/`, `redis/`, `minio/`) are created and made world-writable automatically (required for rootless Podman)
- **Start / Stop** — orchestrates a four-phase startup: stop deps → start deps → run migrations (via `psql` inside the container) → start web server + client
- **Service status** — live table with green/red indicators, port numbers, and Open buttons; auto-refreshes every 10 s
- **Logs** — streaming output from all phases; copy to clipboard or save to file
- **Auto-start** — writes a `~/.config/systemd/user/nkrypt-desktop.service` unit and enables it; no sudo required

## Notes

- Migrations run automatically on every `Start` using `psql` inside the PostgreSQL container — no `golang-migrate` binary needed on the host.
- The `.env` file is generated and managed by the app. Copy `.env.example` to bootstrap defaults, but the app will overwrite it on save.
- On rootless Podman, data directories are `chmod 777`'d so container UIDs can write. A warning dialog appears if this fails.

---

**Need help?** See [dev-docs/](../dev-docs/) or open an [issue](https://github.com/nkrypt-xyz/nkrypt-xyz-web-server/issues).
