package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, display_name, user_name, password_hash, password_salt,
		       is_banned, perm_manage_all_user, perm_create_user, perm_create_bucket,
		       created_at, updated_at
		FROM users WHERE id=$1
	`, id)
	var u model.User
	if err := row.Scan(
		&u.ID, &u.DisplayName, &u.UserName, &u.PasswordHash, &u.PasswordSalt,
		&u.IsBanned, &u.PermManageAllUser, &u.PermCreateUser, &u.PermCreateBucket,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, display_name, user_name, password_hash, password_salt,
		       is_banned, perm_manage_all_user, perm_create_user, perm_create_bucket,
		       created_at, updated_at
		FROM users WHERE user_name=$1
	`, userName)
	var u model.User
	if err := row.Scan(
		&u.ID, &u.DisplayName, &u.UserName, &u.PasswordHash, &u.PasswordSalt,
		&u.IsBanned, &u.PermManageAllUser, &u.PermCreateUser, &u.PermCreateBucket,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) ListAllUsers(ctx context.Context) ([]model.UserListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_name, display_name, is_banned
		FROM users
		ORDER BY user_name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.UserListItem
	for rows.Next() {
		var item model.UserListItem
		if err := rows.Scan(&item.ID, &item.UserName, &item.DisplayName, &item.IsBanned); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (
			id, display_name, user_name, password_hash, password_salt,
			is_banned, perm_manage_all_user, perm_create_user, perm_create_bucket
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		u.ID, u.DisplayName, u.UserName, u.PasswordHash, u.PasswordSalt,
		u.IsBanned, u.PermManageAllUser, u.PermCreateUser, u.PermCreateBucket,
	)
	return err
}

func (r *UserRepository) UpdateUserDisplayName(ctx context.Context, id, displayName string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET display_name=$2, updated_at=NOW() WHERE id=$1
	`, id, displayName)
	return err
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id, hash, salt string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET password_hash=$2, password_salt=$3, updated_at=NOW()
		WHERE id=$1
	`, id, hash, salt)
	return err
}

func (r *UserRepository) UpdateUserBanStatus(ctx context.Context, id string, isBanned bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET is_banned=$2, updated_at=NOW() WHERE id=$1
	`, id, isBanned)
	return err
}

func (r *UserRepository) UpdateUserGlobalPermissions(ctx context.Context, id string, manageAll, createUser, createBucket bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET perm_manage_all_user=$2,
		    perm_create_user=$3,
		    perm_create_bucket=$4,
		    updated_at=NOW()
		WHERE id=$1
	`, id, manageAll, createUser, createBucket)
	return err
}

// QueryUsersByIDs returns users whose IDs are in the given slice.
func (r *UserRepository) QueryUsersByIDs(ctx context.Context, ids []string) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, display_name, user_name, password_hash, password_salt,
		       is_banned, perm_manage_all_user, perm_create_user, perm_create_bucket,
		       created_at, updated_at
		FROM users
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.DisplayName, &u.UserName, &u.PasswordHash, &u.PasswordSalt,
			&u.IsBanned, &u.PermManageAllUser, &u.PermCreateUser, &u.PermCreateBucket,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

// QueryUsersByUserNames returns users whose usernames are in the given slice.
func (r *UserRepository) QueryUsersByUserNames(ctx context.Context, names []string) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, display_name, user_name, password_hash, password_salt,
		       is_banned, perm_manage_all_user, perm_create_user, perm_create_bucket,
		       created_at, updated_at
		FROM users
		WHERE user_name = ANY($1)
	`, names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.DisplayName, &u.UserName, &u.PasswordHash, &u.PasswordSalt,
			&u.IsBanned, &u.PermManageAllUser, &u.PermCreateUser, &u.PermCreateBucket,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}


