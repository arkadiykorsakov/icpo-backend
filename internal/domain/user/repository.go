package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "icpo-backend/internal/errors"
)

// pgUniqueViolation — SQLSTATE нарушения UNIQUE-ограничения.
const pgUniqueViolation = "23505"

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func scanUser(row pgx.Row) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	return scanUser(r.db.QueryRow(ctx, query, email))
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	const query = `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	return scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *repository) Create(ctx context.Context, email, name, passwordHash string) (*User, error) {
	const query = `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, password_hash, created_at, updated_at`

	u, err := scanUser(r.db.QueryRow(ctx, query, email, name, passwordHash))

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return nil, errs.ErrEmailTaken
	}
	return u, err
}
