CREATE TABLE IF NOT EXISTS sessions (
    id              CHAR(16) PRIMARY KEY,
    user_id         CHAR(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_hash    TEXT NOT NULL,
    has_expired     BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at      TIMESTAMPTZ,
    expire_reason   VARCHAR(256),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_api_key_hash ON sessions(api_key_hash);

