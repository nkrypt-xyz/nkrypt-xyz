package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) FindByID(ctx context.Context, bucketID, fileID string) (*model.File, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       size_after_encryption_bytes, created_by_user_id, created_at, updated_at, content_updated_at
		FROM files WHERE bucket_id=$1 AND id=$2
	`, bucketID, fileID)
	var f model.File
	if err := row.Scan(
		&f.ID, &f.BucketID, &f.ParentDirectoryID, &f.Name, &f.MetaData, &f.EncryptedMetaData,
		&f.SizeAfterEncryptionBytes, &f.CreatedByUserID, &f.CreatedAt, &f.UpdatedAt, &f.ContentUpdatedAt,
	); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) ListByDirectory(ctx context.Context, bucketID, parentDirID string) ([]model.File, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       size_after_encryption_bytes, created_by_user_id, created_at, updated_at, content_updated_at
		FROM files
		WHERE bucket_id=$1 AND parent_directory_id=$2
		ORDER BY name
	`, bucketID, parentDirID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.File
	for rows.Next() {
		var f model.File
		if err := rows.Scan(
			&f.ID, &f.BucketID, &f.ParentDirectoryID, &f.Name, &f.MetaData, &f.EncryptedMetaData,
			&f.SizeAfterEncryptionBytes, &f.CreatedByUserID, &f.CreatedAt, &f.UpdatedAt, &f.ContentUpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func (r *FileRepository) FindByNameAndParent(ctx context.Context, name, bucketID, parentDirID string) (*model.File, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       size_after_encryption_bytes, created_by_user_id, created_at, updated_at, content_updated_at
		FROM files
		WHERE bucket_id=$1 AND parent_directory_id=$2 AND name=$3
	`, bucketID, parentDirID, name)
	var f model.File
	if err := row.Scan(
		&f.ID, &f.BucketID, &f.ParentDirectoryID, &f.Name, &f.MetaData, &f.EncryptedMetaData,
		&f.SizeAfterEncryptionBytes, &f.CreatedByUserID, &f.CreatedAt, &f.UpdatedAt, &f.ContentUpdatedAt,
	); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) Create(ctx context.Context, f *model.File) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO files (id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		                   size_after_encryption_bytes, created_by_user_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, f.ID, f.BucketID, f.ParentDirectoryID, f.Name, f.MetaData, f.EncryptedMetaData,
		f.SizeAfterEncryptionBytes, f.CreatedByUserID)
	return err
}

func (r *FileRepository) UpdateName(ctx context.Context, bucketID, fileID, name string) error {
	_, err := r.db.Exec(ctx, `UPDATE files SET name=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, fileID, name)
	return err
}

func (r *FileRepository) UpdateMetaData(ctx context.Context, bucketID, fileID string, metaData []byte) error {
	_, err := r.db.Exec(ctx, `UPDATE files SET meta_data=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, fileID, metaData)
	return err
}

func (r *FileRepository) UpdateEncryptedMetaData(ctx context.Context, bucketID, fileID, encryptedMetaData string) error {
	_, err := r.db.Exec(ctx, `UPDATE files SET encrypted_meta_data=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, fileID, encryptedMetaData)
	return err
}

func (r *FileRepository) UpdateSize(ctx context.Context, bucketID, fileID string, size int64) error {
	_, err := r.db.Exec(ctx, `UPDATE files SET size_after_encryption_bytes=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, fileID, size)
	return err
}

func (r *FileRepository) UpdateContentUpdatedAt(ctx context.Context, bucketID, fileID string) error {
	_, err := r.db.Exec(ctx, `UPDATE files SET content_updated_at=NOW(), updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, fileID)
	return err
}

func (r *FileRepository) Move(ctx context.Context, bucketID, fileID, newParentDirID, newName string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE files SET parent_directory_id=$3, name=$4, updated_at=NOW()
		WHERE bucket_id=$1 AND id=$2
	`, bucketID, fileID, newParentDirID, newName)
	return err
}

func (r *FileRepository) Delete(ctx context.Context, bucketID, fileID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM files WHERE bucket_id=$1 AND id=$2`, bucketID, fileID)
	return err
}
