package eyecolor

import (
	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "eye_colors"

// repository — CRUD справочника цветов глаз целиком берётся из refdata.
type repository struct {
	refdata.Repository[EyeColor]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Repository: refdata.NewRepository[EyeColor](db, TableName),
	}
}
