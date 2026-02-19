package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type BucketRepository struct {
	db *pgxpool.Pool
}

func NewBucketRepository(db *pgxpool.Pool) *BucketRepository {
	return &BucketRepository{db: db}
}

func (r *BucketRepository) FindByID(ctx context.Context, id string) (*model.Bucket, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, crypt_spec, crypt_data, meta_data,
		       created_by_user_id, created_at, updated_at
		FROM buckets WHERE id=$1
	`, id)
	var b model.Bucket
	if err := row.Scan(
		&b.ID, &b.Name, &b.CryptSpec, &b.CryptData, &b.MetaData,
		&b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BucketRepository) FindByName(ctx context.Context, name string) (*model.Bucket, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, crypt_spec, crypt_data, meta_data,
		       created_by_user_id, created_at, updated_at
		FROM buckets WHERE name=$1
	`, name)
	var b model.Bucket
	if err := row.Scan(
		&b.ID, &b.Name, &b.CryptSpec, &b.CryptData, &b.MetaData,
		&b.CreatedByUserID, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BucketRepository) Create(ctx context.Context, b *model.Bucket) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO buckets (id, name, crypt_spec, crypt_data, meta_data, created_by_user_id)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, b.ID, b.Name, b.CryptSpec, b.CryptData, b.MetaData, b.CreatedByUserID)
	return err
}

func (r *BucketRepository) UpdateName(ctx context.Context, id, name string) error {
	_, err := r.db.Exec(ctx, `UPDATE buckets SET name=$2, updated_at=NOW() WHERE id=$1`, id, name)
	return err
}

func (r *BucketRepository) UpdateMetaData(ctx context.Context, id string, metaData []byte) error {
	_, err := r.db.Exec(ctx, `UPDATE buckets SET meta_data=$2, updated_at=NOW() WHERE id=$1`, id, metaData)
	return err
}

func (r *BucketRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM buckets WHERE id=$1`, id)
	return err
}

// FindPermission returns the bucket_user_permissions row for (bucketID, userID).
func (r *BucketRepository) FindPermission(ctx context.Context, bucketID, userID string) (*model.BucketPermission, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bucket_id, user_id, notes,
		       perm_modify, perm_manage_authorization, perm_destroy,
		       perm_view_content, perm_manage_content,
		       created_at, updated_at
		FROM bucket_user_permissions
		WHERE bucket_id=$1 AND user_id=$2
	`, bucketID, userID)
	var p model.BucketPermission
	if err := row.Scan(
		&p.ID, &p.BucketID, &p.UserID, &p.Notes,
		&p.PermModify, &p.PermManageAuthorization, &p.PermDestroy,
		&p.PermViewContent, &p.PermManageContent,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPermissionsByBucketID returns all permission rows for a bucket.
func (r *BucketRepository) ListPermissionsByBucketID(ctx context.Context, bucketID string) ([]model.BucketPermission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bucket_id, user_id, notes,
		       perm_modify, perm_manage_authorization, perm_destroy,
		       perm_view_content, perm_manage_content,
		       created_at, updated_at
		FROM bucket_user_permissions
		WHERE bucket_id=$1
	`, bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.BucketPermission
	for rows.Next() {
		var p model.BucketPermission
		if err := rows.Scan(
			&p.ID, &p.BucketID, &p.UserID, &p.Notes,
			&p.PermModify, &p.PermManageAuthorization, &p.PermDestroy,
			&p.PermViewContent, &p.PermManageContent,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// ListBucketIDsByUserID returns bucket IDs for which the user has any permission.
func (r *BucketRepository) ListBucketIDsByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `SELECT bucket_id FROM bucket_user_permissions WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreatePermission inserts a new bucket_user_permissions row.
func (r *BucketRepository) CreatePermission(ctx context.Context, p *model.BucketPermission) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO bucket_user_permissions (
			bucket_id, user_id, notes,
			perm_modify, perm_manage_authorization, perm_destroy,
			perm_view_content, perm_manage_content
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, p.BucketID, p.UserID, p.Notes,
		p.PermModify, p.PermManageAuthorization, p.PermDestroy,
		p.PermViewContent, p.PermManageContent)
	return err
}

// UpdatePermission overwrites the five permission flags for the given bucket/user.
func (r *BucketRepository) UpdatePermission(ctx context.Context, p *model.BucketPermission) error {
	_, err := r.db.Exec(ctx, `
		UPDATE bucket_user_permissions SET
			perm_modify=$3, perm_manage_authorization=$4, perm_destroy=$5,
			perm_view_content=$6, perm_manage_content=$7,
			updated_at=NOW()
		WHERE bucket_id=$1 AND user_id=$2
	`, p.BucketID, p.UserID,
		p.PermModify, p.PermManageAuthorization, p.PermDestroy,
		p.PermViewContent, p.PermManageContent)
	return err
}

func (r *BucketRepository) DeletePermission(ctx context.Context, bucketID, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM bucket_user_permissions WHERE bucket_id=$1 AND user_id=$2`, bucketID, userID)
	return err
}
