CREATE TABLE IF NOT EXISTS buckets (
    id                      CHAR(16) PRIMARY KEY,
    name                    VARCHAR(64) NOT NULL UNIQUE,
    crypt_spec              VARCHAR(64) NOT NULL,
    crypt_data              VARCHAR(2048) NOT NULL,
    meta_data               JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by_user_id      CHAR(16) NOT NULL REFERENCES users(id),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_buckets_name ON buckets(name);
CREATE INDEX idx_buckets_created_by ON buckets(created_by_user_id);

