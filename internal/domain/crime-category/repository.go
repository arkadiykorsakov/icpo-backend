package crimecategory

import (
	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "crime_categories"

// repository — CRUD справочника категорий преступлений целиком берётся из refdata.
type repository struct {
	refdata.Repository[CrimeCategory]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Repository: refdata.NewRepository[CrimeCategory](db, TableName),
	}
}
