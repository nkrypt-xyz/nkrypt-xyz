CREATE TABLE IF NOT EXISTS bucket_user_permissions (
    id                      BIGSERIAL PRIMARY KEY,
    bucket_id               CHAR(16) NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    user_id                 CHAR(16) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notes                   VARCHAR(256) NOT NULL DEFAULT '',

    perm_modify             BOOLEAN NOT NULL DEFAULT FALSE,
    perm_manage_authorization BOOLEAN NOT NULL DEFAULT FALSE,
    perm_destroy            BOOLEAN NOT NULL DEFAULT FALSE,
    perm_view_content       BOOLEAN NOT NULL DEFAULT FALSE,
    perm_manage_content     BOOLEAN NOT NULL DEFAULT FALSE,

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(bucket_id, user_id)
);

CREATE INDEX idx_bup_bucket_id ON bucket_user_permissions(bucket_id);
CREATE INDEX idx_bup_user_id ON bucket_user_permissions(user_id);

