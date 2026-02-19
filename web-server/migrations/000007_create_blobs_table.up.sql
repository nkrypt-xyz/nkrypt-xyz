CREATE TYPE blob_status AS ENUM ('started', 'finished', 'error');

CREATE TABLE IF NOT EXISTS blobs (
    id                          CHAR(16) PRIMARY KEY,
    bucket_id                   CHAR(16) NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    file_id                     CHAR(16) NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    crypto_meta_header_content  TEXT NOT NULL,
    started_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at                 TIMESTAMPTZ,
    status                      blob_status NOT NULL DEFAULT 'started',
    created_by_user_id          CHAR(16) NOT NULL REFERENCES users(id),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blobs_bucket_file ON blobs(bucket_id, file_id);
CREATE INDEX idx_blobs_file_id ON blobs(file_id);
CREATE INDEX idx_blobs_status ON blobs(status);
CREATE INDEX idx_blobs_finished_at ON blobs(finished_at DESC NULLS LAST);

