#!/usr/bin/env bash
#===============================================================================
# nkrypt.xyz — Debian 12 All-in-One Installer
#===============================================================================
#
# Installs the full nkrypt.xyz stack (PostgreSQL, Redis, MinIO, web-server,
# web-client) on a fresh Debian 12 (Bookworm) machine and wires everything
# together with Caddy as the reverse proxy / automatic-TLS frontend.
#
# Usage:
#   sudo ./install-debian-12.sh \
#       --api-domain    api.example.com \
#       --client-domain app.example.com \
#       --admin-password 'ChangeMeNow!' \
#       [--db-password <pw>]           \
#       [--minio-secret <key>]         \
#       [--acme-email   you@example.com]
#
# Uninstall:
#   sudo ./install-debian-12.sh --uninstall
#
# Required flags:
#   --api-domain      Public domain for the backend API
#   --client-domain   Public domain for the web client SPA
#   --admin-password  Initial admin account password
#
# Optional flags:
#   --db-password       PostgreSQL password  (auto-generated if omitted)
#   --minio-secret      MinIO secret key     (auto-generated if omitted)
#   --acme-email        Email for ACME/Let's Encrypt registration
#   --uninstall         Remove nkrypt.xyz (data is NOT deleted)
#
# Note: If your server has a dynamic IP, configure ddclient (or equivalent)
#       BEFORE running this script so your domains resolve to this machine.
#
# Prerequisites:
#   - Debian 12 (Bookworm) — run as root or with sudo
#   - The repo must be checked out; this script lives at
#     <repo-root>/infrastructure/bare-metal/install-debian-12.sh
#===============================================================================

set -euo pipefail

#-------------------------------------------------------------------------------
# Paths
#-------------------------------------------------------------------------------
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
readonly INSTALL_DIR="/opt/nkrypt"
readonly LOG_FILE="/var/log/nkrypt-install-$(date +%Y%m%d-%H%M%S).log"

#-------------------------------------------------------------------------------
# Versions (override via env vars if needed)
#-------------------------------------------------------------------------------
GO_VERSION="${GO_VERSION:-1.22.4}"
NODE_VERSION="${NODE_VERSION:-20}"
POSTGRES_VERSION="${POSTGRES_VERSION:-16}"
MIGRATE_VERSION="${MIGRATE_VERSION:-v4.17.0}"

#-------------------------------------------------------------------------------
# Colors
#-------------------------------------------------------------------------------
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
BLUE='\033[0;34m'; CYAN='\033[0;36m'; NC='\033[0m'

#-------------------------------------------------------------------------------
# Parameters (set by parse_args)
#-------------------------------------------------------------------------------
API_DOMAIN=""
CLIENT_DOMAIN=""
ADMIN_PASSWORD=""
DB_PASSWORD=""
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY=""
ACME_EMAIL=""
UNINSTALL=false

#-------------------------------------------------------------------------------
# Logging
#-------------------------------------------------------------------------------
log()         { echo -e "$1"; }
log_info()    { log "${BLUE}[INFO]${NC}    $1"; }
log_success() { log "${GREEN}[OK]${NC}      $1"; }
log_warn()    { log "${YELLOW}[WARN]${NC}    $1"; }
log_error()   { log "${RED}[ERROR]${NC}   $1"; }
log_section() { log "\n${CYAN}━━━  $1  ━━━${NC}\n"; }

die() { log_error "$1"; exit 1; }

# Run a command and show its output live.
# main() sets up "exec > >(tee -a $LOG_FILE) 2>&1" so all output already goes
# to both terminal and log — no need for an explicit tee pipe here.
# We use "|| rc=$?" with set -e active: the || operator suppresses errexit on
# the left-hand side, so a failing command is caught cleanly in rc.
run() {
    log_info "▶ $*"
    local rc=0
    "$@" || rc=$?
    [[ "$rc" -eq 0 ]] || die "Command failed (exit $rc): $*  (see ${LOG_FILE})"
}

# Like run() but non-fatal
try() { "$@" 2>/dev/null || true; }

#-------------------------------------------------------------------------------
# Argument parsing
#-------------------------------------------------------------------------------
usage() {
    sed -n '3,57p' "$0" | sed 's/^#//'
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --api-domain)         API_DOMAIN="$2";          shift 2 ;;
            --client-domain)      CLIENT_DOMAIN="$2";       shift 2 ;;
            --admin-password)     ADMIN_PASSWORD="$2";      shift 2 ;;
            --db-password)        DB_PASSWORD="$2";         shift 2 ;;
            --minio-secret)       MINIO_SECRET_KEY="$2";    shift 2 ;;
            --acme-email)         ACME_EMAIL="$2";          shift 2 ;;
            --uninstall)          UNINSTALL=true;           shift   ;;
            -h|--help)            usage; exit 0 ;;
            *) die "Unknown option: $1  (use --help)" ;;
        esac
    done
}

#-------------------------------------------------------------------------------
# Validation
#-------------------------------------------------------------------------------
validate_args() {
    if [[ "$UNINSTALL" == true ]]; then return; fi

    [[ -n "$API_DOMAIN"    ]] || die "--api-domain is required"
    [[ -n "$CLIENT_DOMAIN" ]] || die "--client-domain is required"
    [[ -n "$ADMIN_PASSWORD" ]] || die "--admin-password is required"

    # Auto-generate secrets that weren't supplied
    if [[ -z "$DB_PASSWORD" ]]; then
        DB_PASSWORD="$(openssl rand -hex 16)"
        log_warn "DB password not provided — generated: ${DB_PASSWORD}"
        log_warn "Save this! You will need it if you re-run migrations or inspect the database."
    fi
    if [[ -z "$MINIO_SECRET_KEY" ]]; then
        MINIO_SECRET_KEY="$(openssl rand -hex 20)"
        log_warn "MinIO secret not provided — generated: ${MINIO_SECRET_KEY}"
    fi
}

#-------------------------------------------------------------------------------
# System checks
#-------------------------------------------------------------------------------
check_root() {
    [[ $EUID -eq 0 ]] || die "Run as root (sudo $0 ...)"
}

check_debian() {
    [[ -f /etc/debian_version ]] || die "Debian 12 required. Detected: $(uname -a)"
    local codename
    codename="$(. /etc/os-release && echo "${VERSION_CODENAME:-}")"
    if [[ "$codename" != "bookworm" ]]; then
        log_warn "Tested on Debian 12 (bookworm). Detected: ${codename}. Proceeding anyway."
    fi
}

#-------------------------------------------------------------------------------
# Uninstall
#-------------------------------------------------------------------------------
do_uninstall() {
    log_section "Uninstalling nkrypt.xyz"

    for svc in nkrypt-server minio; do
        try systemctl stop    "$svc"
        try systemctl disable "$svc"
    done

    rm -f /etc/systemd/system/nkrypt-server.service
    rm -f /etc/systemd/system/minio.service
    try systemctl daemon-reload

    # Reset Caddy to empty config (don't remove Caddy itself)
    if [[ -f /etc/caddy/Caddyfile ]]; then
        echo "# nkrypt removed" > /etc/caddy/Caddyfile
        try systemctl reload caddy
    fi

    # Backup cron
    try bash -c '(crontab -l 2>/dev/null | grep -v backup-nkrypt) | crontab -'
    rm -f /usr/local/bin/backup-nkrypt.sh

    # Application directory
    rm -rf "$INSTALL_DIR"

    # Users
    try userdel -r nkrypt 2>/dev/null
    try userdel -r minio  2>/dev/null

    log_success "Uninstall complete."
    log_info "Data NOT removed: /var/lib/postgresql  /var/lib/redis  /mnt/minio"
    log_info "To purge all data: rm -rf /var/lib/postgresql /var/lib/redis /mnt/minio"
    exit 0
}

#-------------------------------------------------------------------------------
# System packages
#-------------------------------------------------------------------------------
install_system_packages() {
    log_section "System packages"
    run apt-get update -qq
    run apt-get install -y \
        git wget curl tar gzip gnupg2 lsb-release ca-certificates \
        apt-transport-https software-properties-common \
        ufw fail2ban unattended-upgrades apt-listchanges \
        openssl
    export PATH="/usr/local/bin:$PATH"
}

#-------------------------------------------------------------------------------
# Go
#-------------------------------------------------------------------------------
install_go() {
    log_section "Go ${GO_VERSION}"

    if command -v go &>/dev/null; then
        local v; v="$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')"
        if [[ "$(printf '%s\n' "$v" "1.21" | sort -V | head -1)" == "1.21" ]]; then
            log_success "Go already installed: $(go version)"; return
        fi
    fi

    local tarball="go${GO_VERSION}.linux-amd64.tar.gz"
    run wget -q "https://go.dev/dl/${tarball}" -O "/tmp/${tarball}"
    run rm -rf /usr/local/go
    run tar -C /usr/local -xzf "/tmp/${tarball}"
    rm -f "/tmp/${tarball}"

    cat > /etc/profile.d/go.sh <<'EOF'
export PATH="$PATH:/usr/local/go/bin"
EOF
    export PATH="$PATH:/usr/local/go/bin"
    log_success "Go installed: $(/usr/local/go/bin/go version)"
}

#-------------------------------------------------------------------------------
# Node.js
#-------------------------------------------------------------------------------
install_nodejs() {
    log_section "Node.js ${NODE_VERSION}"

    if command -v node &>/dev/null; then
        local v; v="$(node -v | grep -oP 'v\K[0-9]+')"
        [[ "$v" -ge 18 ]] && { log_success "Node.js $(node -v) already installed"; return; }
    fi

    run bash -c "curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -"
    run apt-get install -y nodejs
    log_success "Node.js $(node -v) / npm $(npm -v)"
}

#-------------------------------------------------------------------------------
# PostgreSQL
#-------------------------------------------------------------------------------
install_postgresql() {
    log_section "PostgreSQL ${POSTGRES_VERSION}"

    if ! command -v psql &>/dev/null; then
        run bash -c "echo 'deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main' \
                     > /etc/apt/sources.list.d/pgdg.list"
        run bash -c "wget -qO- https://www.postgresql.org/media/keys/ACCC4CF8.asc \
                     | gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg"
        run apt-get update -qq
        run apt-get install -y "postgresql-${POSTGRES_VERSION}" "postgresql-client-${POSTGRES_VERSION}"
    fi

    run systemctl enable --now postgresql

    # Create role + database if they don't exist; always sync the password so
    # re-runs with a different --db-password don't break authentication.
    if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='nkrypt'" | grep -q 1; then
        sudo -u postgres psql >> "$LOG_FILE" 2>&1 <<SQL
CREATE DATABASE nkrypt;
CREATE USER nkrypt WITH ENCRYPTED PASSWORD '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON DATABASE nkrypt TO nkrypt;
ALTER DATABASE nkrypt OWNER TO nkrypt;
\c nkrypt
GRANT ALL ON SCHEMA public TO nkrypt;
SQL
        log_success "PostgreSQL role and database created"
    else
        # Role already exists — update password to match current run
        sudo -u postgres psql -c "ALTER USER nkrypt WITH ENCRYPTED PASSWORD '${DB_PASSWORD}';" \
            >> "$LOG_FILE" 2>&1
        log_info "PostgreSQL role already exists — password updated"
    fi

    local hba="/etc/postgresql/${POSTGRES_VERSION}/main/pg_hba.conf"
    if ! grep -q "host.*nkrypt.*nkrypt.*127.0.0.1" "$hba" 2>/dev/null; then
        cat >> "$hba" <<EOF
host    nkrypt          nkrypt          127.0.0.1/32            scram-sha-256
local   nkrypt          nkrypt                                  scram-sha-256
EOF
        run systemctl restart postgresql
    fi

    log_success "PostgreSQL ready"
}

#-------------------------------------------------------------------------------
# Redis
#-------------------------------------------------------------------------------
install_redis() {
    log_section "Redis"

    if ! command -v redis-server &>/dev/null; then
        run apt-get install -y redis-server
    fi
    run systemctl enable --now redis-server

    redis-cli ping | grep -q PONG || die "Redis ping failed"
    log_success "Redis ready"
}

#-------------------------------------------------------------------------------
# MinIO
#-------------------------------------------------------------------------------
install_minio() {
    log_section "MinIO"

    if [[ ! -x /usr/local/bin/minio ]]; then
        run wget -q "https://dl.min.io/server/minio/release/linux-amd64/minio" \
            -O /usr/local/bin/minio
        chmod +x /usr/local/bin/minio
    fi

    id minio &>/dev/null || run useradd -r -s /usr/sbin/nologin minio
    run mkdir -p /mnt/minio/data /etc/minio
    run chown -R minio:minio /mnt/minio

    cat > /etc/minio/minio.env <<EOF
MINIO_ROOT_USER=${MINIO_ACCESS_KEY}
MINIO_ROOT_PASSWORD=${MINIO_SECRET_KEY}
MINIO_VOLUMES=/mnt/minio/data
MINIO_OPTS="--address :9000 --console-address :9001"
EOF
    chmod 600 /etc/minio/minio.env

    cat > /etc/systemd/system/minio.service <<'UNIT'
[Unit]
Description=MinIO Object Storage
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=minio
Group=minio
EnvironmentFile=/etc/minio/minio.env
ExecStart=/usr/local/bin/minio server $MINIO_OPTS $MINIO_VOLUMES
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
UNIT

    run systemctl daemon-reload
    run systemctl enable minio

    # MinIO bakes root credentials into its data dir on first init.
    # If it fails to start (credentials changed across re-runs), wipe the data
    # dir and let it reinitialize — on a test/fresh install this is safe.
    local minio_rc=0
    systemctl restart minio || minio_rc=$?
    if [[ "$minio_rc" -ne 0 ]]; then
        log_warn "MinIO failed to restart (exit $minio_rc) — wiping data dir and reinitializing"
        systemctl stop minio 2>/dev/null || true
        rm -rf /mnt/minio/data/*
        run systemctl start minio
    fi

    local i=0
    while ! curl -sf http://localhost:9000/minio/health/live &>/dev/null; do
        i=$(( i + 1 )); [[ $i -lt 30 ]] || die "MinIO did not become healthy"
        sleep 2
    done
    log_success "MinIO ready"
}

#-------------------------------------------------------------------------------
# MinIO client (mc) — used to create the initial bucket
#-------------------------------------------------------------------------------
install_mc() {
    command -v mc &>/dev/null && return
    run wget -q "https://dl.min.io/client/mc/release/linux-amd64/mc" -O /usr/local/bin/mc
    chmod +x /usr/local/bin/mc
}

#-------------------------------------------------------------------------------
# golang-migrate
#-------------------------------------------------------------------------------
install_migrate() {
    log_section "golang-migrate ${MIGRATE_VERSION}"
    command -v migrate &>/dev/null && { log_success "migrate already installed"; return; }

    run wget -q \
        "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" \
        -O /tmp/migrate.tar.gz
    run tar xf /tmp/migrate.tar.gz -C /tmp
    run mv /tmp/migrate /usr/local/bin/migrate
    chmod +x /usr/local/bin/migrate
    rm -f /tmp/migrate.tar.gz
    log_success "migrate $(migrate -version 2>&1 | head -1)"
}

#-------------------------------------------------------------------------------
# Caddy
#-------------------------------------------------------------------------------
install_caddy() {
    log_section "Caddy"

    if ! command -v caddy &>/dev/null; then
        run apt-get install -y debian-keyring debian-archive-keyring apt-transport-https
        run bash -c "curl -1sLf \
            'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
            | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg"
        run bash -c "curl -1sLf \
            'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
            | tee /etc/apt/sources.list.d/caddy-stable.list"
        run apt-get update -qq
        run apt-get install -y caddy
    fi

    run systemctl enable caddy
    log_success "Caddy installed"
}

configure_caddy() {
    log_section "Caddy — write Caddyfile"

    local global_block=""
    if [[ -n "$ACME_EMAIL" ]]; then
        global_block="{\n    email ${ACME_EMAIL}\n}\n\n"
    fi

    # Logs go to journald by default (journalctl -u caddy).
    # File-based access logs can be added later once the service is stable.
    cat > /etc/caddy/Caddyfile <<CADDYFILE
$(printf '%b' "$global_block")# ── API (backend) ──────────────────────────────────────────────────────────
${API_DOMAIN} {
    reverse_proxy localhost:9041
}

# ── Web Client (SPA) ────────────────────────────────────────────────────────
${CLIENT_DOMAIN} {
    root * ${INSTALL_DIR}/web-client
    file_server
    try_files {path} /index.html
}
CADDYFILE

    run caddy validate --config /etc/caddy/Caddyfile
    run systemctl restart caddy
    log_success "Caddyfile written and Caddy started"
}

#-------------------------------------------------------------------------------
# Build and deploy backend
#-------------------------------------------------------------------------------
deploy_backend() {
    log_section "Backend — build & deploy"

    # User and directories
    id nkrypt &>/dev/null || run useradd -r -s /bin/bash -d "$INSTALL_DIR" nkrypt
    run mkdir -p "${INSTALL_DIR}"/{bin,migrations}
    run chown -R nkrypt:nkrypt "$INSTALL_DIR"

    # Build
    export PATH="$PATH:/usr/local/go/bin"

    log_info "Downloading Go module dependencies (this may take a few minutes on first run)…"
    run bash -c "cd '${REPO_ROOT}/web-server' && go mod download"

    log_info "Compiling web-server binary…"
    run bash -c "cd '${REPO_ROOT}/web-server' && \
        go build -v -buildvcs=false -ldflags='-s -w' -o '${INSTALL_DIR}/bin/nkrypt-server' ./cmd/server"
    run cp -r "${REPO_ROOT}/web-server/migrations/." "${INSTALL_DIR}/migrations/"
    run chown -R nkrypt:nkrypt "$INSTALL_DIR"
    run chmod 755 "${INSTALL_DIR}/bin/nkrypt-server"

    # Environment file
    cat > "${INSTALL_DIR}/.env" <<EOF
NK_DATABASE_URL=postgres://nkrypt:${DB_PASSWORD}@localhost:5432/nkrypt?sslmode=disable
NK_REDIS_ADDR=localhost:6379
NK_MINIO_ENDPOINT=localhost:9000
NK_MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY}
NK_MINIO_SECRET_KEY=${MINIO_SECRET_KEY}
NK_MINIO_BUCKET_NAME=nkrypt-blobs
NK_IAM_DEFAULT_ADMIN_PASSWORD=${ADMIN_PASSWORD}
NK_LOG_LEVEL=info
NK_LOG_FORMAT=json
NK_SERVER_HTTP_PORT=9041
EOF
    run chown nkrypt:nkrypt "${INSTALL_DIR}/.env"
    run chmod 600 "${INSTALL_DIR}/.env"

    # MinIO bucket
    install_mc
    try mc alias set nkryptlocal http://localhost:9000 \
        "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}"
    try mc mb nkryptlocal/nkrypt-blobs

    # Migrations
    log_info "Running database migrations…"
    run sudo -u nkrypt bash -c "
        export \$(grep -v '^#' '${INSTALL_DIR}/.env' | xargs)
        migrate -path '${INSTALL_DIR}/migrations' \
                -database \"\$NK_DATABASE_URL\" up
    "

    # Systemd unit
    cat > /etc/systemd/system/nkrypt-server.service <<UNIT
[Unit]
Description=nkrypt.xyz Web Server
After=network-online.target postgresql.service redis-server.service minio.service
Wants=network-online.target
Requires=postgresql.service redis-server.service minio.service

[Service]
Type=simple
User=nkrypt
Group=nkrypt
WorkingDirectory=${INSTALL_DIR}
EnvironmentFile=${INSTALL_DIR}/.env
ExecStart=${INSTALL_DIR}/bin/nkrypt-server
Restart=on-failure
RestartSec=5

# Hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${INSTALL_DIR}

[Install]
WantedBy=multi-user.target
UNIT

    run systemctl daemon-reload
    run systemctl enable --now nkrypt-server

    # Health check
    local i=0
    while ! curl -sf http://localhost:9041/healthz &>/dev/null; do
        i=$(( i + 1 )); [[ $i -lt 15 ]] || die "nkrypt-server did not start. Check: journalctl -u nkrypt-server"
        sleep 2
    done
    log_success "Backend running on localhost:9041"
}

#-------------------------------------------------------------------------------
# Build and deploy frontend
#-------------------------------------------------------------------------------
deploy_frontend() {
    log_section "Frontend — build & deploy"

    local spa_out="${INSTALL_DIR}/web-client"
    run mkdir -p "$spa_out"

    # Build core crypto module first (required by web-client)
    if [[ -d "${REPO_ROOT}/nkrypt-xyz-core-nodejs" ]]; then
        log_info "Building nkrypt-xyz-core-nodejs…"
        run bash -c "cd '${REPO_ROOT}/nkrypt-xyz-core-nodejs' && npm install && npm run build"
    fi

    log_info "Building web-client (Vue/Quasar)…"
    run bash -c "cd '${REPO_ROOT}/web-client' && npm install"
    run bash -c "cd '${REPO_ROOT}/web-client' && \
        VITE_DEFAULT_SERVER_URL=https://${API_DOMAIN} \
        VITE_CLIENT_APPLICATION_NAME=nkrypt-web-client \
        npm run build"

    run cp -r "${REPO_ROOT}/web-client/dist/spa/." "$spa_out/"
    run chown -R nkrypt:nkrypt "$spa_out"
    log_success "Frontend deployed to ${spa_out}"
}

#-------------------------------------------------------------------------------
# Firewall (ufw)
#-------------------------------------------------------------------------------
configure_firewall() {
    log_section "Firewall (ufw)"
    run ufw --force reset
    run ufw default deny incoming
    run ufw default allow outgoing
    run ufw allow ssh
    run ufw allow http
    run ufw allow https
    echo "y" | ufw enable >> "$LOG_FILE" 2>&1 || true
    log_success "ufw: SSH + HTTP + HTTPS allowed"
}

#-------------------------------------------------------------------------------
# fail2ban
#-------------------------------------------------------------------------------
configure_fail2ban() {
    log_section "fail2ban"
    run systemctl enable --now fail2ban
    log_success "fail2ban active"
}

#-------------------------------------------------------------------------------
# Automatic security updates
#-------------------------------------------------------------------------------
configure_auto_updates() {
    log_section "Automatic security updates"

    cat > /etc/apt/apt.conf.d/50unattended-upgrades <<'EOF'
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
};
Unattended-Upgrade::AutoFixInterruptedDpkg "true";
Unattended-Upgrade::MinimalSteps "true";
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Automatic-Reboot "false";
EOF

    cat > /etc/apt/apt.conf.d/20auto-upgrades <<'EOF'
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
EOF

    log_success "Automatic security updates enabled"
}

#-------------------------------------------------------------------------------
# Daily PostgreSQL backup
#-------------------------------------------------------------------------------
configure_backups() {
    log_section "Daily database backups"

    cat > /usr/local/bin/backup-nkrypt.sh <<'SCRIPT'
#!/usr/bin/env bash
set -euo pipefail
DATE=$(date +%Y%m%d_%H%M%S)
DIR=/var/backups/nkrypt
mkdir -p "$DIR"
sudo -u postgres pg_dump nkrypt | gzip > "$DIR/db_${DATE}.sql.gz"
find "$DIR" -type f -name 'db_*.sql.gz' -mtime +7 -delete
echo "Backup: $DIR/db_${DATE}.sql.gz"
SCRIPT

    chmod +x /usr/local/bin/backup-nkrypt.sh
    (crontab -l 2>/dev/null | grep -v backup-nkrypt || true; \
     echo "0 2 * * * /usr/local/bin/backup-nkrypt.sh") | crontab -
    log_success "Backups scheduled daily at 02:00 → /var/backups/nkrypt/"
}

#-------------------------------------------------------------------------------
# Summary
#-------------------------------------------------------------------------------
print_summary() {
    log_section "Installation complete"
    log ""
    log "${GREEN}  Web Client  →  https://${CLIENT_DOMAIN}${NC}"
    log "${GREEN}  API         →  https://${API_DOMAIN}${NC}"
    log "${GREEN}  MinIO UI    →  http://localhost:9001  (internal)${NC}"
    log ""
    log_info "Default admin account: admin / (the --admin-password you supplied)"
    log ""
    log_info "Useful commands:"
    log "  systemctl status nkrypt-server postgresql redis-server minio caddy"
    log "  journalctl -u nkrypt-server -f"
    log "  journalctl -u caddy -f"
    log ""
    log_info "Log file: ${LOG_FILE}"
    log ""
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------
main() {
    # Redirect ALL output to log + terminal from this point on.
    # run() and try() rely on this — they do NOT pipe through tee themselves.
    exec > >(tee -a "$LOG_FILE") 2>&1

    echo -e "${CYAN}"
    echo "  ╔══════════════════════════════════════════════════╗"
    echo "  ║      nkrypt.xyz  —  Debian 12 Installer          ║"
    echo "  ╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"

    parse_args "$@"

    if [[ "$UNINSTALL" == true ]]; then
        check_root
        do_uninstall
    fi

    check_root
    check_debian
    validate_args

    log_info "Log: ${LOG_FILE}"

    # ── Dependencies ──────────────────────────────────────────────────
    install_system_packages
    install_go
    install_nodejs
    install_postgresql
    install_redis
    install_minio
    install_migrate
    install_caddy

    # ── Application ───────────────────────────────────────────────────
    deploy_backend
    deploy_frontend

    # ── Reverse proxy ─────────────────────────────────────────────────
    configure_caddy

    # ── Infrastructure ────────────────────────────────────────────────
    configure_firewall
    configure_fail2ban
    configure_auto_updates
    configure_backups

    print_summary
}

main "$@"
