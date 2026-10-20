// pkg/pgerr/pgerr.go
package pgerr

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueViolation — SQLSTATE нарушения UNIQUE-ограничения.
const uniqueViolation = "23505"

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == uniqueViolation &&
		pgErr.ConstraintName == constraint
}
