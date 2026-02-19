CREATE TABLE IF NOT EXISTS users (
    id              CHAR(16) PRIMARY KEY,
    display_name    VARCHAR(128) NOT NULL,
    user_name       VARCHAR(32) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    password_salt   TEXT NOT NULL,
    is_banned       BOOLEAN NOT NULL DEFAULT FALSE,

    -- Global permissions stored as individual boolean columns
    perm_manage_all_user  BOOLEAN NOT NULL DEFAULT FALSE,
    perm_create_user      BOOLEAN NOT NULL DEFAULT FALSE,
    perm_create_bucket    BOOLEAN NOT NULL DEFAULT TRUE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_user_name ON users(user_name);

