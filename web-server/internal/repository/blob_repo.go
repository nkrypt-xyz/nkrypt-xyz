package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type BlobRepository struct {
	db *pgxpool.Pool
}

func NewBlobRepository(db *pgxpool.Pool) *BlobRepository {
	return &BlobRepository{db: db}
}

func (r *BlobRepository) FindByID(ctx context.Context, blobID string) (*model.Blob, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, file_id, crypto_meta_header_content, started_at, finished_at,
		       status, created_by_user_id, created_at, updated_at
		FROM blobs WHERE id=$1
	`, blobID)
	var b model.Blob
	if err := row.Scan(
		&b.ID, &b.BucketID, &b.FileID, &b.CryptoMetaHeaderContent, &b.StartedAt, &b.FinishedAt,
		&b.Status, &b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BlobRepository) FindInProgressBlob(ctx context.Context, bucketID, fileID, blobID string) (*model.Blob, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, file_id, crypto_meta_header_content, started_at, finished_at,
		       status, created_by_user_id, created_at, updated_at
		FROM blobs
		WHERE bucket_id=$1 AND file_id=$2 AND id=$3 AND status='started'
	`, bucketID, fileID, blobID)
	var b model.Blob
	if err := row.Scan(
		&b.ID, &b.BucketID, &b.FileID, &b.CryptoMetaHeaderContent, &b.StartedAt, &b.FinishedAt,
		&b.Status, &b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BlobRepository) FindLatestFinishedBlob(ctx context.Context, bucketID, fileID string) (*model.Blob, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, file_id, crypto_meta_header_content, started_at, finished_at,
		       status, created_by_user_id, created_at, updated_at
		FROM blobs
		WHERE bucket_id=$1 AND file_id=$2 AND status='finished'
		ORDER BY finished_at DESC LIMIT 1
	`, bucketID, fileID)
	var b model.Blob
	if err := row.Scan(
		&b.ID, &b.BucketID, &b.FileID, &b.CryptoMetaHeaderContent, &b.StartedAt, &b.FinishedAt,
		&b.Status, &b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BlobRepository) Create(ctx context.Context, b *model.Blob) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO blobs (id, bucket_id, file_id, crypto_meta_header_content, status, created_by_user_id)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, b.ID, b.BucketID, b.FileID, b.CryptoMetaHeaderContent, b.Status, b.CreatedByUserID)
	return err
}

func (r *BlobRepository) MarkFinished(ctx context.Context, blobID string) error {
	_, err := r.db.Exec(ctx, `UPDATE blobs SET status='finished', finished_at=NOW(), updated_at=NOW() WHERE id=$1`, blobID)
	return err
}

func (r *BlobRepository) MarkErroneous(ctx context.Context, blobID string) error {
	_, err := r.db.Exec(ctx, `UPDATE blobs SET status='error', updated_at=NOW() WHERE id=$1`, blobID)
	return err
}

func (r *BlobRepository) ListBlobsForFile(ctx context.Context, bucketID, fileID string) ([]model.Blob, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, file_id, crypto_meta_header_content, started_at, finished_at,
		       status, created_by_user_id, created_at, updated_at
		FROM blobs
		WHERE bucket_id=$1 AND file_id=$2
	`, bucketID, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Blob
	for rows.Next() {
		var b model.Blob
		if err := rows.Scan(
			&b.ID, &b.BucketID, &b.FileID, &b.CryptoMetaHeaderContent, &b.StartedAt, &b.FinishedAt,
			&b.Status, &b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BlobRepository) ListBlobsForFileExcluding(ctx context.Context, bucketID, fileID, excludeBlobID string) ([]model.Blob, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, file_id, crypto_meta_header_content, started_at, finished_at,
		       status, created_by_user_id, created_at, updated_at
		FROM blobs
		WHERE bucket_id=$1 AND file_id=$2 AND id != $3
	`, bucketID, fileID, excludeBlobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Blob
	for rows.Next() {
		var b model.Blob
		if err := rows.Scan(
			&b.ID, &b.BucketID, &b.FileID, &b.CryptoMetaHeaderContent, &b.StartedAt, &b.FinishedAt,
			&b.Status, &b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BlobRepository) DeleteBlob(ctx context.Context, blobID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM blobs WHERE id=$1`, blobID)
	return err
}

func (r *BlobRepository) DeleteBlobsForFileExcluding(ctx context.Context, bucketID, fileID, keepBlobID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM blobs WHERE bucket_id=$1 AND file_id=$2 AND id != $3`, bucketID, fileID, keepBlobID)
	return err
}

func (r *BlobRepository) DeleteAllBlobsForFile(ctx context.Context, bucketID, fileID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM blobs WHERE bucket_id=$1 AND file_id=$2`, bucketID, fileID)
	return err
}
