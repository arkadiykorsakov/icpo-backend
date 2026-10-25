package country

import (
	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "countries"

// repository — CRUD справочника стран целиком берётся из refdata.
type repository struct {
	refdata.Repository[Country]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Repository: refdata.NewRepository[Country](db, TableName),
	}
}
