package haircolor

import (
	"icpo-backend/internal/shared/refdata"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "hair_colors"

// repository — CRUD справочника цветов волос целиком берётся из refdata.
type repository struct {
	refdata.Repository[HairColor]
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Repository: refdata.NewRepository[HairColor](db, TableName),
	}
}
