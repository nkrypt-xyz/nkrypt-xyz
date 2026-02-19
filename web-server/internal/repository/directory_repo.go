package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type DirectoryRepository struct {
	db *pgxpool.Pool
}

func NewDirectoryRepository(db *pgxpool.Pool) *DirectoryRepository {
	return &DirectoryRepository{db: db}
}

func (r *DirectoryRepository) FindByID(ctx context.Context, bucketID, id string) (*model.Directory, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       created_by_user_id, created_at, updated_at
		FROM directories WHERE bucket_id=$1 AND id=$2
	`, bucketID, id)
	var d model.Directory
	if err := row.Scan(
		&d.ID, &d.BucketID, &d.ParentDirectoryID, &d.Name, &d.MetaData, &d.EncryptedMetaData,
		&d.CreatedByUserID, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DirectoryRepository) FindByNameAndParent(ctx context.Context, name, bucketID string, parentDirID *string) (*model.Directory, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       created_by_user_id, created_at, updated_at
		FROM directories
		WHERE bucket_id=$1 AND name=$2 AND (($3::char(16) IS NULL AND parent_directory_id IS NULL) OR (parent_directory_id = $3))
	`, bucketID, name, parentDirID)
	var d model.Directory
	if err := row.Scan(
		&d.ID, &d.BucketID, &d.ParentDirectoryID, &d.Name, &d.MetaData, &d.EncryptedMetaData,
		&d.CreatedByUserID, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DirectoryRepository) FindRootByBucketID(ctx context.Context, bucketID string) (*model.Directory, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       created_by_user_id, created_at, updated_at
		FROM directories
		WHERE bucket_id=$1 AND parent_directory_id IS NULL
	`, bucketID)
	var d model.Directory
	if err := row.Scan(
		&d.ID, &d.BucketID, &d.ParentDirectoryID, &d.Name, &d.MetaData, &d.EncryptedMetaData,
		&d.CreatedByUserID, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DirectoryRepository) ListRootsByBucketIDs(ctx context.Context, bucketIDs []string) ([]model.Directory, error) {
	if len(bucketIDs) == 0 {
		return nil, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       created_by_user_id, created_at, updated_at
		FROM directories
		WHERE bucket_id = ANY($1) AND parent_directory_id IS NULL
	`, bucketIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Directory
	for rows.Next() {
		var d model.Directory
		if err := rows.Scan(
			&d.ID, &d.BucketID, &d.ParentDirectoryID, &d.Name, &d.MetaData, &d.EncryptedMetaData,
			&d.CreatedByUserID, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *DirectoryRepository) ListChildDirectories(ctx context.Context, bucketID, parentDirID string) ([]model.Directory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data,
		       created_by_user_id, created_at, updated_at
		FROM directories
		WHERE bucket_id=$1 AND parent_directory_id=$2
		ORDER BY name
	`, bucketID, parentDirID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Directory
	for rows.Next() {
		var d model.Directory
		if err := rows.Scan(
			&d.ID, &d.BucketID, &d.ParentDirectoryID, &d.Name, &d.MetaData, &d.EncryptedMetaData,
			&d.CreatedByUserID, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *DirectoryRepository) Create(ctx context.Context, d *model.Directory) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO directories (id, bucket_id, parent_directory_id, name, meta_data, encrypted_meta_data, created_by_user_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, d.ID, d.BucketID, d.ParentDirectoryID, d.Name, d.MetaData, d.EncryptedMetaData, d.CreatedByUserID)
	return err
}

func (r *DirectoryRepository) UpdateName(ctx context.Context, bucketID, id, name string) error {
	_, err := r.db.Exec(ctx, `UPDATE directories SET name=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, id, name)
	return err
}

func (r *DirectoryRepository) UpdateMetaData(ctx context.Context, bucketID, id string, metaData []byte) error {
	_, err := r.db.Exec(ctx, `UPDATE directories SET meta_data=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, id, metaData)
	return err
}

func (r *DirectoryRepository) UpdateEncryptedMetaData(ctx context.Context, bucketID, id, encryptedMetaData string) error {
	_, err := r.db.Exec(ctx, `UPDATE directories SET encrypted_meta_data=$3, updated_at=NOW() WHERE bucket_id=$1 AND id=$2`, bucketID, id, encryptedMetaData)
	return err
}

func (r *DirectoryRepository) Move(ctx context.Context, bucketID, id string, newParentID *string, newName string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE directories SET parent_directory_id=$3, name=$4, updated_at=NOW()
		WHERE bucket_id=$1 AND id=$2
	`, bucketID, id, newParentID, newName)
	return err
}

func (r *DirectoryRepository) Delete(ctx context.Context, bucketID, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM directories WHERE bucket_id=$1 AND id=$2`, bucketID, id)
	return err
}
