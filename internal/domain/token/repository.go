package token

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "icpo-backend/internal/errors"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *repository) GetByHash(ctx context.Context, tokenHash string) (RefreshToken, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1`

	rt := RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return RefreshToken{}, errs.ErrInvalidToken
	}
	if err != nil {
		return RefreshToken{}, err
	}
	return rt, nil
}
func (r *repository) Rotate(
	ctx context.Context,
	oldHash string,
	userID int64,
	newHash string,
	newExpiresAt, now time.Time,
) (bool, error) {
	const revokeQuery = `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2 AND user_id = $3 AND revoked_at IS NULL`
	const insertQuery = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`

	rotated := false
	err := pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, revokeQuery, now, oldHash, userID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return nil
		}
		if _, err := tx.Exec(ctx, insertQuery, userID, newHash, newExpiresAt); err != nil {
			return err
		}
		rotated = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return rotated, nil
}

func (r *repository) Revoke(ctx context.Context, tokenHash string) (bool, error) {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL`
	tag, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *repository) RevokeAllForUser(ctx context.Context, userID int64) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
