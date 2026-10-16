package user

import (
	"context"
	"errors"

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

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
        SELECT id, email, name, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	const query = `
	    SELECT id, email, name, password_hash, created_at, updated_at
	    FROM users
	    WHERE id = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) Create(ctx context.Context, email, name, passwordHash string) (*User, error) {
	const query = `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, password_hash, created_at, updated_at
	`
	user := &User{}
	err := r.db.QueryRow(ctx, query, email, name, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
