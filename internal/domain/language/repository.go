package language

import (
	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "languages"

// repository — CRUD справочника языков целиком берётся из refdata.
type repository struct {
	refdata.Repository[Language]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Repository: refdata.NewRepository[Language](db, TableName),
	}
}
