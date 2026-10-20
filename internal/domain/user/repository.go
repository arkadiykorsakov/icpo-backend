package user

import (
	"context"

	errs "icpo-backend/internal/errors"
	"icpo-backend/pkg/pgerr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// usersEmailKey — имя UNIQUE-констрейнта на users.email.
// Postgres даёт его по умолчанию (<table>_<column>_key) в миграции
// 20261010151320_create_users_table.sql. При переименовании обновить здесь.
const usersEmailKey = "users_email_key"

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func scanUser(row pgx.Row) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
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

	row := r.db.QueryRow(ctx, query, email)

	u, err := scanUser(row)
	if pgerr.IsNotFound(err) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	const query = `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	u, err := scanUser(row)
	if pgerr.IsNotFound(err) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Create(ctx context.Context, candidate User, passwordHash string) (*User, error) {
	const query = `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, password_hash, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, candidate.Email, candidate.Name, passwordHash)

	u, err := scanUser(row)
	if pgerr.IsUniqueViolation(err, usersEmailKey) {
		return nil, errs.ErrEmailTaken
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Update(ctx context.Context, id int64, candidate User) (*User, error) {
	const query = `
		UPDATE users
		SET email = $1, name = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, name, password_hash, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, candidate.Email, candidate.Name, id)

	u, err := scanUser(row)
	if pgerr.IsNotFound(err) {
		return nil, errs.ErrUserNotFound
	}
	if pgerr.IsUniqueViolation(err, usersEmailKey) {
		return nil, errs.ErrEmailTaken
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	const query = `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}
