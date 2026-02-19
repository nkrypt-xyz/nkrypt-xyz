CREATE TABLE IF NOT EXISTS directories (
    id                      CHAR(16) PRIMARY KEY,
    bucket_id               CHAR(16) NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    parent_directory_id     CHAR(16) REFERENCES directories(id) ON DELETE CASCADE,
    name                    VARCHAR(256) NOT NULL,
    meta_data               JSONB NOT NULL DEFAULT '{}'::jsonb,
    encrypted_meta_data     TEXT NOT NULL DEFAULT '',
    created_by_user_id      CHAR(16) NOT NULL REFERENCES users(id),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Root directories have NULL parent_directory_id
    -- Non-root directories must be unique by name within their parent
    UNIQUE NULLS NOT DISTINCT (bucket_id, parent_directory_id, name)
);

CREATE INDEX idx_directories_bucket_id ON directories(bucket_id);
CREATE INDEX idx_directories_parent ON directories(parent_directory_id);
CREATE INDEX idx_directories_bucket_parent ON directories(bucket_id, parent_directory_id);

