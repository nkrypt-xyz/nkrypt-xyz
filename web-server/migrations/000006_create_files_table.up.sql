CREATE TABLE IF NOT EXISTS files (
    id                          CHAR(16) PRIMARY KEY,
    bucket_id                   CHAR(16) NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    parent_directory_id         CHAR(16) NOT NULL REFERENCES directories(id) ON DELETE CASCADE,
    name                        VARCHAR(256) NOT NULL,
    meta_data                   JSONB NOT NULL DEFAULT '{}'::jsonb,
    encrypted_meta_data         TEXT NOT NULL DEFAULT '',
    size_after_encryption_bytes BIGINT NOT NULL DEFAULT 0,
    created_by_user_id          CHAR(16) NOT NULL REFERENCES users(id),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    content_updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(bucket_id, parent_directory_id, name)
);

CREATE INDEX idx_files_bucket_id ON files(bucket_id);
CREATE INDEX idx_files_parent_directory ON files(parent_directory_id);
CREATE INDEX idx_files_bucket_parent ON files(bucket_id, parent_directory_id);

