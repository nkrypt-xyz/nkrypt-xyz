#!/usr/bin/env bash
#===============================================================================
# nkrypt.xyz — Amazon Linux 2 All-in-One Installer
#===============================================================================
#
# Installs the full nkrypt.xyz stack (PostgreSQL 16, Redis, MinIO, web-server,
# web-client) on Amazon Linux 2 and wires everything with Caddy (reverse proxy
# + automatic TLS).
#
# Usage:
#   sudo ./install-amazon-linux-2.sh \
#       --api-domain    api.example.com \
#       --client-domain app.example.com \
#       --admin-password 'ChangeMeNow!' \
#       [--db-password <pw>] [--minio-secret <key>] [--acme-email you@example.com]
#
# Uninstall:
#   sudo ./install-amazon-linux-2.sh --uninstall
#
# Required: --api-domain, --client-domain, --admin-password
# Optional: --db-password, --minio-secret, --acme-email, --uninstall
#
# Prerequisites:
#   - Amazon Linux 2 (x86_64) — run as root or with sudo
#   - Repo checked out; script at <repo>/infrastructure/bare-metal/install-amazon-linux-2.sh
#   - If dynamic IP: configure ddclient (or equivalent) before running.
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
# Versions
#-------------------------------------------------------------------------------
GO_VERSION="${GO_VERSION:-1.22.4}"
NODE_VERSION="${NODE_VERSION:-20}"
POSTGRES_VERSION="${POSTGRES_VERSION:-16}"
MIGRATE_VERSION="${MIGRATE_VERSION:-v4.17.0}"
CADDY_VERSION="${CADDY_VERSION:-2.8.4}"

#-------------------------------------------------------------------------------
# Colors
#-------------------------------------------------------------------------------
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
BLUE='\033[0;34m'; CYAN='\033[0;36m'; NC='\033[0m'

#-------------------------------------------------------------------------------
# Parameters
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

run() {
    log_info "▶ $*"
    local rc=0
    "$@" || rc=$?
    [[ "$rc" -eq 0 ]] || die "Command failed (exit $rc): $*  (see ${LOG_FILE})"
}

try() { "$@" 2>/dev/null || true; }

#-------------------------------------------------------------------------------
# Argument parsing
#-------------------------------------------------------------------------------
usage() {
    sed -n '3,45p' "$0" | sed 's/^#//'
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
    [[ -n "$API_DOMAIN"     ]] || die "--api-domain is required"
    [[ -n "$CLIENT_DOMAIN"  ]] || die "--client-domain is required"
    [[ -n "$ADMIN_PASSWORD" ]] || die "--admin-password is required"

    if [[ -z "$DB_PASSWORD" ]]; then
        DB_PASSWORD="$(openssl rand -hex 16)"
        log_warn "DB password not provided — generated: ${DB_PASSWORD}"
        log_warn "Save this if you re-run migrations or inspect the database."
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

check_amazon_linux_2() {
    [[ -f /etc/os-release ]] || die "Cannot detect OS (no /etc/os-release)"
    local id version
    id="$(. /etc/os-release && echo "${ID:-}")"
    version="$(. /etc/os-release && echo "${VERSION_ID:-}")"
    if [[ "$id" != "amzn" ]] || [[ "$version" != "2" ]]; then
        die "This script is for Amazon Linux 2. Detected: id=${id} version=${version}. Use install-debian-12.sh for Debian."
    fi
    if [[ "$(uname -m)" != "x86_64" ]]; then
        log_warn "Tested on x86_64 only. Arch: $(uname -m). Proceeding anyway."
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

    if [[ -f /etc/caddy/Caddyfile ]]; then
        echo "# nkrypt removed" > /etc/caddy/Caddyfile
        try systemctl restart caddy
    fi

    local new_cron; new_cron="$(crontab -l 2>/dev/null | grep -v backup-nkrypt || true)"
    if [[ -n "$new_cron" ]]; then echo "$new_cron" | crontab -; else crontab -r 2>/dev/null || true; fi
    rm -f /usr/local/bin/backup-nkrypt.sh

    rm -rf "$INSTALL_DIR"
    try userdel -r nkrypt 2>/dev/null
    try userdel -r minio  2>/dev/null

    log_success "Uninstall complete."
    log_info "Data NOT removed: PostgreSQL (/var/lib/pgsql), Redis (/var/lib/redis), MinIO (/mnt/minio)"
    exit 0
}

#-------------------------------------------------------------------------------
# System packages
#-------------------------------------------------------------------------------
install_system_packages() {
    log_section "System packages"

    run yum install -y git wget curl tar gzip openssl

    # EPEL (needed for Redis, fail2ban)
    if ! rpm -q epel-release &>/dev/null; then
        run amazon-linux-extras install -y epel
    fi

    run yum install -y firewalld fail2ban
    export PATH="/usr/local/bin:$PATH"
}

#-------------------------------------------------------------------------------
# Go
#-------------------------------------------------------------------------------
install_go() {
    log_section "Go ${GO_VERSION}"

    if command -v go &>/dev/null; then
        local v; v="$(go version | sed -n 's/.*go\([0-9]*\.[0-9]*\).*/\1/p')"
        if [[ -n "$v" ]] && [[ "$(printf '%s\n' "$v" "1.21" | sort -V | head -1)" == "1.21" ]]; then
            log_success "Go already installed: $(go version)"; return
        fi
    fi

    local tarball="go${GO_VERSION}.linux-amd64.tar.gz"
    run wget -q "https://go.dev/dl/${tarball}" -O "/tmp/${tarball}"
    run rm -rf /usr/local/go
    run tar -C /usr/local -xzf "/tmp/${tarball}"
    rm -f "/tmp/${tarball}"

    echo 'export PATH="$PATH:/usr/local/go/bin"' > /etc/profile.d/go.sh
    export PATH="$PATH:/usr/local/go/bin"
    log_success "Go installed: $(/usr/local/go/bin/go version)"
}

#-------------------------------------------------------------------------------
# Node.js (NodeSource repo for AL2)
#-------------------------------------------------------------------------------
install_nodejs() {
    log_section "Node.js ${NODE_VERSION}"

    if command -v node &>/dev/null; then
        local v; v="$(node -v | sed -n 's/^v\([0-9]*\).*/\1/p')"
        [[ -n "$v" ]] && [[ "$v" -ge 18 ]] && { log_success "Node.js $(node -v) already installed"; return; }
    fi

    run bash -c "curl -sL https://rpm.nodesource.com/setup_${NODE_VERSION}.x | bash -"
    run yum install -y nodejs
    log_success "Node.js $(node -v) / npm $(npm -v)"
}

#-------------------------------------------------------------------------------
# PostgreSQL 16 (PGDG repo for RHEL/CentOS 7 — AL2 compatible)
#-------------------------------------------------------------------------------
install_postgresql() {
    log_section "PostgreSQL ${POSTGRES_VERSION}"

    local pg_installed=false
    if command -v psql &>/dev/null; then
        local v; v="$(psql --version | sed -n 's/.* \([0-9]*\).*/\1/p')"
        [[ "$v" -ge 16 ]] && pg_installed=true
    fi

    if [[ "$pg_installed" != true ]]; then
        run yum install -y "https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm"
        run yum install -y "postgresql${POSTGRES_VERSION}-server" "postgresql${POSTGRES_VERSION}"

        if [[ ! -d /var/lib/pgsql/${POSTGRES_VERSION}/data/base ]]; then
            run "/usr/pgsql-${POSTGRES_VERSION}/bin/postgresql-${POSTGRES_VERSION}-setup" initdb
        fi
    fi

    run systemctl enable "postgresql-${POSTGRES_VERSION}"
    run systemctl start  "postgresql-${POSTGRES_VERSION}"

    if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='nkrypt'" 2>/dev/null | grep -q 1; then
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
        sudo -u postgres psql -c "ALTER USER nkrypt WITH ENCRYPTED PASSWORD '${DB_PASSWORD}';" >> "$LOG_FILE" 2>&1
        log_info "PostgreSQL role already exists — password updated"
    fi

    local hba="/var/lib/pgsql/${POSTGRES_VERSION}/data/pg_hba.conf"
    if ! grep -q "host.*nkrypt.*nkrypt.*127.0.0.1" "$hba" 2>/dev/null; then
        cat >> "$hba" <<EOF
host    nkrypt          nkrypt          127.0.0.1/32            md5
local   nkrypt          nkrypt                                  md5
EOF
        run systemctl restart "postgresql-${POSTGRES_VERSION}"
    fi

    log_success "PostgreSQL ready"
}

#-------------------------------------------------------------------------------
# Redis
#-------------------------------------------------------------------------------
install_redis() {
    log_section "Redis"

    if ! command -v redis-server &>/dev/null; then
        run yum install -y redis
    fi
    run systemctl enable --now redis

    redis-cli ping | grep -q PONG || die "Redis ping failed"
    log_success "Redis ready"
}

#-------------------------------------------------------------------------------
# MinIO
#-------------------------------------------------------------------------------
install_minio() {
    log_section "MinIO"

    if [[ ! -x /usr/local/bin/minio ]]; then
        run wget -q "https://dl.min.io/server/minio/release/linux-amd64/minio" -O /usr/local/bin/minio
        chmod +x /usr/local/bin/minio
    fi

    id minio &>/dev/null || run useradd -r -s /sbin/nologin minio
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
# MinIO client (mc)
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
# Caddy (binary install — no official yum package for AL2)
#-------------------------------------------------------------------------------
install_caddy() {
    log_section "Caddy"

    if command -v caddy &>/dev/null; then
        log_success "Caddy already installed: $(caddy version 2>&1 | head -1)"; return
    fi

    local caddy_tgz="caddy_${CADDY_VERSION}_linux_amd64.tar.gz"
    run wget -q "https://github.com/caddyserver/caddy/releases/download/v${CADDY_VERSION}/${caddy_tgz}" -O "/tmp/${caddy_tgz}"
    run tar -xzf "/tmp/${caddy_tgz}" -C /tmp
    run mv /tmp/caddy /usr/local/bin/caddy
    chmod +x /usr/local/bin/caddy
    rm -f "/tmp/${caddy_tgz}"

    id caddy &>/dev/null || run useradd -r -s /sbin/nologin -d /var/lib/caddy caddy
    run mkdir -p /etc/caddy
    run chown -R caddy:caddy /etc/caddy

    cat > /etc/systemd/system/caddy.service <<'UNIT'
[Unit]
Description=Caddy
Documentation=https://caddyserver.com/docs/
After=network-online.target
Wants=network-online.target

[Service]
User=caddy
Group=caddy
ExecStart=/usr/local/bin/caddy run --config /etc/caddy/Caddyfile
ExecReload=/usr/local/bin/caddy reload --config /etc/caddy/Caddyfile
TimeoutStopSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
UNIT

    run systemctl daemon-reload
    run systemctl enable caddy
    log_success "Caddy installed"
}

#-------------------------------------------------------------------------------
# Caddy — write Caddyfile
#-------------------------------------------------------------------------------
configure_caddy() {
    log_section "Caddy — write Caddyfile"

    local global_block=""
    if [[ -n "$ACME_EMAIL" ]]; then
        global_block="{\n    email ${ACME_EMAIL}\n}\n\n"
    fi

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

    run chown caddy:caddy /etc/caddy/Caddyfile
    run caddy validate --config /etc/caddy/Caddyfile
    run systemctl restart caddy
    log_success "Caddyfile written and Caddy started"
}

#-------------------------------------------------------------------------------
# Build and deploy backend
#-------------------------------------------------------------------------------
deploy_backend() {
    log_section "Backend — build & deploy"

    id nkrypt &>/dev/null || run useradd -r -s /bin/bash -d "$INSTALL_DIR" nkrypt
    run mkdir -p "${INSTALL_DIR}"/{bin,migrations}
    run chown -R nkrypt:nkrypt "$INSTALL_DIR"

    export PATH="$PATH:/usr/local/go/bin"

    log_info "Downloading Go module dependencies (this may take a few minutes on first run)…"
    run bash -c "cd '${REPO_ROOT}/web-server' && go mod download"

    log_info "Compiling web-server binary…"
    run bash -c "cd '${REPO_ROOT}/web-server' && \
        go build -v -buildvcs=false -ldflags='-s -w' -o '${INSTALL_DIR}/bin/nkrypt-server' ./cmd/server"
    run cp -r "${REPO_ROOT}/web-server/migrations/." "${INSTALL_DIR}/migrations/"
    run chown -R nkrypt:nkrypt "$INSTALL_DIR"
    run chmod 755 "${INSTALL_DIR}/bin/nkrypt-server"

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

    install_mc
    try mc alias set nkryptlocal http://localhost:9000 \
        "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}"
    try mc mb nkryptlocal/nkrypt-blobs

    log_info "Running database migrations…"
    run sudo -u nkrypt bash -c "
        export \$(grep -v '^#' '${INSTALL_DIR}/.env' | xargs)
        migrate -path '${INSTALL_DIR}/migrations' \
                -database \"\$NK_DATABASE_URL\" up
    "

    cat > /etc/systemd/system/nkrypt-server.service <<UNIT
[Unit]
Description=nkrypt.xyz Web Server
After=network-online.target postgresql-${POSTGRES_VERSION}.service redis.service minio.service
Wants=network-online.target
Requires=postgresql-${POSTGRES_VERSION}.service redis.service minio.service

[Service]
Type=simple
User=nkrypt
Group=nkrypt
WorkingDirectory=${INSTALL_DIR}
EnvironmentFile=${INSTALL_DIR}/.env
ExecStart=${INSTALL_DIR}/bin/nkrypt-server
Restart=on-failure
RestartSec=5
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
    # Caddy (running as caddy) needs read access to serve static files
    run chmod -R o+rX "$spa_out"
    log_success "Frontend deployed to ${spa_out}"
}

#-------------------------------------------------------------------------------
# Firewall (firewalld)
#-------------------------------------------------------------------------------
configure_firewall() {
    log_section "Firewall (firewalld)"

    run systemctl enable firewalld
    run systemctl start firewalld 2>/dev/null || true
    run firewall-cmd --permanent --set-default-zone=public
    run firewall-cmd --permanent --add-service=ssh
    run firewall-cmd --permanent --add-service=http
    run firewall-cmd --permanent --add-service=https
    run firewall-cmd --reload
    log_success "firewalld: SSH + HTTP + HTTPS allowed"
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
# Automatic security updates (yum-cron)
#-------------------------------------------------------------------------------
configure_auto_updates() {
    log_section "Automatic security updates"

    run yum install -y yum-cron
    if [[ -f /etc/yum/yum-cron.conf ]]; then
        sed -i 's/^apply_updates = .*/apply_updates = yes/' /etc/yum/yum-cron.conf
        sed -i 's/^update_cmd = .*/update_cmd = security/' /etc/yum/yum-cron.conf
    fi
    run systemctl enable --now yum-cron
    log_success "yum-cron enabled (security updates)"
}

#-------------------------------------------------------------------------------
# Daily PostgreSQL backup (use PGDG pg_dump path)
#-------------------------------------------------------------------------------
configure_backups() {
    log_section "Daily database backups"

    local pg_dump_bin="/usr/pgsql-${POSTGRES_VERSION}/bin/pg_dump"
    cat > /usr/local/bin/backup-nkrypt.sh <<SCRIPT
#!/usr/bin/env bash
set -euo pipefail
DATE=\$(date +%Y%m%d_%H%M%S)
DIR=/var/backups/nkrypt
mkdir -p "\$DIR"
sudo -u postgres ${pg_dump_bin} nkrypt | gzip > "\$DIR/db_\${DATE}.sql.gz"
find "\$DIR" -type f -name 'db_*.sql.gz' -mtime +7 -delete
echo "Backup: \$DIR/db_\${DATE}.sql.gz"
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
    log "  systemctl status nkrypt-server postgresql-${POSTGRES_VERSION} redis minio caddy"
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
    exec > >(tee -a "$LOG_FILE") 2>&1

    echo -e "${CYAN}"
    echo "  ╔══════════════════════════════════════════════════╗"
    echo "  ║   nkrypt.xyz  —  Amazon Linux 2 Installer         ║"
    echo "  ╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"

    parse_args "$@"

    if [[ "$UNINSTALL" == true ]]; then
        check_root
        do_uninstall
    fi

    check_root
    check_amazon_linux_2
    validate_args

    log_info "Log: ${LOG_FILE}"

    install_system_packages
    install_go
    install_nodejs
    install_postgresql
    install_redis
    install_minio
    install_migrate
    install_caddy

    deploy_backend
    deploy_frontend

    configure_caddy

    configure_firewall
    configure_fail2ban
    configure_auto_updates
    configure_backups

    print_summary
}

main "$@"
