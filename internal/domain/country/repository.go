package country

import (
	"context"

	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "countries"

// repository — CRUD справочника стран целиком берётся из refdata.
type repository struct {
	db *pgxpool.Pool
	refdata.Repository[Country]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db:         db,
		Repository: refdata.NewRepository[Country](db, TableName),
	}
}

func (r *repository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM countries WHERE id = $1)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
