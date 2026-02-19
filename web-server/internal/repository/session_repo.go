package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, s *model.Session) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (
			id, user_id, api_key_hash, has_expired, expired_at, expire_reason,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		s.ID, s.UserID, s.APIKeyHash, s.HasExpired, s.ExpiredAt, s.ExpireReason,
		s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *SessionRepository) FindSessionByID(ctx context.Context, id string) (*model.Session, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, api_key_hash, has_expired, expired_at, expire_reason,
		       created_at, updated_at
		FROM sessions WHERE id=$1
	`, id)
	var s model.Session
	if err := row.Scan(
		&s.ID, &s.UserID, &s.APIKeyHash, &s.HasExpired, &s.ExpiredAt, &s.ExpireReason,
		&s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) ListSessionsByUserID(ctx context.Context, userID string, limit int32) ([]model.Session, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, api_key_hash, has_expired, expired_at, expire_reason,
		       created_at, updated_at
		FROM sessions
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Session
	for rows.Next() {
		var s model.Session
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.APIKeyHash, &s.HasExpired, &s.ExpiredAt, &s.ExpireReason,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *SessionRepository) ExpireSessionByID(ctx context.Context, id, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sessions
		SET has_expired=true,
		    expired_at=NOW(),
		    expire_reason=$2,
		    updated_at=NOW()
		WHERE id=$1
	`, id, reason)
	return err
}

func (r *SessionRepository) ExpireAllSessionsByUserID(ctx context.Context, userID, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sessions
		SET has_expired=true,
		    expired_at=NOW(),
		    expire_reason=$2,
		    updated_at=NOW()
		WHERE user_id=$1 AND has_expired=false
	`, userID, reason)
	return err
}

